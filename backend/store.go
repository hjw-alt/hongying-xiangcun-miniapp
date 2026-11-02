package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

var ErrNotFound = errors.New("记录不存在")

var databaseNamePattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

type MySQLStore struct {
	DB *sql.DB
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(6)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	store := &MySQLStore{DB: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func CreateDatabase(dsn string) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}
	databaseName := cfg.DBName
	if !databaseNamePattern.MatchString(cfg.DBName) {
		return fmt.Errorf("数据库名只能使用字母、数字或下划线")
	}
	cfg.DBName = ""
	server, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	defer server.Close()
	_, err = server.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", databaseName))
	return err
}

func (s *MySQLStore) migrate() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS artworks (
			id CHAR(36) PRIMARY KEY,
			title VARCHAR(80) NOT NULL,
			author VARCHAR(80) NOT NULL,
			theme VARCHAR(20) NOT NULL,
			description VARCHAR(500) NOT NULL DEFAULT '',
			tags VARCHAR(200) NOT NULL DEFAULT '',
			image_url VARCHAR(500) NOT NULL DEFAULT '/static/seed/placeholder.png',
			likes INT UNSIGNED NOT NULL DEFAULT 0,
			owner_key VARCHAR(64) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_theme_created (theme, created_at),
			INDEX idx_likes (likes)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS comments (
			id CHAR(36) PRIMARY KEY,
			artwork_id CHAR(36) NOT NULL,
			nickname VARCHAR(20) NOT NULL DEFAULT '匿名同学',
			content VARCHAR(200) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_comments_artwork FOREIGN KEY (artwork_id) REFERENCES artworks(id) ON DELETE CASCADE,
			INDEX idx_comments_artwork (artwork_id, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`ALTER TABLE artworks MODIFY COLUMN theme VARCHAR(20) NOT NULL`,
		`UPDATE artworks SET theme = '纸影人物' WHERE theme NOT IN ('纸影人物', '湘潭印象', '戏台故事')`,
		`INSERT INTO artworks (id, title, author, theme, description, tags, image_url, likes)
		SELECT 'seed-shadow-001', '穆桂英影偶', '吴影记', '纸影人物',
			'以湘潭纸影的七分脸和镂空纹样为基础，用水彩重绘穆桂英挂帅的英气。', '穆桂英,镂空,七分脸', '/static/seed/art-puppet.png', 32
		WHERE NOT EXISTS (SELECT 1 FROM artworks WHERE id = 'seed-shadow-001')`,
		`INSERT INTO artworks (id, title, author, theme, description, tags, image_url, likes)
		SELECT 'seed-shadow-002', '窑湾戏台', '陈星野', '湘潭印象',
			'窑湾老街的青瓦石巷外，叠加一场纸影戏开演前的暖黄灯光。', '窑湾,戏台,夜景', '/static/seed/art-town.png', 24
		WHERE NOT EXISTS (SELECT 1 FROM artworks WHERE id = 'seed-shadow-002')`,
		`INSERT INTO artworks (id, title, author, theme, description, tags, image_url, likes)
		SELECT 'seed-shadow-003', '一担戏箱进村小', '苏禾', '戏台故事',
			'传承人挑着戏箱走进乡村小学，孩子们在幕布后第一次举起自己的纸影。', '传承人,乡村美育,课堂', '/static/seed/art-class.png', 41
		WHERE NOT EXISTS (SELECT 1 FROM artworks WHERE id = 'seed-shadow-003')`,
		`UPDATE artworks SET image_url='/static/seed/art-puppet.png' WHERE id='seed-shadow-001' AND image_url LIKE '/static/seed/muguiying%'`,
		`UPDATE artworks SET image_url='/static/seed/art-town.png' WHERE id='seed-shadow-002' AND image_url LIKE '/static/seed/yaowan%'`,
		`UPDATE artworks SET image_url='/static/seed/art-class.png' WHERE id='seed-shadow-003' AND image_url LIKE '/static/seed/xitong%'`,
		`INSERT INTO comments (id, artwork_id, nickname, content)
		SELECT 'seed-comment-001', 'seed-shadow-001', '湘潭小满',
			'镂空纹样处理得很细，保留纸影脆感又有绘画语言。'
		WHERE NOT EXISTS (SELECT 1 FROM comments WHERE id = 'seed-comment-001')`,
		`INSERT INTO comments (id, artwork_id, nickname, content)
		SELECT 'seed-comment-002', 'seed-shadow-002', '湘江边',
			'窑湾夜色和戏台灯光的关系处理得很动人。'
		WHERE NOT EXISTS (SELECT 1 FROM comments WHERE id = 'seed-comment-002')`,
	}
	for index, statement := range statements {
		if _, err := s.DB.Exec(statement); err != nil {
			return fmt.Errorf("base migration statement %d: %w", index+1, err)
		}
	}
	exists, err := s.columnExists("artworks", "owner_key")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.DB.Exec("ALTER TABLE artworks ADD COLUMN owner_key VARCHAR(64) NOT NULL DEFAULT '' AFTER likes"); err != nil {
			return err
		}
	}
	if err := s.migrateContent(); err != nil {
		return fmt.Errorf("migrateContent: %w", err)
	}
	if err := s.migrateShadowPlays(); err != nil {
		return fmt.Errorf("migrateShadowPlays: %w", err)
	}
	if err := s.migrateRedTourism(); err != nil {
		return fmt.Errorf("migrateRedTourism: %w", err)
	}
	if err := s.migrateShop(); err != nil {
		return fmt.Errorf("migrateShop: %w", err)
	}
	return nil
}

func (s *MySQLStore) columnExists(tableName, columnName string) (bool, error) {
	var count int
	err := s.DB.QueryRow(`SELECT COUNT(1) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?`,
		tableName, columnName).Scan(&count)
	return count > 0, err
}

func (s *MySQLStore) ensureColumn(tableName, columnName, definition string) error {
	exists, err := s.columnExists(tableName, columnName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = s.DB.Exec("ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + definition)
	return err
}

func (s *MySQLStore) ListArtworks(ctx context.Context, keyword, theme, sort, ownerKey string) ([]Artwork, error) {
	query := `SELECT id, title, author, theme, description, tags, image_url, likes, owner_key, created_at, updated_at FROM artworks WHERE 1=1`
	args := []any{}
	if ownerKey != "" {
		query += " AND owner_key = ?"
		args = append(args, ownerKey)
	}
	if theme != "" && theme != "全部" {
		query += " AND theme = ?"
		args = append(args, theme)
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query += " AND (title LIKE ? OR author LIKE ? OR description LIKE ? OR tags LIKE ?)"
		args = append(args, like, like, like, like)
	}
	if sort == "likes" {
		query += " ORDER BY likes DESC, created_at DESC"
	} else {
		query += " ORDER BY created_at DESC"
	}
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArtworks(rows)
}

func scanArtworks(rows *sql.Rows) ([]Artwork, error) {
	items := []Artwork{}
	for rows.Next() {
		var item Artwork
		var tags string
		if err := rows.Scan(&item.ID, &item.Title, &item.Author, &item.Theme, &item.Description,
			&tags, &item.ImageURL, &item.Likes, &item.OwnerKey, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Tags = splitTags(tags)
		items = append(items, item)
	}
	return items, rows.Err()
}

func splitTags(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

func (s *MySQLStore) GetArtwork(ctx context.Context, id string) (Artwork, error) {
	items, err := s.ListArtworksByID(ctx, id)
	if err != nil {
		return Artwork{}, err
	}
	if len(items) == 0 {
		return Artwork{}, ErrNotFound
	}
	return items[0], nil
}

func (s *MySQLStore) ListArtworksByID(ctx context.Context, id string) ([]Artwork, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, title, author, theme, description, tags, image_url, likes, owner_key, created_at, updated_at
		FROM artworks WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArtworks(rows)
}

func (s *MySQLStore) CreateArtwork(ctx context.Context, item Artwork) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO artworks
		(id, title, author, theme, description, tags, image_url, likes, owner_key)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?)`,
		item.ID, item.Title, item.Author, item.Theme, item.Description, strings.Join(item.Tags, ","), item.ImageURL, item.OwnerKey)
	return err
}

func (s *MySQLStore) UpdateArtwork(ctx context.Context, item Artwork) error {
	result, err := s.DB.ExecContext(ctx, `UPDATE artworks SET
		title = ?, author = ?, theme = ?, description = ?, tags = ?, image_url = ?
		WHERE id = ?`,
		item.Title, item.Author, item.Theme, item.Description, strings.Join(item.Tags, ","), item.ImageURL, item.ID)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) DeleteArtwork(ctx context.Context, id string) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM artworks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) LikeArtwork(ctx context.Context, id string) error {
	result, err := s.DB.ExecContext(ctx, "UPDATE artworks SET likes = likes + 1 WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) CreateComment(ctx context.Context, item Comment) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO comments (id, artwork_id, nickname, content) VALUES (?, ?, ?, ?)`,
		item.ID, item.ArtworkID, item.Nickname, item.Content)
	return err
}

func (s *MySQLStore) ListComments(ctx context.Context, artworkID string) ([]Comment, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, artwork_id, nickname, content, created_at
		FROM comments WHERE artwork_id = ? ORDER BY created_at DESC`, artworkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Comment{}
	for rows.Next() {
		var item Comment
		if err := rows.Scan(&item.ID, &item.ArtworkID, &item.Nickname, &item.Content, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) DeleteComment(ctx context.Context, artworkID, commentID string) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM comments WHERE artwork_id = ? AND id = ?", artworkID, commentID)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) Stats(ctx context.Context) (Stats, error) {
	var result Stats
	err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(likes), 0),
		COALESCE(SUM(theme = '纸影人物'), 0), COALESCE(SUM(theme = '湘潭印象'), 0),
		COALESCE(SUM(theme = '戏台故事'), 0)
		FROM artworks`).Scan(&result.ArtworkTotal, &result.LikeTotal, &result.PuppetTotal, &result.ImpressTotal, &result.StoryTotal)
	return result, err
}

func requireOneRow(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
