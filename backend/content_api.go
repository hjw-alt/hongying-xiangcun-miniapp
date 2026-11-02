package main

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const maxVideoSize = 100 << 20

var userKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
var phonePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

func userKey(r *http.Request) string {
	key := r.Header.Get("X-User-Key")
	if key == "" {
		key = r.URL.Query().Get("userKey")
	}
	key = strings.TrimSpace(key)
	if !userKeyPattern.MatchString(key) {
		return ""
	}
	return key
}

func requireUserKey(w http.ResponseWriter, r *http.Request) string {
	key := userKey(r)
	if key == "" {
		writeError(w, http.StatusUnauthorized, "缺少用户标识")
	}
	return key
}

func clipList(items []Clip) []Clip {
	if items == nil {
		return []Clip{}
	}
	return items
}

func (a *API) handleListClips(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListClips(ctx)
	if err != nil {
		handleStoreError(w, err, "读取影戏剧目")
		return
	}
	writeJSON(w, http.StatusOK, clipList(items))
}

func normalizeClip(input Clip) (Clip, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Uploader = strings.TrimSpace(input.Uploader)
	input.Description = strings.TrimSpace(input.Description)
	input.VideoURL = strings.TrimSpace(input.VideoURL)
	input.CoverURL = strings.TrimSpace(input.CoverURL)
	if input.Title == "" {
		return Clip{}, errors.New("剧目标题不能为空")
	}
	if len([]rune(input.Title)) > 40 {
		return Clip{}, errors.New("剧目标题不能超过 40 个字")
	}
	if len([]rune(input.Description)) > 500 {
		return Clip{}, errors.New("剧目介绍不能超过 500 个字")
	}
	if input.VideoURL == "" {
		return Clip{}, errors.New("请上传影戏片段")
	}
	if input.Uploader == "" {
		input.Uploader = "匿名同学"
	}
	if input.CoverURL == "" {
		input.CoverURL = defaultCover
	}
	return input, nil
}

func (a *API) handleCreateClip(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	var input Clip
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeClip(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = newID()
	item.OwnerKey = userKey(r)
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateClip(ctx, item); err != nil {
		handleStoreError(w, err, "上传剧目")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) handleLikeClip(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.LikeClip(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "剧目点赞")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleFavorite(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	itemType := r.PathValue("type")
	itemID := r.PathValue("id")
	if itemType != "artwork" && itemType != "clip" {
		writeError(w, http.StatusBadRequest, "收藏类型不正确")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	favorited, err := a.Store.ToggleFavorite(ctx, key, itemType, itemID)
	if err != nil {
		handleStoreError(w, err, "更新收藏")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"favorited": favorited})
}

func (a *API) handleFavoriteStatus(w http.ResponseWriter, r *http.Request) {
	key := userKey(r)
	itemType := r.PathValue("type")
	itemID := r.PathValue("id")
	if itemType != "artwork" && itemType != "clip" {
		writeError(w, http.StatusBadRequest, "收藏类型不正确")
		return
	}
	if key == "" {
		writeJSON(w, http.StatusOK, map[string]bool{"favorited": false})
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	favorited, err := a.Store.FavoriteStatus(ctx, key, itemType, itemID)
	if err != nil {
		handleStoreError(w, err, "读取收藏状态")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"favorited": favorited})
}

func (a *API) handleListFavoriteArtworks(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListFavoriteArtworks(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取收藏影偶")
		return
	}
	writeJSON(w, http.StatusOK, listResponse{Items: items, Total: len(items)})
}

func (a *API) handleListFavoriteClips(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListFavoriteClips(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取收藏剧目")
		return
	}
	writeJSON(w, http.StatusOK, clipList(items))
}

func (a *API) handleListActivities(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListActivities(ctx)
	if err != nil {
		handleStoreError(w, err, "读取研学活动")
		return
	}
	if items == nil {
		items = []Activity{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleGetActivity(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.GetActivity(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取研学活动详情")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input Booking
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	if input.Name == "" || len([]rune(input.Name)) > 20 {
		writeError(w, http.StatusBadRequest, "联系人需为 1-20 个字")
		return
	}
	if !phonePattern.MatchString(input.Phone) {
		writeError(w, http.StatusBadRequest, "请输入正确的手机号")
		return
	}
	if input.Participants < 1 || input.Participants > 10 {
		writeError(w, http.StatusBadRequest, "预约人数需为 1-10 人")
		return
	}
	input.OwnerKey = key
	input.ID = newID()
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.CreateBooking(ctx, input)
	if err != nil {
		if err.Error() == "名额不足" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleStoreError(w, err, "创建预约")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) handleListBookings(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListBookings(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取预约记录")
		return
	}
	if items == nil {
		items = []Booking{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleCancelBooking(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CancelBooking(ctx, key, r.PathValue("id")); err != nil {
		if err.Error() == "该订单已取消" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleStoreError(w, err, "取消预约")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleUploadVideo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxVideoSize)
	if err := r.ParseMultipartForm(maxVideoSize); err != nil {
		writeError(w, http.StatusBadRequest, "视频不能超过 100MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "缺少 file 字段")
		return
	}
	defer file.Close()
	ext, ok := videoExtension(header)
	if !ok {
		writeError(w, http.StatusBadRequest, "仅支持 mp4、m4v 或 mov 视频")
		return
	}
	url, err := saveUploadedFile(file, ext)
	if err != nil {
		handleStoreError(w, err, "保存视频")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"url": url})
}

func videoExtension(header *multipart.FileHeader) (string, bool) {
	switch header.Header.Get("Content-Type") {
	case "video/mp4":
		return ".mp4", true
	case "video/x-m4v":
		return ".m4v", true
	case "video/quicktime":
		return ".mov", true
	default:
		return "", false
	}
}

func saveUploadedFile(file multipart.File, ext string) (string, error) {
	uploadDir := filepath.Join("static", "uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}
	filename := newID() + ext
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filepath.Join(uploadDir, filename))
		return "", err
	}
	return "/static/uploads/" + filename, nil
}

func (a *API) handleMyArtworks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListArtworks(ctx,
		r.URL.Query().Get("keyword"), r.URL.Query().Get("theme"), r.URL.Query().Get("sort"), userKey(r))
	if err != nil {
		handleStoreError(w, err, "读取我的作品")
		return
	}
	writeJSON(w, http.StatusOK, listResponse{Items: items, Total: len(items)})
}
