package main

import (
	"errors"
	"net/http"
	"strings"
)

func (a *API) handleListBanners(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListBanners(ctx)
	if err != nil {
		handleStoreError(w, err, "读取首页轮播")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleListVillages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListVillages(ctx, r.URL.Query().Get("keyword"))
	if err != nil {
		handleStoreError(w, err, "读取红色村落")
		return
	}
	writeJSON(w, http.StatusOK, listResponseVillage(items))
}

type villageList struct {
	Items []Village `json:"items"`
	Total int       `json:"total"`
}

func listResponseVillage(items []Village) villageList {
	if items == nil {
		items = []Village{}
	}
	return villageList{Items: items, Total: len(items)}
}

func (a *API) handleGetVillage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	village, err := a.Store.GetVillage(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取村落详情")
		return
	}
	spots, err := a.Store.ListSpots(ctx, village.ID)
	if err != nil {
		handleStoreError(w, err, "读取村落景点")
		return
	}
	stories, err := a.Store.ListStories(ctx, village.ID, "", "")
	if err != nil {
		handleStoreError(w, err, "读取村落故事")
		return
	}
	products, err := a.Store.ListProducts(ctx, "", "", "", "")
	if err != nil {
		handleStoreError(w, err, "读取村落商品")
		return
	}
	related := []Product{}
	for _, item := range products {
		if item.VillageID == village.ID {
			related = append(related, item)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"village":  village,
		"spots":    spots,
		"stories":  stories,
		"products": related,
	})
}

func (a *API) handleListSpots(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListSpots(ctx, r.URL.Query().Get("villageId"))
	if err != nil {
		handleStoreError(w, err, "读取景点")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleListRoutes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListRoutes(ctx)
	if err != nil {
		handleStoreError(w, err, "读取红色线路")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleListStories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListStories(ctx, r.URL.Query().Get("villageId"), r.URL.Query().Get("category"), r.URL.Query().Get("keyword"))
	if err != nil {
		handleStoreError(w, err, "读取红迹故事")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func (a *API) handleGetStory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.GetStory(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取故事详情")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListProducts(ctx,
		r.URL.Query().Get("category"), r.URL.Query().Get("keyword"),
		r.URL.Query().Get("sort"), r.URL.Query().Get("farmerId"))
	if err != nil {
		handleStoreError(w, err, "读取商品列表")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func (a *API) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.GetProduct(ctx, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取商品详情")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleListFarmers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListFarmers(ctx)
	if err != nil {
		handleStoreError(w, err, "读取农户信息")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleFarmSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()
	summary, err := a.Store.FarmSummary(ctx)
	if err != nil {
		handleStoreError(w, err, "读取助农统计")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (a *API) handleListTasks(w http.ResponseWriter, r *http.Request) {
	key := userKey(r)
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListTasks(ctx, r.URL.Query().Get("villageId"))
	if err != nil {
		handleStoreError(w, err, "读取研学任务")
		return
	}
	records := []TaskRecord{}
	if key != "" {
		records, err = a.Store.ListTaskRecords(ctx, key)
		if err != nil {
			handleStoreError(w, err, "读取打卡记录")
			return
		}
	}
	done := map[string]bool{}
	points := 0
	var profile UserProfile
	if key != "" {
		for _, record := range records {
			done[record.TaskID] = true
		}
		profile, err = a.Store.GetUserProfile(ctx, key)
		if err != nil {
			handleStoreError(w, err, "读取成长账户")
			return
		}
		points = profile.Points
		if profile.Medals == nil {
			profile.Medals = []Medal{}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":   items,
		"done":    done,
		"points":  points,
		"records": records,
		"profile": profile,
	})
}

func (a *API) handleCheckinTask(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input TaskRecord
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	input.TaskID = r.PathValue("id")
	input.UserKey = key
	input.Note = strings.TrimSpace(input.Note)
	if len([]rune(input.Note)) > 300 {
		writeError(w, http.StatusBadRequest, "打卡心得不能超过 300 个字")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	record, err := a.Store.CreateTaskRecord(ctx, input)
	if err != nil {
		if err.Error() == "该任务已经打卡" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleStoreError(w, err, "提交打卡")
		return
	}
	profile, profileErr := a.Store.GetUserProfile(ctx, key)
	if profileErr != nil {
		handleStoreError(w, profileErr, "读取成长账户")
		return
	}
	medals, medalsErr := a.Store.ListMedals(ctx, key)
	if medalsErr != nil {
		handleStoreError(w, medalsErr, "读取徽章")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"record": record, "profile": profile, "medals": medals})
}

func (a *API) handleListTaskRecords(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListTaskRecords(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取打卡记录")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	var input struct {
		Order
		AddressID   string   `json:"addressId"`
		CartItemIDs []string `json:"cartItemIds"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	item := input.Order
	item.Receiver = strings.TrimSpace(item.Receiver)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Address = strings.TrimSpace(item.Address)
	item.Remark = strings.TrimSpace(item.Remark)

	ctx, cancel := requestContext(r)
	defer cancel()

	// 支持直接选择已保存的收货地址
	if strings.TrimSpace(input.AddressID) != "" {
		address, err := a.Store.GetAddress(ctx, key, strings.TrimSpace(input.AddressID))
		if err != nil {
			writeError(w, http.StatusBadRequest, "收货地址不存在，请重新选择")
			return
		}
		item.Receiver = address.Receiver
		item.Phone = address.Phone
		if address.Region != "" {
			item.Address = address.Region + " " + address.Detail
		} else {
			item.Address = address.Detail
		}
	}

	if item.Receiver == "" || len([]rune(item.Receiver)) > 20 {
		writeError(w, http.StatusBadRequest, "请填写 1-20 字的收货人姓名")
		return
	}
	if !phonePattern.MatchString(item.Phone) {
		writeError(w, http.StatusBadRequest, "请输入正确的手机号")
		return
	}
	if len([]rune(item.Address)) < 5 {
		writeError(w, http.StatusBadRequest, "请填写详细收货地址")
		return
	}
	if len(item.Items) == 0 || len(item.Items) > 10 {
		writeError(w, http.StatusBadRequest, "订单商品数量不正确")
		return
	}
	item.UserKey = key
	item.SupportFarm = true
	created, err := a.Store.CreateOrder(ctx, item)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusBadRequest, "商品不存在或已下架")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	// 从购物车下单时，清理已结算条目
	if len(input.CartItemIDs) > 0 {
		_ = a.Store.DeleteCartItems(ctx, key, input.CartItemIDs)
	}
	profile, profileErr := a.Store.GetUserProfile(ctx, key)
	if profileErr != nil {
		handleStoreError(w, profileErr, "读取成长账户")
		return
	}
	medals, medalsErr := a.Store.ListMedals(ctx, key)
	if medalsErr != nil {
		handleStoreError(w, medalsErr, "读取徽章")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"order": created, "profile": profile, "medals": medals})
}

func (a *API) handleListOrders(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListOrders(ctx, key, r.URL.Query().Get("status"))
	if err != nil {
		handleStoreError(w, err, "读取订单")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// handleToggleFavorite 支持商品、故事、村落等新类型收藏。
func (a *API) handleToggleFavoriteItem(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	itemType := r.PathValue("type")
	if itemType != "product" && itemType != "story" && itemType != "village" {
		writeError(w, http.StatusBadRequest, "收藏类型不正确")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	favorited, err := a.Store.ToggleFavorite(ctx, key, itemType, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "更新收藏")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"favorited": favorited})
}

func (a *API) handleFavoriteItemStatus(w http.ResponseWriter, r *http.Request) {
	key := userKey(r)
	itemType := r.PathValue("type")
	if itemType != "product" && itemType != "story" && itemType != "village" {
		writeError(w, http.StatusBadRequest, "收藏类型不正确")
		return
	}
	if key == "" {
		writeJSON(w, http.StatusOK, map[string]bool{"favorited": false})
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	favorited, err := a.Store.FavoriteStatus(ctx, key, itemType, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取收藏状态")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"favorited": favorited})
}


func (a *API) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	profile, err := a.Store.GetUserProfile(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取成长账户")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (a *API) handleListMedals(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListMedals(ctx, key)
	if err != nil {
		handleStoreError(w, err, "读取徽章")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func (a *API) handleGetMedal(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	item, err := a.Store.GetMedal(ctx, key, r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "读取徽章详情")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) handleListPointLedger(w http.ResponseWriter, r *http.Request) {
	key := requireUserKey(w, r)
	if key == "" {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListPointLedger(ctx, key, 50)
	if err != nil {
		handleStoreError(w, err, "读取积分流水")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}
