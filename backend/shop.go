package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// Address 收货地址。
type Address struct {
	ID        string    `json:"id"`
	UserKey   string    `json:"-"`
	Receiver  string    `json:"receiver"`
	Phone     string    `json:"phone"`
	Region    string    `json:"region"`
	Detail    string    `json:"detail"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
}

// CartItem 购物车条目，附带商品快照信息。
type CartItem struct {
	ID         string    `json:"id"`
	UserKey    string    `json:"-"`
	ProductID  string    `json:"productId"`
	Quantity   int       `json:"quantity"`
	CustomText string    `json:"customText"`
	CreatedAt  time.Time `json:"createdAt"`
	Product    *Product  `json:"product,omitempty"`
}

func (s *MySQLStore) migrateShop() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS addresses (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			receiver VARCHAR(40) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			region VARCHAR(120) NOT NULL DEFAULT '',
			detail VARCHAR(200) NOT NULL,
			is_default TINYINT(1) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_addresses_user (user_key, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS cart_items (
			id CHAR(36) PRIMARY KEY,
			user_key VARCHAR(64) NOT NULL,
			product_id VARCHAR(40) NOT NULL,
			quantity INT NOT NULL DEFAULT 1,
			custom_text VARCHAR(200) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_cart_user (user_key, created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, statement := range statements {
		if _, err := s.DB.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLStore) ListAddresses(ctx context.Context, userKey string) ([]Address, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,receiver,phone,region,detail,is_default,created_at
		FROM addresses WHERE user_key=? ORDER BY is_default DESC, created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Address{}
	for rows.Next() {
		var item Address
		var isDefault int
		if err := rows.Scan(&item.ID, &item.Receiver, &item.Phone, &item.Region, &item.Detail, &isDefault, &item.CreatedAt); err != nil {
			return nil, err
		}
		item.IsDefault = isDefault == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) GetAddress(ctx context.Context, userKey, id string) (Address, error) {
	var item Address
	var isDefault int
	err := s.DB.QueryRowContext(ctx, `SELECT id,receiver,phone,region,detail,is_default,created_at
		FROM addresses WHERE id=? AND user_key=?`, id, userKey).
		Scan(&item.ID, &item.Receiver, &item.Phone, &item.Region, &item.Detail, &isDefault, &item.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Address{}, ErrNotFound
	}
	item.IsDefault = isDefault == 1
	return item, err
}

func (s *MySQLStore) CreateAddress(ctx context.Context, item Address) (Address, error) {
	var count int
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM addresses WHERE user_key=?`, item.UserKey).Scan(&count); err != nil {
		return Address{}, err
	}
	if count == 0 {
		item.IsDefault = true
	}
	item.ID = newID()
	item.CreatedAt = time.Now()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return Address{}, err
	}
	defer tx.Rollback()
	if item.IsDefault {
		if _, err := tx.ExecContext(ctx, `UPDATE addresses SET is_default=0 WHERE user_key=?`, item.UserKey); err != nil {
			return Address{}, err
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO addresses
		(id,user_key,receiver,phone,region,detail,is_default) VALUES (?,?,?,?,?,?,?)`,
		item.ID, item.UserKey, item.Receiver, item.Phone, item.Region, item.Detail, boolToInt(item.IsDefault)); err != nil {
		return Address{}, err
	}
	if err := tx.Commit(); err != nil {
		return Address{}, err
	}
	return item, nil
}

func (s *MySQLStore) UpdateAddress(ctx context.Context, item Address) (Address, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return Address{}, err
	}
	defer tx.Rollback()
	if item.IsDefault {
		if _, err := tx.ExecContext(ctx, `UPDATE addresses SET is_default=0 WHERE user_key=?`, item.UserKey); err != nil {
			return Address{}, err
		}
	}
	result, err := tx.ExecContext(ctx, `UPDATE addresses SET receiver=?, phone=?, region=?, detail=?, is_default=?
		WHERE id=? AND user_key=?`,
		item.Receiver, item.Phone, item.Region, item.Detail, boolToInt(item.IsDefault), item.ID, item.UserKey)
	if err != nil {
		return Address{}, err
	}
	if err := requireOneRow(result); err != nil {
		return Address{}, err
	}
	if err := tx.Commit(); err != nil {
		return Address{}, err
	}
	return s.GetAddress(ctx, item.UserKey, item.ID)
}

func (s *MySQLStore) DeleteAddress(ctx context.Context, userKey, id string) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM addresses WHERE id=? AND user_key=?`, id, userKey)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) SetDefaultAddress(ctx context.Context, userKey, id string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `UPDATE addresses SET is_default=0 WHERE user_key=?`, userKey); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE addresses SET is_default=1 WHERE id=? AND user_key=?`, id, userKey)
	if err != nil {
		return err
	}
	if err := requireOneRow(result); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *MySQLStore) ListCartItems(ctx context.Context, userKey string) ([]CartItem, error) {
	// JOIN 查询时商品字段需要加表前缀，避免与 cart_items 的同名列冲突
	const joinedProductColumns = `p.id,p.name,p.category,p.price_cents,p.points_price,p.reward_points,p.reward_exp,p.original_cents,p.stock,p.sales,
		p.cover_url,p.images,p.summary,p.detail,p.story_card,p.tags,p.customizable,p.farmer_id,p.farmer_name,
		p.village_id,p.status,p.sort`
	rows, err := s.DB.QueryContext(ctx, `SELECT c.id,c.product_id,c.quantity,c.custom_text,c.created_at,
		`+joinedProductColumns+`
		FROM cart_items c LEFT JOIN products p ON p.id=c.product_id
		WHERE c.user_key=? ORDER BY c.created_at DESC`, userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CartItem{}
	for rows.Next() {
		var item CartItem
		var product Product
		var images, tags string
		var custom int
		var productID sql.NullString
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.CustomText, &item.CreatedAt,
			&productID, &product.Name, &product.Category, &product.PriceCents, &product.PointsPrice, &product.RewardPoints, &product.RewardExp, &product.OriginalCent, &product.Stock,
			&product.Sales, &product.CoverURL, &images, &product.Summary, &product.Detail, &product.StoryCard, &tags,
			&custom, &product.FarmerID, &product.FarmerName, &product.VillageID, &product.Status, &product.Sort); err != nil {
			return nil, err
		}
		if productID.Valid {
			product.ID = productID.String
			product.Images = parseList(images)
			product.Tags = parseList(tags)
			product.Customizable = custom == 1
			item.Product = &product
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MySQLStore) AddCartItem(ctx context.Context, item CartItem) (CartItem, error) {
	var existingID string
	var existingQty int
	err := s.DB.QueryRowContext(ctx, `SELECT id,quantity FROM cart_items
		WHERE user_key=? AND product_id=? AND custom_text=?`, item.UserKey, item.ProductID, item.CustomText).
		Scan(&existingID, &existingQty)
	if err == nil {
		quantity := existingQty + item.Quantity
		if quantity > 99 {
			quantity = 99
		}
		if _, err := s.DB.ExecContext(ctx, `UPDATE cart_items SET quantity=? WHERE id=?`, quantity, existingID); err != nil {
			return CartItem{}, err
		}
		item.ID = existingID
		item.Quantity = quantity
		return item, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return CartItem{}, err
	}
	var stock int
	if err := s.DB.QueryRowContext(ctx, `SELECT stock FROM products WHERE id=? AND status='上架'`, item.ProductID).Scan(&stock); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CartItem{}, ErrNotFound
		}
		return CartItem{}, err
	}
	if item.Quantity < 1 {
		item.Quantity = 1
	}
	if stock <= 0 {
		return CartItem{}, errors.New("该商品暂时缺货")
	}
	item.ID = newID()
	item.CreatedAt = time.Now()
	_, err = s.DB.ExecContext(ctx, `INSERT INTO cart_items (id,user_key,product_id,quantity,custom_text)
		VALUES (?,?,?,?,?)`, item.ID, item.UserKey, item.ProductID, item.Quantity, item.CustomText)
	return item, err
}

func (s *MySQLStore) UpdateCartItem(ctx context.Context, userKey, id string, quantity int, customText string) error {
	if quantity < 1 {
		quantity = 1
	}
	if quantity > 99 {
		quantity = 99
	}
	result, err := s.DB.ExecContext(ctx, `UPDATE cart_items SET quantity=?, custom_text=? WHERE id=? AND user_key=?`,
		quantity, strings.TrimSpace(customText), id, userKey)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) DeleteCartItem(ctx context.Context, userKey, id string) error {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM cart_items WHERE id=? AND user_key=?`, id, userKey)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

func (s *MySQLStore) DeleteCartItems(ctx context.Context, userKey string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		if _, err := s.DB.ExecContext(ctx, `DELETE FROM cart_items WHERE id=? AND user_key=?`, id, userKey); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLStore) CountCartItems(ctx context.Context, userKey string) (int, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity),0) FROM cart_items WHERE user_key=?`, userKey).Scan(&count)
	return count, err
}
