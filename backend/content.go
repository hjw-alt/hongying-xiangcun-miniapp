package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Clip struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"videoUrl"`
	CoverURL    string    `json:"coverUrl"`
	Uploader    string    `json:"uploader"`
	Likes       int       `json:"likes"`
	OwnerKey    string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Activity struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Venue       string    `json:"venue"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	PriceCents  int       `json:"priceCents"`
	Capacity    int       `json:"capacity"`
	Remaining   int       `json:"remaining"`
	CoverURL    string    `json:"coverUrl"`
	Content     string    `json:"content"`
	Highlights  []string  `json:"highlights"`
	Notice      string    `json:"notice"`
}

type Booking struct {
	ID            string    `json:"id"`
	OwnerKey      string    `json:"-"`
	ActivityID    string    `json:"activityId"`
	ActivityTitle string    `json:"activityTitle"`
	Venue         string    `json:"venue"`
	StartTime     time.Time `json:"startTime"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Participants  int       `json:"participants"`
	TotalFeeCents int       `json:"totalFeeCents"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}

const defaultCover = "/static/seed/placeholder.png"

func (s *MySQLStore) migrateContent() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS clips (
			id CHAR(36) PRIMARY KEY,
			title VARCHAR(80) NOT NULL,
			description VARCHAR(500) NOT NULL DEFAULT '',
			video_url VARCHAR(500) NOT NULL,
			cover_url VARCHAR(500) NOT NULL DEFAULT '/static/seed/placeholder.png',
			uploader VARCHAR(80) NOT NULL DEFAULT '匿名同学',
			likes INT UNSIGNED NOT NULL DEFAULT 0,
			owner_key VARCHAR(64) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_clips_created (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS activities (
			id VARCHAR(40) PRIMARY KEY,
			title VARCHAR(80) NOT NULL,
			category VARCHAR(20) NOT NULL,
			venue VARCHAR(120) NOT NULL,
			description VARCHAR(500) NOT NULL DEFAULT '',
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			price_cents INT NOT NULL DEFAULT 0,
			capacity INT NOT NULL DEFAULT 0,
			remaining INT NOT NULL DEFAULT 0,
			cover_url VARCHAR(500) NOT NULL DEFAULT '/static/seed/placeholder.png',
			content TEXT NOT NULL,
			highlights VARCHAR(600) NOT NULL DEFAULT '[]',
			notice VARCHAR(600) NOT NULL DEFAULT ''
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			activity_id VARCHAR(40) NOT NULL,
			name VARCHAR(20) NOT NULL,
			phone VARCHAR(11) NOT NULL,
			participants INT NOT NULL DEFAULT 1,
			total_fee_cents INT NOT NULL DEFAULT 0,
			status VARCHAR(20) NOT NULL DEFAULT '已预约',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_booking_activity FOREIGN KEY (activity_id) REFERENCES activities(id),
			INDEX idx_bookings_user (user_key, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS favorites (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			item_type VARCHAR(10) NOT NULL,
			item_id VARCHAR(40) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_favorite (user_key, item_type, item_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`INSERT INTO activities (id, title, category, venue, description, start_time, end_time, price_cents, capacity, remaining, cover_url, content, highlights, notice)
		SELECT 'activity-workshop', '简易纸影人DIY工作坊', 'DIY体验', '湘潭市非遗展示中心',
			'从画稿、剪刻到装杆，完成一个可操纵的简易纸影人。', '2026-09-20 14:00:00', '2026-09-20 16:30:00', 6800, 24, 24, '/static/seed/activity-workshop.png',
			'先看一段《穆桂英挂帅》选段，认识湘潭纸影的七分脸与镂空纹样；随后在传承人指导下画稿、剪刻、上色、刷桐油、装竹签，最后走到幕布后，操纵自己做的影偶亮一次相。',
			'["传承人现场示范操纵手法","亲手完成一个可操纵纸影人","作品可带走并参加线下展示","零基础可参与，含全部材料"]',
			'建议 12 岁以上参加，12 岁以下需家长陪同；材料含颜料与刻刀，现场统一发放并指导使用；请提前 10 分钟到场签到；活动开始前 24 小时可取消预约。'
		WHERE NOT EXISTS (SELECT 1 FROM activities WHERE id='activity-workshop')`,
		`UPDATE activities SET cover_url='/static/seed/activity-workshop.png'
		WHERE id='activity-workshop' AND cover_url LIKE '/static/seed/muguiying%'`,
		`INSERT INTO activities (id, title, category, venue, description, start_time, end_time, price_cents, capacity, remaining, cover_url, content, highlights, notice)
		SELECT 'activity-campus', '纸影戏进校园研学课', '进校园', '湖南科技大学美术学院',
			'面向美术专业学生的纸影临摹与舞台色彩公开课。', '2026-09-25 09:30:00', '2026-09-25 11:30:00', 0, 60, 60, '/static/seed/activity-campus.png',
			'从湘潭纸影的造型规律讲到舞台色彩：拆解七分脸、镂空与透光三者之间的关系，现场拆开一件真实影偶看内部结构，再用水彩或数字绘画完成一次影偶再设计练习。',
			'["非遗传承人与专业教师双主讲","现场拆解真实影偶结构","领取临摹线稿与配色卡","成果可用于课程作业与参赛素材"]',
			'凭学生证入场；请自带速写本与常用绘画工具；课堂内请勿使用闪光灯拍摄；名额有限，请按预约时间到场。'
		WHERE NOT EXISTS (SELECT 1 FROM activities WHERE id='activity-campus')`,
		`UPDATE activities SET cover_url='/static/seed/activity-campus.png'
		WHERE id='activity-campus' AND cover_url LIKE '/static/seed/xitong%'`,
		`INSERT INTO activities (id, title, category, venue, description, start_time, end_time, price_cents, capacity, remaining, cover_url, content, highlights, notice)
		SELECT 'activity-night', '湘江畔纸影夜展演', '展演活动', '窑湾历史文化街区戏台',
			'《穆桂英挂帅》《打猎回书》选段展演与幕后操纵体验。', '2026-09-28 19:30:00', '2026-09-28 21:00:00', 4900, 80, 80, '/static/seed/activity-night.png',
			'在窑湾老戏台看一场露天纸影戏：《穆桂英挂帅》《打猎回书》两出选段连演。演出结束后开放幕后，观众可以上台近距离看影偶、灯箱和操纵杆，听传承人讲一段旧时戏班的故事。',
			'["湘江夜色下的露天戏台","两出经典选段连演","演出后可上台体验操纵","允许拍照，适合结伴同行"]',
			'露天场地，雨天顺延并以短信通知为准；建议提前 20 分钟入场；现场不设固定座位，请听从工作人员引导；儿童需家长陪同。'
		WHERE NOT EXISTS (SELECT 1 FROM activities WHERE id='activity-night')`,
		`UPDATE activities SET cover_url='/static/seed/activity-night.png'
		WHERE id='activity-night' AND cover_url LIKE '/static/seed/yaowan%'`,
		`UPDATE activities SET
			content = CASE id
				WHEN 'activity-workshop' THEN '先看一段《穆桂英挂帅》选段，认识湘潭纸影的七分脸与镂空纹样；随后在传承人指导下画稿、剪刻、上色、刷桐油、装竹签，最后走到幕布后，操纵自己做的影偶亮一次相。'
				WHEN 'activity-campus' THEN '从湘潭纸影的造型规律讲到舞台色彩：拆解七分脸、镂空与透光三者之间的关系，现场拆开一件真实影偶看内部结构，再用水彩或数字绘画完成一次影偶再设计练习。'
				WHEN 'activity-night' THEN '在窑湾老戏台看一场露天纸影戏：《穆桂英挂帅》《打猎回书》两出选段连演。演出结束后开放幕后，观众可以上台近距离看影偶、灯箱和操纵杆，听传承人讲一段旧时戏班的故事。'
				ELSE content END,
			highlights = CASE id
				WHEN 'activity-workshop' THEN '["传承人现场示范操纵手法","亲手完成一个可操纵纸影人","作品可带走并参加线下展示","零基础可参与，含全部材料"]'
				WHEN 'activity-campus' THEN '["非遗传承人与专业教师双主讲","现场拆解真实影偶结构","领取临摹线稿与配色卡","成果可用于课程作业与参赛素材"]'
				WHEN 'activity-night' THEN '["湘江夜色下的露天戏台","两出经典选段连演","演出后可上台体验操纵","允许拍照，适合结伴同行"]'
				ELSE highlights END,
			notice = CASE id
				WHEN 'activity-workshop' THEN '建议 12 岁以上参加，12 岁以下需家长陪同；材料含颜料与刻刀，现场统一发放并指导使用；请提前 10 分钟到场签到；活动开始前 24 小时可取消预约。'
				WHEN 'activity-campus' THEN '凭学生证入场；请自带速写本与常用绘画工具；课堂内请勿使用闪光灯拍摄；名额有限，请按预约时间到场。'
				WHEN 'activity-night' THEN '露天场地，雨天顺延并以短信通知为准；建议提前 20 分钟入场；现场不设固定座位，请听从工作人员引导；儿童需家长陪同。'
				ELSE notice END
		WHERE id IN ('activity-workshop','activity-campus','activity-night')
			AND (content IS NULL OR content = '')`,
	}
	activityColumns := []struct {
		name      string
		statement string
	}{
		{"content", "ALTER TABLE activities ADD COLUMN content TEXT NOT NULL AFTER cover_url"},
		{"highlights", "ALTER TABLE activities ADD COLUMN highlights VARCHAR(600) NOT NULL DEFAULT '[]' AFTER content"},
		{"notice", "ALTER TABLE activities ADD COLUMN notice VARCHAR(600) NOT NULL DEFAULT '' AFTER highlights"},
	}
	for _, column := range activityColumns {
		exists, err := s.columnExists("activities", column.name)
		if err != nil {
			return err
		}
		if !exists {
			if _, err := s.DB.Exec(column.statement); err != nil {
				return err
			}
		}
	}
	for index, statement := range statements {
		if _, err := s.DB.Exec(statement); err != nil {
			return fmt.Errorf("content migration statement %d: %w", index+1, err)
		}
	}
	return nil
}

func scanClips(rows *sql.Rows) ([]Clip, error) {
	items := []Clip{}
	for rows.Next() {
		var item Clip
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.VideoURL, &item.CoverURL,
			&item.Uploader, &item.Likes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) ListClips(ctx context.Context) ([]Clip, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, title, description, video_url, cover_url, uploader, likes, created_at
		FROM clips ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanClips(rows)
}

func (s *MySQLStore) CreateClip(ctx context.Context, item Clip) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO clips
		(id, title, description, video_url, cover_url, uploader, owner_key)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Title, item.Description, item.VideoURL, item.CoverURL, item.Uploader, item.OwnerKey)
	return err
}

func (s *MySQLStore) UpdateClip(ctx context.Context, item Clip) error {
	result, err := s.DB.ExecContext(ctx, `UPDATE clips SET title=?, description=?, video_url=?, cover_url=?, uploader=?
		WHERE id=?`,
		item.Title, item.Description, item.VideoURL, item.CoverURL, item.Uploader, item.ID)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) DeleteClip(ctx context.Context, id string) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM clips WHERE id=?", id)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) LikeClip(ctx context.Context, id string) (Clip, error) {
	result, err := s.DB.ExecContext(ctx, "UPDATE clips SET likes=likes+1 WHERE id=?", id)
	if err != nil {
		return Clip{}, err
	}
	if err := requireOneRow(result); err != nil {
		return Clip{}, err
	}
	var item Clip
	err = s.DB.QueryRowContext(ctx, `SELECT id, title, description, video_url, cover_url, uploader, likes, created_at
		FROM clips WHERE id=?`, id).Scan(&item.ID, &item.Title, &item.Description, &item.VideoURL,
		&item.CoverURL, &item.Uploader, &item.Likes, &item.CreatedAt)
	return item, err
}

func (s *MySQLStore) favoriteExists(ctx context.Context, userKey, itemType, itemID string) (bool, error) {
	var exists bool
	err := s.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_key=? AND item_type=? AND item_id=?)`,
		userKey, itemType, itemID).Scan(&exists)
	return exists, err
}

func (s *MySQLStore) ToggleFavorite(ctx context.Context, userKey, itemType, itemID string) (bool, error) {
	var exists bool
	var err error
	// 收藏对象所在的业务表
	tables := map[string]string{
		"artwork": "artworks",
		"clip":    "clips",
		"product": "products",
		"story":   "stories",
		"village": "villages",
	}
	table, ok := tables[itemType]
	if !ok {
		return false, errors.New("不支持的收藏类型")
	}
	var count int
	err = s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM "+table+" WHERE id=?", itemID).Scan(&count)
	if err == nil && count == 0 {
		return false, ErrNotFound
	}
	exists, err = s.favoriteExists(ctx, userKey, itemType, itemID)
	if err != nil {
		return false, err
	}
	if exists {
		_, err = s.DB.ExecContext(ctx, `DELETE FROM favorites WHERE user_key=? AND item_type=? AND item_id=?`, userKey, itemType, itemID)
		return false, err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO favorites (user_key, item_type, item_id) VALUES (?, ?, ?)`, userKey, itemType, itemID)
	return true, err
}

func (s *MySQLStore) FavoriteStatus(ctx context.Context, userKey, itemType, itemID string) (bool, error) {
	return s.favoriteExists(ctx, userKey, itemType, itemID)
}

func (s *MySQLStore) ListFavoriteArtworks(ctx context.Context, userKey string) ([]Artwork, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT a.id,a.title,a.author,a.theme,a.description,a.tags,a.image_url,a.likes,a.owner_key,a.created_at,a.updated_at
		FROM favorites f JOIN artworks a ON f.item_id=a.id
		WHERE f.user_key=? AND f.item_type='artwork' ORDER BY f.created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArtworks(rows)
}

func (s *MySQLStore) ListFavoriteClips(ctx context.Context, userKey string) ([]Clip, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT c.id,c.title,c.description,c.video_url,c.cover_url,c.uploader,c.likes,c.created_at
		FROM favorites f JOIN clips c ON f.item_id=c.id
		WHERE f.user_key=? AND f.item_type='clip' ORDER BY f.created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanClips(rows)
}

func (s *MySQLStore) ListActivities(ctx context.Context) ([]Activity, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,title,category,venue,description,start_time,end_time,price_cents,capacity,remaining,cover_url,content,highlights,notice
		FROM activities ORDER BY start_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Activity{}
	for rows.Next() {
		var item Activity
		var highlights string
		if err := rows.Scan(&item.ID, &item.Title, &item.Category, &item.Venue, &item.Description, &item.StartTime,
			&item.EndTime, &item.PriceCents, &item.Capacity, &item.Remaining, &item.CoverURL,
			&item.Content, &highlights, &item.Notice); err != nil {
			return nil, err
		}
		item.Highlights = parseHighlights(highlights)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) GetActivity(ctx context.Context, id string) (Activity, error) {
	var item Activity
	var highlights string
	err := s.DB.QueryRowContext(ctx, `SELECT id,title,category,venue,description,start_time,end_time,price_cents,capacity,remaining,cover_url,content,highlights,notice
		FROM activities WHERE id=?`, id).Scan(&item.ID, &item.Title, &item.Category, &item.Venue, &item.Description,
		&item.StartTime, &item.EndTime, &item.PriceCents, &item.Capacity, &item.Remaining, &item.CoverURL,
		&item.Content, &highlights, &item.Notice)
	if errors.Is(err, sql.ErrNoRows) {
		return Activity{}, ErrNotFound
	}
	if err == nil {
		item.Highlights = parseHighlights(highlights)
	}
	return item, err
}

func parseHighlights(raw string) []string {
	items := []string{}
	if strings.TrimSpace(raw) == "" {
		return items
	}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}
	return items
}

func (s *MySQLStore) CreateActivity(ctx context.Context, item Activity, highlights string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO activities
		(id, title, category, venue, description, start_time, end_time, price_cents, capacity, remaining, cover_url, content, highlights, notice)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Title, item.Category, item.Venue, item.Description, item.StartTime, item.EndTime,
		item.PriceCents, item.Capacity, item.Remaining, item.CoverURL, item.Content, highlights, item.Notice)
	return err
}

func (s *MySQLStore) UpdateActivity(ctx context.Context, item Activity, highlights string) error {
	result, err := s.DB.ExecContext(ctx, `UPDATE activities SET
		title=?, category=?, venue=?, description=?, start_time=?, end_time=?, price_cents=?,
		capacity=?, remaining=?, cover_url=?, content=?, highlights=?, notice=?
		WHERE id=?`,
		item.Title, item.Category, item.Venue, item.Description, item.StartTime, item.EndTime,
		item.PriceCents, item.Capacity, item.Remaining, item.CoverURL, item.Content, highlights,
		item.Notice, item.ID)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) DeleteActivity(ctx context.Context, id string) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM activities WHERE id=?", id)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) AdminSummary(ctx context.Context) (map[string]int, error) {
	counts := map[string]int{}
	queries := []struct {
		key   string
		query string
	}{
		{"artworks", "SELECT COUNT(1) FROM artworks"},
		{"clips", "SELECT COUNT(1) FROM clips"},
		{"activities", "SELECT COUNT(1) FROM activities"},
		{"bookings", "SELECT COUNT(1) FROM bookings"},
		{"shadowPlays", "SELECT COUNT(1) FROM shadow_plays"},
	}
	for _, item := range queries {
		var count int
		if err := s.DB.QueryRowContext(ctx, item.query).Scan(&count); err != nil {
			return nil, err
		}
		counts[item.key] = count
	}
	return counts, nil
}

func (s *MySQLStore) CreateBooking(ctx context.Context, item Booking) (Booking, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return Booking{}, err
	}
	defer tx.Rollback()

	var activity Activity
	err = tx.QueryRowContext(ctx, `SELECT id,title,venue,start_time,price_cents,remaining FROM activities WHERE id=? FOR UPDATE`,
		item.ActivityID).Scan(&activity.ID, &activity.Title, &activity.Venue, &activity.StartTime, &activity.PriceCents, &activity.Remaining)
	if errors.Is(err, sql.ErrNoRows) {
		return Booking{}, ErrNotFound
	}
	if err != nil {
		return Booking{}, err
	}
	if activity.Remaining < item.Participants {
		return Booking{}, errors.New("名额不足")
	}
	item.ID = newID()
	item.ActivityTitle = activity.Title
	item.Venue = activity.Venue
	item.StartTime = activity.StartTime
	item.TotalFeeCents = activity.PriceCents * item.Participants
	item.Status = "已预约"
	_, err = tx.ExecContext(ctx, `INSERT INTO bookings
		(id,user_key,activity_id,name,phone,participants,total_fee_cents,status)
		VALUES (?,?,?,?,?,?,?,?)`,
		item.ID, item.OwnerKey, item.ActivityID, item.Name, item.Phone, item.Participants, item.TotalFeeCents, item.Status)
	if err != nil {
		return Booking{}, err
	}
	_, err = tx.ExecContext(ctx, "UPDATE activities SET remaining=remaining-? WHERE id=?", item.Participants, item.ActivityID)
	if err != nil {
		return Booking{}, err
	}
	if err = tx.Commit(); err != nil {
		return Booking{}, err
	}
	return item, nil
}

func (s *MySQLStore) ListBookings(ctx context.Context, userKey string) ([]Booking, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT b.id,b.activity_id,a.title,a.venue,a.start_time,b.name,b.phone,
		b.participants,b.total_fee_cents,b.status,b.created_at
		FROM bookings b JOIN activities a ON b.activity_id=a.id
		WHERE b.user_key=? ORDER BY b.created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Booking{}
	for rows.Next() {
		var item Booking
		if err := rows.Scan(&item.ID, &item.ActivityID, &item.ActivityTitle, &item.Venue, &item.StartTime, &item.Name,
			&item.Phone, &item.Participants, &item.TotalFeeCents, &item.Status, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) CancelBooking(ctx context.Context, userKey, bookingID string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var activityID string
	var participants int
	var status string
	err = tx.QueryRowContext(ctx, `SELECT activity_id,participants,status FROM bookings WHERE id=? AND user_key=? FOR UPDATE`,
		bookingID, userKey).Scan(&activityID, &participants, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if status != "已预约" {
		return errors.New("该订单已取消")
	}
	if _, err := tx.ExecContext(ctx, "UPDATE bookings SET status='已取消' WHERE id=?", bookingID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE activities SET remaining=remaining+? WHERE id=?", participants, activityID); err != nil {
		return err
	}
	return tx.Commit()
}

func cleanText(value string, maxLength int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > maxLength {
		return string(runes[:maxLength])
	}
	return value
}
