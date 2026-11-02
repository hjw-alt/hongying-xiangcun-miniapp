package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxJSONSize = 1 << 20
const maxImageSize = 8 << 20

type API struct {
	Store *MySQLStore
}

type listResponse struct {
	Items []Artwork `json:"items"`
	Total int       `json:"total"`
}

type detailResponse struct {
	Artwork  Artwork   `json:"artwork"`
	Comments []Comment `json:"comments"`
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Admin-Token, X-User-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requestContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Second)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func handleStoreError(w http.ResponseWriter, err error, action string) {
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "作品或评论不存在")
		return
	}
	log.Printf("%s: %v", action, err)
	writeError(w, http.StatusInternalServerError, action+"失败，请稍后再试")
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DB.PingContext(ctx); err != nil {
		writeError(w, http.StatusServiceUnavailable, "数据库连接不可用")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "shadowart-xiangtan", "database": "mysql"})
}

func (a *API) handleStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	result, err := a.Store.Stats(ctx)
	if err != nil {
		handleStoreError(w, err, "读取统计")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (a *API) handleListArtworks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListArtworks(ctx,
		r.URL.Query().Get("keyword"), r.URL.Query().Get("theme"), r.URL.Query().Get("sort"), "")
	if err != nil {
		handleStoreError(w, err, "读取作品列表")
		return
	}
	writeJSON(w, http.StatusOK, listResponse{Items: items, Total: len(items)})
}

func normalizeArtwork(input Artwork) (Artwork, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)
	input.Theme = strings.TrimSpace(input.Theme)
	input.Description = strings.TrimSpace(input.Description)
	input.ImageURL = strings.TrimSpace(input.ImageURL)
	if input.Title == "" {
		return Artwork{}, errors.New("作品名称不能为空")
	}
	if len([]rune(input.Title)) > 40 {
		return Artwork{}, errors.New("作品名称不能超过 40 个字")
	}
	if input.Author == "" {
		return Artwork{}, errors.New("作者不能为空")
	}
	if len([]rune(input.Author)) > 30 {
		return Artwork{}, errors.New("作者不能超过 30 个字")
	}
	if input.Theme != "纸影人物" && input.Theme != "湘潭印象" && input.Theme != "戏台故事" {
		return Artwork{}, errors.New("请选择纸影人物、湘潭印象或戏台故事")
	}
	if len([]rune(input.Description)) > 500 {
		return Artwork{}, errors.New("作品介绍不能超过 500 个字")
	}
	if input.ImageURL == "" {
		input.ImageURL = defaultCover
	}
	cleanTags := make([]string, 0, len(input.Tags))
	for _, tag := range input.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if len([]rune(tag)) > 12 {
			return Artwork{}, errors.New("单个标签不能超过 12 个字")
		}
		if len(cleanTags) < 6 {
			cleanTags = append(cleanTags, tag)
		}
	}
	input.Tags = cleanTags
	return input, nil
}

func (a *API) handleCreateArtwork(w http.ResponseWriter, r *http.Request) {
	var input Artwork
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeArtwork(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = newID()
	item.OwnerKey = userKey(r)
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateArtwork(ctx, item); err != nil {
		handleStoreError(w, err, "创建作品")
		return
	}
	item, err = a.Store.GetArtwork(ctx, item.ID)
	if err != nil {
		handleStoreError(w, err, "读取新作品")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxJSONSize))
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "请求数据格式不正确")
		return err
	}
	return nil
}

func (a *API) detailPayload(ctx context.Context, id string) (detailResponse, error) {
	item, err := a.Store.GetArtwork(ctx, id)
	if err != nil {
		return detailResponse{}, err
	}
	comments, err := a.Store.ListComments(ctx, id)
	if err != nil {
		return detailResponse{}, err
	}
	return detailResponse{Artwork: item, Comments: comments}, nil
}

func (a *API) handleGetArtwork(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	payload, err := a.detailPayload(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取作品详情")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (a *API) handleUpdateArtwork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input Artwork
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeArtwork(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = id
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.UpdateArtwork(ctx, item); err != nil {
		handleStoreError(w, err, "更新作品")
		return
	}
	item, err = a.Store.GetArtwork(ctx, id)
	if err != nil {
		handleStoreError(w, err, "读取更新结果")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleDeleteArtwork(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteArtwork(ctx, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除作品")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleLikeArtwork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.LikeArtwork(ctx, id); err != nil {
		handleStoreError(w, err, "作品点赞")
		return
	}
	item, err := a.Store.GetArtwork(ctx, id)
	if err != nil {
		handleStoreError(w, err, "读取点赞结果")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	artworkID := r.PathValue("id")
	var input Comment
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	input.Nickname = strings.TrimSpace(input.Nickname)
	input.Content = strings.TrimSpace(input.Content)
	if input.Nickname == "" {
		input.Nickname = "匿名同学"
	}
	if len([]rune(input.Nickname)) > 20 {
		writeError(w, http.StatusBadRequest, "昵称不能超过 20 个字")
		return
	}
	if input.Content == "" || len([]rune(input.Content)) > 200 {
		writeError(w, http.StatusBadRequest, "评论需为 1-200 个字")
		return
	}
	input.ID = newID()
	input.ArtworkID = artworkID
	ctx, cancel := requestContext(r)
	defer cancel()
	if _, err := a.Store.GetArtwork(ctx, artworkID); err != nil {
		handleStoreError(w, err, "校验评论作品")
		return
	}
	if err := a.Store.CreateComment(ctx, input); err != nil {
		handleStoreError(w, err, "创建评论")
		return
	}
	item, err := a.Store.ListComments(ctx, artworkID)
	if err != nil {
		handleStoreError(w, err, "读取评论")
		return
	}
	for _, comment := range item {
		if comment.ID == input.ID {
			writeJSON(w, http.StatusCreated, comment)
			return
		}
	}
	writeJSON(w, http.StatusCreated, input)
}

func (a *API) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	err := a.Store.DeleteComment(ctx, r.PathValue("id"), r.PathValue("commentID"))
	if err != nil {
		handleStoreError(w, err, "删除评论")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)
	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		writeError(w, http.StatusBadRequest, "图片不能超过 8MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "缺少 file 字段")
		return
	}
	defer file.Close()
	ext, ok := imageExtension(header)
	if !ok {
		writeError(w, http.StatusBadRequest, "仅支持 jpg、png 或 webp 图片")
		return
	}
	uploadDir := filepath.Join("static", "uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		handleStoreError(w, err, "创建图片目录")
		return
	}
	filename := newID() + ext
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		handleStoreError(w, err, "保存图片")
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filepath.Join(uploadDir, filename))
		handleStoreError(w, err, "写入图片")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"url": "/static/uploads/" + filename})
}

func imageExtension(header *multipart.FileHeader) (string, bool) {
	switch header.Header.Get("Content-Type") {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/webp":
		return ".webp", true
	default:
		return "", false
	}
}

func newID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}
	return hex.EncodeToString(raw)
}
