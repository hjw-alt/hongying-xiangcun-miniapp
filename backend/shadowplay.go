package main

import (
	"context"
	"database/sql"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ShadowPlay struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Preset    string    `json:"preset"`
	Mode      string    `json:"mode"`
	SourceURL string    `json:"sourceUrl"`
	GIFURL    string    `json:"gifUrl"`
	PosterURL string    `json:"posterUrl"`
	Frames    int       `json:"frames"`
	OwnerKey  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

var validPresets = map[string]bool{
	"亮相": true,
	"舞剑": true,
	"赶路": true,
}

func (s *MySQLStore) migrateShadowPlays() error {
	statement := `CREATE TABLE IF NOT EXISTS shadow_plays (
		id CHAR(36) PRIMARY KEY,
		owner_key VARCHAR(64) NOT NULL DEFAULT '',
		title VARCHAR(60) NOT NULL,
		author VARCHAR(30) NOT NULL DEFAULT '匿名同学',
		preset VARCHAR(20) NOT NULL DEFAULT '亮相',
		mode VARCHAR(10) NOT NULL DEFAULT '彩色',
		source_url VARCHAR(500) NOT NULL DEFAULT '',
		gif_url VARCHAR(500) NOT NULL DEFAULT '',
		poster_url VARCHAR(500) NOT NULL DEFAULT '',
		frames INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_shadow_plays_owner (owner_key, created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
	if _, err := s.DB.Exec(statement); err != nil {
		return err
	}
	exists, err := s.columnExists("shadow_plays", "mode")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.DB.Exec("ALTER TABLE shadow_plays ADD COLUMN mode VARCHAR(10) NOT NULL DEFAULT '彩色' AFTER preset"); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLStore) CreateShadowPlay(ctx context.Context, item ShadowPlay) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO shadow_plays
		(id, owner_key, title, author, preset, mode, source_url, gif_url, poster_url, frames)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.OwnerKey, item.Title, item.Author, item.Preset, item.Mode,
		item.SourceURL, item.GIFURL, item.PosterURL, item.Frames)
	return err
}

func scanShadowPlays(rows *sql.Rows) ([]ShadowPlay, error) {
	items := []ShadowPlay{}
	for rows.Next() {
		var item ShadowPlay
		if err := rows.Scan(&item.ID, &item.Title, &item.Author, &item.Preset, &item.Mode, &item.SourceURL,
			&item.GIFURL, &item.PosterURL, &item.Frames, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListShadowPlays(ctx context.Context, ownerKey string) ([]ShadowPlay, error) {
	query := `SELECT id, title, author, preset, mode, source_url, gif_url, poster_url, frames, created_at
		FROM shadow_plays`
	args := []any{}
	if ownerKey != "" {
		query += ` WHERE owner_key=?`
		args = append(args, ownerKey)
	}
	query += ` ORDER BY created_at DESC LIMIT 40`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShadowPlays(rows)
}

func (s *MySQLStore) GetShadowPlay(ctx context.Context, id string) (ShadowPlay, error) {
	var item ShadowPlay
	err := s.DB.QueryRowContext(ctx, `SELECT id, title, author, preset, mode, source_url, gif_url, poster_url, frames, created_at
		FROM shadow_plays WHERE id=?`, id).
		Scan(&item.ID, &item.Title, &item.Author, &item.Preset, &item.Mode, &item.SourceURL,
			&item.GIFURL, &item.PosterURL, &item.Frames, &item.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ShadowPlay{}, ErrNotFound
	}
	return item, err
}

func (s *MySQLStore) DeleteShadowPlay(ctx context.Context, id, ownerKey string) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM shadow_plays WHERE id=? AND owner_key=?`, id, ownerKey)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func saveGeneratedFile(data []byte, ext string) (string, error) {
	uploadDir := filepath.Join("static", "uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}
	filename := newID() + ext
	if err := os.WriteFile(filepath.Join(uploadDir, filename), data, 0o644); err != nil {
		return "", err
	}
	return "/static/uploads/" + filename, nil
}

func localStaticPath(url string) (string, error) {
	cleaned := strings.TrimSpace(url)
	if !strings.HasPrefix(cleaned, "/static/") {
		return "", errors.New("请先上传一张照片")
	}
	relative := strings.TrimPrefix(cleaned, "/static/")
	full := filepath.Join("static", filepath.FromSlash(relative))
	absolute, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	root, err := filepath.Abs("static")
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absolute, root) {
		return "", errors.New("照片路径不合法")
	}
	return absolute, nil
}

func decodeSourceImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("找不到上传的照片")
	}
	defer file.Close()
	src, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.New("照片格式无法解析，请换一张 JPG 或 PNG 图片")
	}
	return src, nil
}

func normalizeShadowPlay(input ShadowPlay) (ShadowPlay, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)
	input.Preset = strings.TrimSpace(input.Preset)
	input.Mode = strings.TrimSpace(input.Mode)
	input.SourceURL = strings.TrimSpace(input.SourceURL)
	if input.Title == "" {
		input.Title = "我的纸影戏"
	}
	if len([]rune(input.Title)) > 60 {
		return ShadowPlay{}, errors.New("剧目名称不能超过 60 个字")
	}
	if input.Author == "" {
		input.Author = "匿名同学"
	}
	if len([]rune(input.Author)) > 30 {
		return ShadowPlay{}, errors.New("署名不能超过 30 个字")
	}
	if !validPresets[input.Preset] {
		input.Preset = "亮相"
	}
	if input.Mode != "墨影" {
		input.Mode = "彩色"
	}
	if input.SourceURL == "" {
		return ShadowPlay{}, errors.New("请先拍照或选择一张照片")
	}
	return input, nil
}

func (a *API) handleCreateShadowPlay(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input ShadowPlay
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeShadowPlay(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sourcePath, err := localStaticPath(item.SourceURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	src, err := decodeSourceImage(sourcePath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	gifData, posterData, frames, err := GenerateShadowPlay(src, item.Preset, item.Mode != "墨影")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "生成纸影戏失败，请换一张光线更清晰的照片")
		return
	}
	gifURL, err := saveGeneratedFile(gifData, ".gif")
	if err != nil {
		handleStoreError(w, err, "保存影戏短片")
		return
	}
	posterURL, err := saveGeneratedFile(posterData, ".jpg")
	if err != nil {
		handleStoreError(w, err, "保存影戏封面")
		return
	}

	item.ID = newID()
	item.OwnerKey = key
	item.GIFURL = gifURL
	item.PosterURL = posterURL
	item.Frames = frames

	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateShadowPlay(ctx, item); err != nil {
		handleStoreError(w, err, "保存影戏记录")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) handleListShadowPlays(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListShadowPlays(ctx, userKey(r))
	if err != nil {
		handleStoreError(w, err, "读取我的影戏")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleGetShadowPlay(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.GetShadowPlay(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取影戏详情")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleDeleteShadowPlay(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteShadowPlay(ctx, r.PathValue("id"), key); err != nil {
		handleStoreError(w, err, "删除影戏")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
