package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func serializeHighlights(items []string) string {
	cleaned := []string{}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if len(cleaned) < 8 {
			cleaned = append(cleaned, item)
		}
	}
	encoded, err := json.Marshal(cleaned)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func normalizeActivity(input Activity) (Activity, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Category = strings.TrimSpace(input.Category)
	input.Venue = strings.TrimSpace(input.Venue)
	input.Description = strings.TrimSpace(input.Description)
	input.Content = strings.TrimSpace(input.Content)
	input.Notice = strings.TrimSpace(input.Notice)
	input.CoverURL = strings.TrimSpace(input.CoverURL)

	if input.Title == "" {
		return Activity{}, errors.New("活动名称不能为空")
	}
	if len([]rune(input.Title)) > 60 {
		return Activity{}, errors.New("活动名称不能超过 60 个字")
	}
	if input.Category == "" {
		return Activity{}, errors.New("请填写活动分类")
	}
	if len([]rune(input.Category)) > 20 {
		return Activity{}, errors.New("活动分类不能超过 20 个字")
	}
	if input.Venue == "" {
		return Activity{}, errors.New("活动地点不能为空")
	}
	if len([]rune(input.Venue)) > 120 {
		return Activity{}, errors.New("活动地点不能超过 120 个字")
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return Activity{}, errors.New("请填写活动开始与结束时间")
	}
	if !input.EndTime.After(input.StartTime) {
		return Activity{}, errors.New("结束时间需晚于开始时间")
	}
	if input.PriceCents < 0 {
		return Activity{}, errors.New("费用不能为负数")
	}
	if input.Capacity < 0 {
		return Activity{}, errors.New("名额不能为负数")
	}
	if input.Remaining < 0 {
		input.Remaining = 0
	}
	if input.Remaining > input.Capacity {
		input.Remaining = input.Capacity
	}
	if len([]rune(input.Description)) > 500 {
		return Activity{}, errors.New("活动简介不能超过 500 个字")
	}
	if len([]rune(input.Content)) > 5000 {
		return Activity{}, errors.New("活动介绍不能超过 5000 个字")
	}
	if len([]rune(input.Notice)) > 1000 {
		return Activity{}, errors.New("参与须知不能超过 1000 个字")
	}
	for _, item := range input.Highlights {
		if len([]rune(strings.TrimSpace(item))) > 50 {
			return Activity{}, errors.New("单条活动亮点不能超过 50 个字")
		}
	}
	if input.CoverURL == "" {
		input.CoverURL = defaultCover
	}
	return input, nil
}

func (a *API) handleAdminCreateArtwork(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
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
	item.OwnerKey = "admin"
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateArtwork(ctx, item); err != nil {
		handleStoreError(w, err, "创建作品")
		return
	}
	created, err := a.Store.GetArtwork(ctx, item.ID)
	if err != nil {
		handleStoreError(w, err, "读取新作品")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (a *API) handleAdminUpdateArtwork(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
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
	updated, err := a.Store.GetArtwork(ctx, id)
	if err != nil {
		handleStoreError(w, err, "读取更新结果")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *API) handleAdminDeleteArtwork(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteArtwork(ctx, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除作品")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleAdminCreateClip(w http.ResponseWriter, r *http.Request) {
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
	item.OwnerKey = "admin"
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateClip(ctx, item); err != nil {
		handleStoreError(w, err, "创建剧目")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) handleAdminUpdateClip(w http.ResponseWriter, r *http.Request) {
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
	item.ID = r.PathValue("id")
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.UpdateClip(ctx, item); err != nil {
		handleStoreError(w, err, "更新剧目")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleAdminDeleteClip(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteClip(ctx, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除剧目")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleAdminCreateActivity(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	var input Activity
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeActivity(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if item.ID == "" {
		item.ID = newID()
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.CreateActivity(ctx, item, serializeHighlights(item.Highlights)); err != nil {
		handleStoreError(w, err, "创建研学活动")
		return
	}
	created, err := a.Store.GetActivity(ctx, item.ID)
	if err != nil {
		handleStoreError(w, err, "读取新活动")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (a *API) handleAdminUpdateActivity(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	var input Activity
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeActivity(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = r.PathValue("id")
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.UpdateActivity(ctx, item, serializeHighlights(item.Highlights)); err != nil {
		handleStoreError(w, err, "更新研学活动")
		return
	}
	updated, err := a.Store.GetActivity(ctx, item.ID)
	if err != nil {
		handleStoreError(w, err, "读取更新结果")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *API) handleAdminDeleteActivity(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteActivity(ctx, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除研学活动")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
