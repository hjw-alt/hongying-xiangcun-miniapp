package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

func jsonList(items []string) string {
	if items == nil {
		items = []string{}
	}
	encoded, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func parseList(raw string) []string {
	items := []string{}
	if strings.TrimSpace(raw) == "" {
		return items
	}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}
	return items
}

func parseOrderItems(raw string) []OrderItem {
	items := []OrderItem{}
	if strings.TrimSpace(raw) == "" {
		return items
	}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []OrderItem{}
	}
	return items
}

func (s *MySQLStore) migrateRedTourism() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS banners (
			id VARCHAR(40) PRIMARY KEY,
			title VARCHAR(80) NOT NULL,
			subtitle VARCHAR(160) NOT NULL DEFAULT '',
			image_url VARCHAR(500) NOT NULL DEFAULT '',
			link_type VARCHAR(20) NOT NULL DEFAULT '',
			link_id VARCHAR(40) NOT NULL DEFAULT '',
			sort INT NOT NULL DEFAULT 0,
			active TINYINT(1) NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS villages (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(80) NOT NULL,
			region VARCHAR(80) NOT NULL DEFAULT '',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			history TEXT NOT NULL,
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			lat DECIMAL(10,6) NOT NULL DEFAULT 0,
			lng DECIMAL(10,6) NOT NULL DEFAULT 0,
			honors VARCHAR(600) NOT NULL DEFAULT '[]',
			themes VARCHAR(600) NOT NULL DEFAULT '[]',
			spot_count INT NOT NULL DEFAULT 0,
			story_count INT NOT NULL DEFAULT 0,
			product_count INT NOT NULL DEFAULT 0,
			sort INT NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS spots (
			id VARCHAR(40) PRIMARY KEY,
			village_id VARCHAR(40) NOT NULL,
			name VARCHAR(80) NOT NULL,
			category VARCHAR(30) NOT NULL DEFAULT '',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			image_url VARCHAR(500) NOT NULL DEFAULT '',
			audio_url VARCHAR(500) NOT NULL DEFAULT '',
			lat DECIMAL(10,6) NOT NULL DEFAULT 0,
			lng DECIMAL(10,6) NOT NULL DEFAULT 0,
			sort INT NOT NULL DEFAULT 0,
			INDEX idx_spots_village (village_id, sort)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS routes (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(80) NOT NULL,
			theme VARCHAR(40) NOT NULL DEFAULT '',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			days INT NOT NULL DEFAULT 1,
			distance_km DECIMAL(6,1) NOT NULL DEFAULT 0,
			spot_ids VARCHAR(600) NOT NULL DEFAULT '[]',
			spot_names VARCHAR(600) NOT NULL DEFAULT '[]',
			sort INT NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS stories (
			id VARCHAR(40) PRIMARY KEY,
			village_id VARCHAR(40) NOT NULL DEFAULT '',
			title VARCHAR(120) NOT NULL,
			category VARCHAR(30) NOT NULL DEFAULT '',
			era VARCHAR(40) NOT NULL DEFAULT '',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			content MEDIUMTEXT NOT NULL,
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			audio_url VARCHAR(500) NOT NULL DEFAULT '',
			ar_image VARCHAR(500) NOT NULL DEFAULT '',
			views INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_stories_village (village_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS farmers (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(80) NOT NULL,
			village_id VARCHAR(40) NOT NULL DEFAULT '',
			village_name VARCHAR(80) NOT NULL DEFAULT '',
			intro VARCHAR(500) NOT NULL DEFAULT '',
			avatar_url VARCHAR(500) NOT NULL DEFAULT '',
			project VARCHAR(200) NOT NULL DEFAULT '',
			support_target INT NOT NULL DEFAULT 0,
			support_count INT NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			category VARCHAR(30) NOT NULL DEFAULT '',
			price_cents INT NOT NULL DEFAULT 0,
			points_price INT NOT NULL DEFAULT 0,
			reward_points INT NOT NULL DEFAULT 0,
			reward_exp INT NOT NULL DEFAULT 0,
			original_cents INT NOT NULL DEFAULT 0,
			stock INT NOT NULL DEFAULT 0,
			sales INT NOT NULL DEFAULT 0,
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			images VARCHAR(1200) NOT NULL DEFAULT '[]',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			detail TEXT NOT NULL,
			story_card VARCHAR(1000) NOT NULL DEFAULT '',
			tags VARCHAR(400) NOT NULL DEFAULT '[]',
			customizable TINYINT(1) NOT NULL DEFAULT 0,
			farmer_id VARCHAR(40) NOT NULL DEFAULT '',
			farmer_name VARCHAR(80) NOT NULL DEFAULT '',
			village_id VARCHAR(40) NOT NULL DEFAULT '',
			status VARCHAR(20) NOT NULL DEFAULT '上架',
			sort INT NOT NULL DEFAULT 0,
			INDEX idx_products_category (category, status, sort)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(40) PRIMARY KEY,
			village_id VARCHAR(40) NOT NULL DEFAULT '',
			spot_id VARCHAR(40) NOT NULL DEFAULT '',
			title VARCHAR(120) NOT NULL,
			type VARCHAR(30) NOT NULL DEFAULT '',
			summary VARCHAR(500) NOT NULL DEFAULT '',
			points INT NOT NULL DEFAULT 0,
			reward_exp INT NOT NULL DEFAULT 0,
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			sort INT NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS task_records (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			task_id VARCHAR(40) NOT NULL,
			task_title VARCHAR(120) NOT NULL DEFAULT '',
			village_id VARCHAR(40) NOT NULL DEFAULT '',
			photo_url VARCHAR(500) NOT NULL DEFAULT '',
			note VARCHAR(300) NOT NULL DEFAULT '',
			points INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_task_records_user (user_key, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS orders (
			id CHAR(36) PRIMARY KEY,
			order_no VARCHAR(32) NOT NULL,
			user_key VARCHAR(64) NOT NULL,
			items TEXT NOT NULL,
			total_cents INT NOT NULL DEFAULT 0,
			status VARCHAR(20) NOT NULL DEFAULT '待发货',
			receiver VARCHAR(40) NOT NULL DEFAULT '',
			phone VARCHAR(20) NOT NULL DEFAULT '',
			address VARCHAR(200) NOT NULL DEFAULT '',
			remark VARCHAR(300) NOT NULL DEFAULT '',
			support_farm TINYINT(1) NOT NULL DEFAULT 1,
			payment_type VARCHAR(20) NOT NULL DEFAULT 'cash',
			points_spent INT NOT NULL DEFAULT 0,
			reward_points INT NOT NULL DEFAULT 0,
			reward_exp INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_orders_user (user_key, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS user_profiles (
			user_key VARCHAR(64) PRIMARY KEY,
			nickname VARCHAR(80) NOT NULL DEFAULT '湘村体验者',
			points INT NOT NULL DEFAULT 0,
			experience INT NOT NULL DEFAULT 0,
			level INT NOT NULL DEFAULT 2,
			coupons INT NOT NULL DEFAULT 0,
			balance_cents INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS point_ledger (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			change_amount INT NOT NULL,
			balance_after INT NOT NULL DEFAULT 0,
			event_type VARCHAR(40) NOT NULL DEFAULT '',
			event_id VARCHAR(80) NOT NULL DEFAULT '',
			description VARCHAR(200) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_point_event (user_key,event_type,event_id),
			INDEX idx_point_ledger_user (user_key,created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS medals (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(80) NOT NULL,
			subtitle VARCHAR(120) NOT NULL DEFAULT '',
			description VARCHAR(500) NOT NULL DEFAULT '',
			icon_url VARCHAR(500) NOT NULL DEFAULT '',
			cover_url VARCHAR(500) NOT NULL DEFAULT '',
			condition_type VARCHAR(40) NOT NULL DEFAULT '',
			condition_value INT NOT NULL DEFAULT 1,
			sort INT NOT NULL DEFAULT 0,
			active TINYINT(1) NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS user_medals (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			medal_id VARCHAR(40) NOT NULL,
			unlocked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_user_medal (user_key,medal_id),
			INDEX idx_user_medals_user (user_key,unlocked_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS point_rules (
			id VARCHAR(40) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			event_type VARCHAR(40) NOT NULL DEFAULT '',
			points INT NOT NULL DEFAULT 0,
			experience INT NOT NULL DEFAULT 0,
			description VARCHAR(300) NOT NULL DEFAULT '',
			active TINYINT(1) NOT NULL DEFAULT 1,
			sort INT NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, statement := range statements {
		if _, err := s.DB.Exec(statement); err != nil {
			return err
		}
	}
	for _, column := range []struct { table, name, definition string }{
		{"products", "points_price", "INT NOT NULL DEFAULT 0"},
		{"products", "reward_points", "INT NOT NULL DEFAULT 0"},
		{"products", "reward_exp", "INT NOT NULL DEFAULT 0"},
		{"tasks", "reward_exp", "INT NOT NULL DEFAULT 0"},
		{"orders", "payment_type", "VARCHAR(20) NOT NULL DEFAULT 'cash'"},
		{"orders", "points_spent", "INT NOT NULL DEFAULT 0"},
		{"orders", "reward_points", "INT NOT NULL DEFAULT 0"},
		{"orders", "reward_exp", "INT NOT NULL DEFAULT 0"},
	} {
		if err := s.ensureColumn(column.table, column.name, column.definition); err != nil {
			return err
		}
	}
	if _, err := s.DB.Exec(`UPDATE tasks SET reward_exp=points WHERE reward_exp=0`); err != nil {
		return err
	}
	if err := s.seedRedTourism(); err != nil {
		return err
	}
	return s.seedRewards()
}

func (s *MySQLStore) bumpVillageCount(ctx context.Context, table string, villageID string) error {
	if villageID == "" {
		return nil
	}
	column := map[string]string{
		"spots":    "spot_count",
		"stories":  "story_count",
		"products": "product_count",
	}[table]
	if column == "" {
		return nil
	}
	_, err := s.DB.ExecContext(ctx, "UPDATE villages SET "+column+"=(SELECT COUNT(1) FROM "+table+" WHERE village_id=?) WHERE id=?", villageID, villageID)
	return err
}

func (s *MySQLStore) syncVillageCounts(ctx context.Context) error {
	pairs := [][2]string{{"spots", "spot_count"}, {"stories", "story_count"}, {"products", "product_count"}}
	for _, pair := range pairs {
		statement := "UPDATE villages SET " + pair[1] + "=(SELECT COUNT(1) FROM " + pair[0] + " WHERE " + pair[0] + ".village_id=villages.id)"
		if _, err := s.DB.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLStore) GetVillage(ctx context.Context, id string) (Village, error) {
	var item Village
	var honors, themes string
	err := s.DB.QueryRowContext(ctx, `SELECT id,name,region,summary,history,cover_url,lat,lng,honors,themes,
		spot_count,story_count,product_count,sort FROM villages WHERE id=?`, id).
		Scan(&item.ID, &item.Name, &item.Region, &item.Summary, &item.History, &item.CoverURL, &item.Lat, &item.Lng,
			&honors, &themes, &item.SpotCount, &item.StoryCount, &item.ProductCount, &item.Sort)
	if errors.Is(err, sql.ErrNoRows) {
		return Village{}, ErrNotFound
	}
	if err != nil {
		return Village{}, err
	}
	item.Honors = parseList(honors)
	item.Themes = parseList(themes)
	return item, nil
}
