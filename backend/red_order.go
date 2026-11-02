package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func formatOrderNo(now time.Time) string {
	return fmt.Sprintf("HX%s%04d", now.Format("20060102150405"), now.Nanosecond()/1000000)
}

// CreateOrder 创建订单，同时扣减库存、累加销量与助农数据，并完成积分/经验结算。
func (s *MySQLStore) CreateOrder(ctx context.Context, item Order) (Order, error) {
	if len(item.Items) == 0 {
		return Order{}, errors.New("订单中没有商品")
	}
	if item.PaymentType == "" {
		item.PaymentType = "cash"
	}
	if item.PaymentType != "cash" && item.PaymentType != "points" {
		return Order{}, errors.New("支付方式不正确")
	}
	now := time.Now()
	item.ID = newID()
	item.OrderNo = formatOrderNo(now)
	item.CreatedAt = now

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()
	if err := ensureUserProfileTx(ctx, tx, item.UserKey); err != nil {
		return Order{}, err
	}
	cashRule, err := activePointRuleTx(ctx, tx, "cash_order")
	if err != nil {
		return Order{}, err
	}
	pointsRule, err := activePointRuleTx(ctx, tx, "points_order")
	if err != nil {
		return Order{}, err
	}

	total := 0
	pointsSpent := 0
	rewardPoints := 0
	rewardExp := 0
	farmers := map[string]bool{}
	for index := range item.Items {
		row := item.Items[index]
		var price, pointsPrice, productRewardPoints, productRewardExp, stock int
		var name, farmerID string
		err := tx.QueryRowContext(ctx, `SELECT name,price_cents,points_price,reward_points,reward_exp,stock,farmer_id
			FROM products WHERE id=? AND status='上架' FOR UPDATE`, row.ProductID).
			Scan(&name, &price, &pointsPrice, &productRewardPoints, &productRewardExp, &stock, &farmerID)
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, errors.New("商品不存在或已下架")
		}
		if err != nil {
			return Order{}, err
		}
		quantity := row.Quantity
		if quantity < 1 {
			quantity = 1
		}
		if stock < quantity {
			return Order{}, fmt.Errorf("「%s」库存不足", name)
		}
		if item.PaymentType == "points" && pointsPrice <= 0 {
			return Order{}, fmt.Errorf("「%s」暂不支持积分兑换", name)
		}
		if item.PaymentType == "cash" && productRewardPoints <= 0 && cashRule.points > 0 {
			productRewardPoints = cashRule.points
		}
		item.Items[index].Name = name
		item.Items[index].PriceCents = price
		item.Items[index].PointsPrice = pointsPrice
		item.Items[index].Quantity = quantity
		total += price * quantity
		pointsSpent += pointsPrice * quantity
		rewardPoints += productRewardPoints * quantity
		productExp := productRewardExp
		if productExp <= 0 {
			if item.PaymentType == "points" && pointsRule.experience > 0 {
				productExp = pointsRule.experience
			} else if item.PaymentType == "cash" && cashRule.experience > 0 {
				productExp = cashRule.experience
			} else {
				productExp = productRewardPoints
			}
		}
		rewardExp += productExp * quantity
		if _, err := tx.ExecContext(ctx, `UPDATE products SET stock=stock-?, sales=sales+? WHERE id=?`,
			quantity, quantity, row.ProductID); err != nil {
			return Order{}, err
		}
		if farmerID != "" {
			farmers[farmerID] = true
		}
	}
	for farmerID := range farmers {
		if _, err := tx.ExecContext(ctx, `UPDATE farmers SET support_count=support_count+1 WHERE id=?`, farmerID); err != nil {
			return Order{}, err
		}
	}

	item.TotalCents = total
	item.PointsSpent = 0
	item.RewardPoints = 0
	item.RewardExp = rewardExp
	item.Status = "待发货"
	if item.PaymentType == "points" {
		item.PointsSpent = pointsSpent
		if err := applyPointChangeTx(ctx, tx, item.UserKey, -pointsSpent, "points_order", item.ID,
			fmt.Sprintf("积分兑换订单 %s", item.OrderNo)); err != nil {
			return Order{}, err
		}
	} else {
		item.RewardPoints = rewardPoints
		if rewardPoints > 0 {
			if err := applyPointChangeTx(ctx, tx, item.UserKey, rewardPoints, "cash_order", item.ID,
				fmt.Sprintf("现金购买奖励：%s", item.OrderNo)); err != nil {
				return Order{}, err
			}
		}
	}
	if err := addExperienceTx(ctx, tx, item.UserKey, rewardExp); err != nil {
		return Order{}, err
	}

	encoded, err := json.Marshal(item.Items)
	if err != nil {
		return Order{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO orders
		(id,order_no,user_key,items,total_cents,status,receiver,phone,address,remark,support_farm,payment_type,points_spent,reward_points,reward_exp)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		item.ID, item.OrderNo, item.UserKey, string(encoded), item.TotalCents, item.Status,
		item.Receiver, item.Phone, item.Address, item.Remark, boolToInt(item.SupportFarm), item.PaymentType,
		item.PointsSpent, item.RewardPoints, item.RewardExp); err != nil {
		return Order{}, err
	}
	if _, err := unlockMedalsTx(ctx, tx, item.UserKey); err != nil {
		return Order{}, err
	}
	if err := tx.Commit(); err != nil {
		return Order{}, err
	}
	return item, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func scanOrders(rows *sql.Rows) ([]Order, error) {
	items := []Order{}
	for rows.Next() {
		var item Order
		var encoded string
		var support int
		if err := rows.Scan(&item.ID, &item.OrderNo, &item.UserKey, &encoded, &item.TotalCents, &item.Status,
			&item.Receiver, &item.Phone, &item.Address, &item.Remark, &support, &item.PaymentType,
			&item.PointsSpent, &item.RewardPoints, &item.RewardExp, &item.CreatedAt); err != nil {
			return nil, err
		}
		item.Items = parseOrderItems(encoded)
		item.SupportFarm = support == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

const orderColumns = `id,order_no,user_key,items,total_cents,status,receiver,phone,address,remark,support_farm,payment_type,points_spent,reward_points,reward_exp,created_at`

func (s *MySQLStore) ListOrders(ctx context.Context, userKey, status string) ([]Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders`
	args := []any{}
	conditions := []string{}
	if userKey != "" {
		conditions = append(conditions, "user_key=?")
		args = append(args, userKey)
	}
	if status != "" {
		conditions = append(conditions, "status=?")
		args = append(args, status)
	}
	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for _, extra := range conditions[1:] {
			query += " AND " + extra
		}
	}
	query += ` ORDER BY created_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (s *MySQLStore) UpdateOrderStatus(ctx context.Context, id, status string) error {
	result, err := s.DB.ExecContext(ctx, `UPDATE orders SET status=? WHERE id=?`, status, id)
	if err != nil {
		return err
	}
	return requireOneRow(result)
}

// FarmSummary 助农专区统计数据。
func (s *MySQLStore) FarmSummary(ctx context.Context) (map[string]int, error) {
	summary := map[string]int{}
	queries := []struct {
		key   string
		query string
	}{
		{"farmers", "SELECT COUNT(1) FROM farmers"},
		{"products", "SELECT COUNT(1) FROM products WHERE category='湘村风物助农' AND status='上架'"},
		{"orders", "SELECT COUNT(1) FROM orders WHERE support_farm=1"},
		{"supportCount", "SELECT COALESCE(SUM(support_count),0) FROM farmers"},
	}
	for _, item := range queries {
		var value int
		if err := s.DB.QueryRowContext(ctx, item.query).Scan(&value); err != nil {
			return nil, err
		}
		summary[item.key] = value
	}
	return summary, nil
}

// RedSummary 后台首页统计。
func (s *MySQLStore) RedSummary(ctx context.Context) (map[string]int, error) {
	summary := map[string]int{}
	queries := []struct {
		key   string
		query string
	}{
		{"villages", "SELECT COUNT(1) FROM villages"},
		{"spots", "SELECT COUNT(1) FROM spots"},
		{"routes", "SELECT COUNT(1) FROM routes"},
		{"stories", "SELECT COUNT(1) FROM stories"},
		{"products", "SELECT COUNT(1) FROM products"},
		{"farmers", "SELECT COUNT(1) FROM farmers"},
		{"tasks", "SELECT COUNT(1) FROM tasks"},
		{"orders", "SELECT COUNT(1) FROM orders"},
		{"taskRecords", "SELECT COUNT(1) FROM task_records"},
		{"banners", "SELECT COUNT(1) FROM banners"},
		{"profiles", "SELECT COUNT(1) FROM user_profiles"},
		{"medals", "SELECT COUNT(1) FROM medals WHERE active=1"},
		{"pointLedger", "SELECT COUNT(1) FROM point_ledger"},
	}
	for _, item := range queries {
		var value int
		if err := s.DB.QueryRowContext(ctx, item.query).Scan(&value); err != nil {
			return nil, err
		}
		summary[item.key] = value
	}
	return summary, nil
}
