package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *MySQLStore) ListBanners(ctx context.Context) ([]Banner, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,title,subtitle,image_url,link_type,link_id,sort,active
		FROM banners WHERE active=1 ORDER BY sort ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Banner{}
	for rows.Next() {
		var item Banner
		var active int
		if err := rows.Scan(&item.ID, &item.Title, &item.Subtitle, &item.ImageURL, &item.LinkType, &item.LinkID, &item.Sort, &active); err != nil {
			return nil, err
		}
		item.Active = active == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListVillages(ctx context.Context, keyword string) ([]Village, error) {
	query := `SELECT id,name,region,summary,history,cover_url,lat,lng,honors,themes,
		spot_count,story_count,product_count,sort FROM villages`
	args := []any{}
	if strings.TrimSpace(keyword) != "" {
		query += ` WHERE name LIKE ? OR region LIKE ? OR summary LIKE ?`
		like := "%" + strings.TrimSpace(keyword) + "%"
		args = append(args, like, like, like)
	}
	query += ` ORDER BY sort ASC`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Village{}
	for rows.Next() {
		var item Village
		var honors, themes string
		if err := rows.Scan(&item.ID, &item.Name, &item.Region, &item.Summary, &item.History, &item.CoverURL,
			&item.Lat, &item.Lng, &honors, &themes, &item.SpotCount, &item.StoryCount, &item.ProductCount, &item.Sort); err != nil {
			return nil, err
		}
		item.Honors = parseList(honors)
		item.Themes = parseList(themes)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListSpots(ctx context.Context, villageID string) ([]Spot, error) {
	query := `SELECT id,village_id,name,category,summary,image_url,audio_url,lat,lng,sort FROM spots`
	args := []any{}
	if villageID != "" {
		query += ` WHERE village_id=?`
		args = append(args, villageID)
	}
	query += ` ORDER BY village_id ASC, sort ASC`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Spot{}
	for rows.Next() {
		var item Spot
		if err := rows.Scan(&item.ID, &item.VillageID, &item.Name, &item.Category, &item.Summary,
			&item.ImageURL, &item.AudioURL, &item.Lat, &item.Lng, &item.Sort); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListRoutes(ctx context.Context) ([]Route, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,theme,summary,cover_url,days,distance_km,spot_ids,spot_names,sort
		FROM routes ORDER BY sort ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Route{}
	for rows.Next() {
		var item Route
		var spotIDs, spotNames string
		if err := rows.Scan(&item.ID, &item.Name, &item.Theme, &item.Summary, &item.CoverURL, &item.Days,
			&item.DistanceKM, &spotIDs, &spotNames, &item.Sort); err != nil {
			return nil, err
		}
		item.SpotIDs = parseList(spotIDs)
		item.SpotNames = parseList(spotNames)
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanStories(rows *sql.Rows) ([]Story, error) {
	items := []Story{}
	for rows.Next() {
		var item Story
		if err := rows.Scan(&item.ID, &item.VillageID, &item.Title, &item.Category, &item.Era, &item.Summary,
			&item.Content, &item.CoverURL, &item.AudioURL, &item.ARImage, &item.Views, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

const storyColumns = `id,village_id,title,category,era,summary,content,cover_url,audio_url,ar_image,views,created_at`

func (s *MySQLStore) ListStories(ctx context.Context, villageID, category, keyword string) ([]Story, error) {
	query := `SELECT ` + storyColumns + ` FROM stories WHERE 1=1`
	args := []any{}
	if villageID != "" {
		query += ` AND village_id=?`
		args = append(args, villageID)
	}
	if category != "" && category != "全部" {
		query += ` AND category=?`
		args = append(args, category)
	}
	if strings.TrimSpace(keyword) != "" {
		query += ` AND (title LIKE ? OR summary LIKE ?)`
		like := "%" + strings.TrimSpace(keyword) + "%"
		args = append(args, like, like)
	}
	query += ` ORDER BY views DESC, created_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStories(rows)
}

func (s *MySQLStore) GetStory(ctx context.Context, id string) (Story, error) {
	var item Story
	err := s.DB.QueryRowContext(ctx, `SELECT `+storyColumns+` FROM stories WHERE id=?`, id).
		Scan(&item.ID, &item.VillageID, &item.Title, &item.Category, &item.Era, &item.Summary, &item.Content,
			&item.CoverURL, &item.AudioURL, &item.ARImage, &item.Views, &item.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Story{}, ErrNotFound
	}
	if err != nil {
		return Story{}, err
	}
	_, _ = s.DB.ExecContext(ctx, `UPDATE stories SET views=views+1 WHERE id=?`, id)
	item.Views++
	return item, nil
}

const productColumns = `id,name,category,price_cents,points_price,reward_points,reward_exp,original_cents,stock,sales,cover_url,images,summary,detail,
	story_card,tags,customizable,farmer_id,farmer_name,village_id,status,sort`

func scanProducts(rows *sql.Rows) ([]Product, error) {
	items := []Product{}
	for rows.Next() {
		var item Product
		var images, tags string
		var custom int
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.PriceCents, &item.PointsPrice, &item.RewardPoints, &item.RewardExp, &item.OriginalCent, &item.Stock,
			&item.Sales, &item.CoverURL, &images, &item.Summary, &item.Detail, &item.StoryCard, &tags, &custom,
			&item.FarmerID, &item.FarmerName, &item.VillageID, &item.Status, &item.Sort); err != nil {
			return nil, err
		}
		item.Images = parseList(images)
		item.Tags = parseList(tags)
		item.Customizable = custom == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListProducts(ctx context.Context, category, keyword, sort, farmerID string) ([]Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE status='上架'`
	args := []any{}
	if category != "" && category != "全部" {
		query += ` AND category=?`
		args = append(args, category)
	}
	if farmerID != "" {
		query += ` AND farmer_id=?`
		args = append(args, farmerID)
	}
	if strings.TrimSpace(keyword) != "" {
		query += ` AND (name LIKE ? OR summary LIKE ?)`
		like := "%" + strings.TrimSpace(keyword) + "%"
		args = append(args, like, like)
	}
	switch sort {
	case "sales":
		query += ` ORDER BY sales DESC, sort ASC`
	case "price":
		query += ` ORDER BY price_cents ASC, sort ASC`
	default:
		query += ` ORDER BY sort ASC`
	}
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (s *MySQLStore) GetProduct(ctx context.Context, id string) (Product, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT `+productColumns+` FROM products WHERE id=?`, id)
	if err != nil {
		return Product{}, err
	}
	defer rows.Close()
	items, err := scanProducts(rows)
	if err != nil {
		return Product{}, err
	}
	if len(items) == 0 {
		return Product{}, ErrNotFound
	}
	return items[0], nil
}

func (s *MySQLStore) ListFarmers(ctx context.Context) ([]Farmer, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,village_id,village_name,intro,avatar_url,project,
		support_target,support_count FROM farmers ORDER BY support_count DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Farmer{}
	for rows.Next() {
		var item Farmer
		if err := rows.Scan(&item.ID, &item.Name, &item.VillageID, &item.VillageName, &item.Intro,
			&item.AvatarURL, &item.Project, &item.SupportTarget, &item.SupportCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListTasks(ctx context.Context, villageID string) ([]Task, error) {
	query := `SELECT id,village_id,spot_id,title,type,summary,points,reward_exp,cover_url,sort FROM tasks`
	args := []any{}
	if villageID != "" {
		query += ` WHERE village_id=?`
		args = append(args, villageID)
	}
	query += ` ORDER BY sort ASC`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Task{}
	for rows.Next() {
		var item Task
		if err := rows.Scan(&item.ID, &item.VillageID, &item.SpotID, &item.Title, &item.Type, &item.Summary,
			&item.Points, &item.RewardExp, &item.CoverURL, &item.Sort); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) CreateTaskRecord(ctx context.Context, item TaskRecord) (TaskRecord, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return TaskRecord{}, err
	}
	defer tx.Rollback()
	if err := ensureUserProfileTx(ctx, tx, item.UserKey); err != nil {
		return TaskRecord{}, err
	}
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM task_records WHERE user_key=? AND task_id=?`,
		item.UserKey, item.TaskID).Scan(&count); err != nil {
		return TaskRecord{}, err
	}
	if count > 0 {
		return TaskRecord{}, errors.New("该任务已经打卡")
	}
	var task Task
	err = tx.QueryRowContext(ctx, `SELECT id,village_id,title,points,reward_exp FROM tasks WHERE id=?`, item.TaskID).
		Scan(&task.ID, &task.VillageID, &task.Title, &task.Points, &task.RewardExp)
	if errors.Is(err, sql.ErrNoRows) {
		return TaskRecord{}, ErrNotFound
	}
	if err != nil {
		return TaskRecord{}, err
	}
	rule, err := activePointRuleTx(ctx, tx, "task_checkin")
	if err != nil {
		return TaskRecord{}, err
	}
	if task.Points <= 0 && rule.points > 0 {
		task.Points = rule.points
	}
	item.ID = newID()
	item.TaskTitle = task.Title
	item.VillageID = task.VillageID
	item.Points = task.Points
	item.CreatedAt = time.Now()
	if _, err := tx.ExecContext(ctx, `INSERT INTO task_records
		(id,user_key,task_id,task_title,village_id,photo_url,note,points) VALUES (?,?,?,?,?,?,?,?)`,
		item.ID, item.UserKey, item.TaskID, item.TaskTitle, item.VillageID, item.PhotoURL, item.Note, item.Points); err != nil {
		return TaskRecord{}, err
	}
	if err := applyPointChangeTx(ctx, tx, item.UserKey, item.Points, "task_checkin", item.ID, fmt.Sprintf("完成研学任务：%s", item.TaskTitle)); err != nil {
		return TaskRecord{}, err
	}
	rewardExp := task.RewardExp
	if rewardExp <= 0 {
		if rule.experience > 0 {
			rewardExp = rule.experience
		} else {
			rewardExp = task.Points
		}
	}
	if err := addExperienceTx(ctx, tx, item.UserKey, rewardExp); err != nil {
		return TaskRecord{}, err
	}
	if _, err := unlockMedalsTx(ctx, tx, item.UserKey); err != nil {
		return TaskRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return TaskRecord{}, err
	}
	return item, nil
}

func (s *MySQLStore) ListTaskRecords(ctx context.Context, userKey string) ([]TaskRecord, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,task_id,task_title,village_id,photo_url,note,points,created_at
		FROM task_records WHERE user_key=? ORDER BY created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TaskRecord{}
	for rows.Next() {
		var item TaskRecord
		if err := rows.Scan(&item.ID, &item.TaskID, &item.TaskTitle, &item.VillageID, &item.PhotoURL,
			&item.Note, &item.Points, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
