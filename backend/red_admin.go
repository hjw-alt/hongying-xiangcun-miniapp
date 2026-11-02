package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type adminField struct {
	Column string
	Kind   string
}

type adminResource struct {
	Name   string
	Table  string
	Fields []adminField
	Order  string
}

func adminResources() map[string]adminResource {
	text := func(columns ...string) []adminField {
		fields := []adminField{}
		for _, column := range columns {
			fields = append(fields, adminField{Column: column, Kind: "text"})
		}
		return fields
	}
	resources := map[string]adminResource{
		"banners": {Name: "banners", Table: "banners", Order: "sort ASC",
			Fields: append(text("id", "title", "subtitle", "image_url", "link_type", "link_id"),
				adminField{"sort", "number"}, adminField{"active", "bool"})},
		"villages": {Name: "villages", Table: "villages", Order: "sort ASC",
			Fields: append(text("id", "name", "region", "summary", "history", "cover_url"),
				adminField{"lat", "float"}, adminField{"lng", "float"}, adminField{"honors", "list"},
				adminField{"themes", "list"}, adminField{"sort", "number"})},
		"spots": {Name: "spots", Table: "spots", Order: "village_id ASC, sort ASC",
			Fields: append(text("id", "village_id", "name", "category", "summary", "image_url", "audio_url"),
				adminField{"lat", "float"}, adminField{"lng", "float"}, adminField{"sort", "number"})},
		"routes": {Name: "routes", Table: "routes", Order: "sort ASC",
			Fields: append(text("id", "name", "theme", "summary", "cover_url"),
				adminField{"days", "number"}, adminField{"distance_km", "float"}, adminField{"spot_ids", "list"},
				adminField{"spot_names", "list"}, adminField{"sort", "number"})},
		"stories": {Name: "stories", Table: "stories", Order: "views DESC",
			Fields: append(text("id", "village_id", "title", "category", "era", "summary", "content", "cover_url", "audio_url", "ar_image"),
				adminField{"views", "number"})},
		"products": {Name: "products", Table: "products", Order: "sort ASC",
			Fields: append(text("id", "name", "category", "cover_url", "summary", "detail", "story_card", "farmer_id", "farmer_name", "village_id", "status"),
				adminField{"price_cents", "number"}, adminField{"points_price", "number"}, adminField{"reward_points", "number"}, adminField{"reward_exp", "number"},
				adminField{"original_cents", "number"}, adminField{"stock", "number"}, adminField{"sales", "number"},
				adminField{"images", "list"}, adminField{"tags", "list"},
				adminField{"customizable", "bool"}, adminField{"sort", "number"})},
		"farmers": {Name: "farmers", Table: "farmers", Order: "support_count DESC",
			Fields: append(text("id", "name", "village_id", "village_name", "intro", "avatar_url", "project"),
				adminField{"support_target", "number"}, adminField{"support_count", "number"})},
		"tasks": {Name: "tasks", Table: "tasks", Order: "sort ASC",
			Fields: append(text("id", "village_id", "spot_id", "title", "type", "summary", "cover_url"),
				adminField{"points", "number"}, adminField{"reward_exp", "number"}, adminField{"sort", "number"})},
		"medals": {Name: "medals", Table: "medals", Order: "sort ASC",
			Fields: append(text("id", "name", "subtitle", "description", "icon_url", "cover_url", "condition_type"),
				adminField{"condition_value", "number"}, adminField{"sort", "number"}, adminField{"active", "bool"})},
		"point_rules": {Name: "point_rules", Table: "point_rules", Order: "sort ASC",
			Fields: append(text("id", "name", "event_type", "description"),
				adminField{"points", "number"}, adminField{"experience", "number"}, adminField{"active", "bool"}, adminField{"sort", "number"})},
	}
	return resources
}

func resourceByPath(name string) (adminResource, bool) {
	item, ok := adminResources()[name]
	return item, ok
}

func fieldValue(field adminField, raw any) (any, error) {
	switch field.Kind {
	case "number":
		return toInt(raw), nil
	case "float":
		return toFloat(raw), nil
	case "bool":
		return toBoolInt(raw), nil
	case "list":
		switch value := raw.(type) {
		case []any:
			items := []string{}
			for _, entry := range value {
				items = append(items, strings.TrimSpace(fmt.Sprintf("%v", entry)))
			}
			return jsonList(items), nil
		case string:
			trimmed := strings.TrimSpace(value)
			if strings.HasPrefix(trimmed, "[") {
				return trimmed, nil
			}
			parts := strings.FieldsFunc(trimmed, func(r rune) bool { return r == '\n' || r == ',' || r == '，' })
			items := []string{}
			for _, part := range parts {
				if item := strings.TrimSpace(part); item != "" {
					items = append(items, item)
				}
			}
			return jsonList(items), nil
		default:
			return "[]", nil
		}
	default:
		if raw == nil {
			return "", nil
		}
		return strings.TrimSpace(fmt.Sprintf("%v", raw)), nil
	}
}

func toInt(raw any) int {
	switch value := raw.(type) {
	case float64:
		return int(value)
	case int:
		return value
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func toFloat(raw any) float64 {
	switch value := raw.(type) {
	case float64:
		return value
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func toBoolInt(raw any) int {
	switch value := raw.(type) {
	case bool:
		if value {
			return 1
		}
		return 0
	case string:
		if value == "true" || value == "1" {
			return 1
		}
		return 0
	case float64:
		if value != 0 {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func decodeAdminBody(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	payload := map[string]any{}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxJSONSize))
	if err := decoder.Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "请求数据格式不正确")
		return nil, err
	}
	return payload, nil
}

func (a *API) handleAdminResourceList(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	resource, ok := resourceByPath(r.PathValue("resource"))
	if !ok {
		writeError(w, http.StatusNotFound, "未知的管理模块")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	rows, err := a.Store.DB.QueryContext(ctx, "SELECT * FROM "+resource.Table+" ORDER BY "+resource.Order)
	if err != nil {
		handleStoreError(w, err, "读取数据")
		return
	}
	defer rows.Close()
	items, err := rowsToMaps(rows)
	if err != nil {
		handleStoreError(w, err, "解析数据")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func rowsToMaps(rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	items := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for index := range values {
			pointers[index] = &values[index]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}
		item := map[string]any{}
		for index, column := range columns {
			value := values[index]
			if raw, ok := value.([]byte); ok {
				item[column] = string(raw)
			} else {
				item[column] = value
			}
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) handleAdminResourceCreate(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	resource, ok := resourceByPath(r.PathValue("resource"))
	if !ok {
		writeError(w, http.StatusNotFound, "未知的管理模块")
		return
	}
	payload, err := decodeAdminBody(w, r)
	if err != nil {
		return
	}
	columns := []string{}
	placeholders := []string{}
	values := []any{}
	for _, field := range resource.Fields {
		raw, exists := payload[field.Column]
		if !exists {
			if field.Column == "id" {
				raw = newID()
			} else {
				continue
			}
		}
		value, err := fieldValue(field, raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		columns = append(columns, field.Column)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}
	if len(columns) == 0 {
		writeError(w, http.StatusBadRequest, "没有可写入的字段")
		return
	}
	statement := "INSERT INTO " + resource.Table + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(placeholders, ",") + ")"
	ctx, cancel := requestContext(r)
	defer cancel()
	if _, err := a.Store.DB.ExecContext(ctx, statement, values...); err != nil {
		handleStoreError(w, err, "新增数据")
		return
	}
	_ = a.Store.syncVillageCounts(ctx)
	writeJSON(w, http.StatusCreated, map[string]any{"ok": true})
}

func (a *API) handleAdminResourceUpdate(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	resource, ok := resourceByPath(r.PathValue("resource"))
	if !ok {
		writeError(w, http.StatusNotFound, "未知的管理模块")
		return
	}
	payload, err := decodeAdminBody(w, r)
	if err != nil {
		return
	}
	assignments := []string{}
	values := []any{}
	for _, field := range resource.Fields {
		if field.Column == "id" {
			continue
		}
		raw, exists := payload[field.Column]
		if !exists {
			continue
		}
		value, err := fieldValue(field, raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		assignments = append(assignments, field.Column+"=?")
		values = append(values, value)
	}
	if len(assignments) == 0 {
		writeError(w, http.StatusBadRequest, "没有可更新的字段")
		return
	}
	values = append(values, r.PathValue("id"))
	ctx, cancel := requestContext(r)
	defer cancel()
	result, err := a.Store.DB.ExecContext(ctx,
		"UPDATE "+resource.Table+" SET "+strings.Join(assignments, ",")+" WHERE id=?", values...)
	if err != nil {
		handleStoreError(w, err, "更新数据")
		return
	}
	if err := requireOneRow(result); err != nil {
		handleStoreError(w, err, "更新数据")
		return
	}
	_ = a.Store.syncVillageCounts(ctx)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *API) handleAdminResourceDelete(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	resource, ok := resourceByPath(r.PathValue("resource"))
	if !ok {
		writeError(w, http.StatusNotFound, "未知的管理模块")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	result, err := a.Store.DB.ExecContext(ctx, "DELETE FROM "+resource.Table+" WHERE id=?", r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err, "删除数据")
		return
	}
	if err := requireOneRow(result); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "记录不存在")
			return
		}
		handleStoreError(w, err, "删除数据")
		return
	}
	_ = a.Store.syncVillageCounts(ctx)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *API) handleAdminRedSummary(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	summary, err := a.Store.RedSummary(ctx)
	if err != nil {
		handleStoreError(w, err, "读取后台统计")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (a *API) handleAdminOrders(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.Store.ListOrders(ctx, "", r.URL.Query().Get("status"))
	if err != nil {
		handleStoreError(w, err, "读取订单")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func (a *API) handleAdminOrderStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	var input struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	allowed := map[string]bool{"待发货": true, "已发货": true, "已完成": true, "已取消": true}
	if !allowed[input.Status] {
		writeError(w, http.StatusBadRequest, "订单状态不正确")
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.Store.UpdateOrderStatus(ctx, r.PathValue("id"), input.Status); err != nil {
		handleStoreError(w, err, "更新订单状态")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
