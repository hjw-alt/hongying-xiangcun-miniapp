package main

import (
	"errors"
	"net/http"
	"strings"
)

func normalizeAddress(input Address) (Address, error) {
	input.Receiver = strings.TrimSpace(input.Receiver)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Region = strings.TrimSpace(input.Region)
	input.Detail = strings.TrimSpace(input.Detail)
	if input.Receiver == "" || len([]rune(input.Receiver)) > 20 {
		return Address{}, errors.New("收货人需为 1-20 个字")
	}
	if !phonePattern.MatchString(input.Phone) {
		return Address{}, errors.New("请输入正确的手机号")
	}
	if len([]rune(input.Region)) > 60 {
		return Address{}, errors.New("地区不能超过 60 个字")
	}
	if len([]rune(input.Detail)) < 4 || len([]rune(input.Detail)) > 100 {
		return Address{}, errors.New("详细地址需为 4-100 个字")
	}
	return input, nil
}

func (a *API) handleListAddresses(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListAddresses(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取收货地址")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input Address
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeAddress(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.UserKey = key
	ctx, cancel := requestContext(r)
	defer cancel()
	created, err := a.Store.CreateAddress(ctx, item)
	if err != nil {
		handleStoreError(w, err, "保存收货地址")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (a *API) handleUpdateAddress(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input Address
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item, err := normalizeAddress(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ID = r.PathValue("id")
	item.UserKey = key
	ctx, cancel := requestContext(r)
	defer cancel()
	updated, err := a.Store.UpdateAddress(ctx, item)
	if err != nil {
		handleStoreError(w, err, "更新收货地址")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *API) handleDeleteAddress(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteAddress(ctx, key, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除收货地址")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleSetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.SetDefaultAddress(ctx, key, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "设置默认地址")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleListCart(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListCartItems(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取购物车")
		return
	}
	total := 0
	for _, item := range items {
		if item.Product != nil {
			total += item.Product.PriceCents * item.Quantity
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": total})
}

func (a *API) handleAddCartItem(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input CartItem
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	input.ProductID = strings.TrimSpace(input.ProductID)
	input.CustomText = strings.TrimSpace(input.CustomText)
	if input.ProductID == "" {
		writeError(w, http.StatusBadRequest, "缺少商品信息")
		return
	}
	if len([]rune(input.CustomText)) > 150 {
		writeError(w, http.StatusBadRequest, "定制要求不能超过 150 个字")
		return
	}
	input.UserKey = key
	if input.Quantity < 1 {
		input.Quantity = 1
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.AddCartItem(ctx, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusBadRequest, "商品不存在或已下架")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	count, _ := a.Store.CountCartItems(ctx, key)
	writeJSON(w, http.StatusCreated, map[string]any{"item": item, "count": count})
}

func (a *API) handleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input struct {
		Quantity   int    `json:"quantity"`
		CustomText string `json:"customText"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.UpdateCartItem(ctx, key, r.PathValue("id"), input.Quantity, input.CustomText); err != nil {
		handleStoreError(w, err, "更新购物车")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *API) handleDeleteCartItem(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.DeleteCartItem(ctx, key, r.PathValue("id")); err != nil {
		handleStoreError(w, err, "删除购物车商品")
		return
	}
	count, _ := a.Store.CountCartItems(ctx, key)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "count": count})
}

func (a *API) handleCartCount(w http.ResponseWriter, r *http.Request) {
	key := userKey(r)
	if key == "" {
		writeJSON(w, http.StatusOK, map[string]int{"count": 0})
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	count, err := a.Store.CountCartItems(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取购物车数量")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}
