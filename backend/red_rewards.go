package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

type levelRule struct {
	Level int
	Name  string
	Exp   int
}

var levelRules = []levelRule{
	{Level: 2, Name: "初见", Exp: 0},
	{Level: 3, Name: "熟客", Exp: 30},
	{Level: 4, Name: "领路人", Exp: 100},
	{Level: 5, Name: "湘村挚友", Exp: 300},
	{Level: 6, Name: "红色合伙人", Exp: 800},
}

type medalSeed struct {
	id, name, subtitle, description, icon, cover, conditionType string
	conditionValue, sort                                        int
}

func (s *MySQLStore) seedRewards() error {
	medals := []medalSeed{
		{"medal-first-checkin", "红迹初见", "完成第一次研学打卡", "在湘潭红色村落留下你的第一枚红迹。", "/static/red/product-badge.png", "/static/red/banner-village.png", "first_checkin", 1, 1},
		{"medal-checkin-3", "红色行者", "完成 3 次研学打卡", "一步一个脚印，把故事走成自己的记忆。", "/static/red/product-badge.png", "/static/red/route-shaoshan.png", "checkin_count", 3, 2},
		{"medal-points-100", "心意满百", "累计获得 100 心意值", "每一次打卡与下单，都是给湘村的温柔支持。", "/static/red/product-badge.png", "/static/red/banner-market.png", "points_total", 100, 3},
		{"medal-first-order", "湘村带回家", "完成第一笔商城订单", "把一份湘村好物带回家，也把故事带回家。", "/static/red/product-badge.png", "/static/red/banner-market.png", "first_order", 1, 4},
		{"medal-order-3", "红文创收藏家", "完成 3 笔商城订单", "用每一件好物，收藏一段湘潭记忆。", "/static/red/product-badge.png", "/static/red/product-bag.png", "order_count", 3, 5},
		{"medal-exp-300", "湘村挚友", "成长经验达到 300", "持续探索红色村落，成为真正的湘村挚友。", "/static/red/product-badge.png", "/static/red/village-xingtang.png", "experience", 300, 6},
	}
	for _, item := range medals {
		if _, err := s.DB.Exec(`INSERT INTO medals
			(id,name,subtitle,description,icon_url,cover_url,condition_type,condition_value,sort,active)
			VALUES (?,?,?,?,?,?,?,?,?,1)
			ON DUPLICATE KEY UPDATE name=VALUES(name), subtitle=VALUES(subtitle), description=VALUES(description),
			icon_url=VALUES(icon_url), cover_url=VALUES(cover_url), condition_type=VALUES(condition_type),
			condition_value=VALUES(condition_value), sort=VALUES(sort), active=1`,
			item.id, item.name, item.subtitle, item.description, item.icon, item.cover,
			item.conditionType, item.conditionValue, item.sort); err != nil {
			return err
		}
	}
	rules := []struct {
		id, name, eventType, description string
		points, experience, sort         int
	}{
		{"rule-task-checkin", "研学任务打卡", "task_checkin", "完成一次有效研学任务后发放积分与经验，具体数值由任务配置决定。", 20, 20, 1},
		{"rule-cash-order", "现金购买商品", "cash_order", "现金订单完成创建后，按商品配置发放积分与经验。", 10, 15, 2},
		{"rule-points-order", "积分兑换商品", "points_order", "积分兑换不重复发放积分，但仍可获得商品配置的成长经验。", 0, 8, 3},
	}
	for _, item := range rules {
		if _, err := s.DB.Exec(`INSERT INTO point_rules
			(id,name,event_type,points,experience,description,active,sort) VALUES (?,?,?,?,?,?,1,?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), event_type=VALUES(event_type), points=VALUES(points),
			experience=VALUES(experience), description=VALUES(description), active=1, sort=VALUES(sort)`,
			item.id, item.name, item.eventType, item.points, item.experience, item.description, item.sort); err != nil {
			return err
		}
	}
	return nil
}

// pointRuleValues holds the active default reward for one event type.
// Item/task-level configuration remains authoritative; these values only fill
// in fields that were intentionally left at zero by the administrator.
type pointRuleValues struct {
	points     int
	experience int
}

func activePointRuleTx(ctx context.Context, tx *sql.Tx, eventType string) (pointRuleValues, error) {
	var rule pointRuleValues
	err := tx.QueryRowContext(ctx, `SELECT points,experience FROM point_rules
		WHERE event_type=? AND active=1 ORDER BY sort ASC,id ASC LIMIT 1`, eventType).
		Scan(&rule.points, &rule.experience)
	if errors.Is(err, sql.ErrNoRows) {
		return pointRuleValues{}, nil
	}
	if err != nil {
		return pointRuleValues{}, err
	}
	return rule, nil
}

func ensureUserProfile(ctx context.Context, db interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, userKey string) error {
	if userKey == "" {
		return errors.New("缺少用户标识")
	}
	_, err := db.ExecContext(ctx, `INSERT IGNORE INTO user_profiles
		(user_key,nickname,points,experience,level,coupons,balance_cents) VALUES (?, '湘村体验者', 0, 0, 2, 0, 0)`, userKey)
	return err
}

func levelForExperience(exp int) levelRule {
	result := levelRules[0]
	for _, rule := range levelRules {
		if exp >= rule.Exp {
			result = rule
		} else {
			break
		}
	}
	return result
}

func nextLevelRule(level int) *levelRule {
	for index, rule := range levelRules {
		if rule.Level == level && index+1 < len(levelRules) {
			next := levelRules[index+1]
			return &next
		}
	}
	return nil
}

func profileValues(experience int) (level int, name string, current, next, progress int) {
	if experience < 0 {
		experience = 0
	}
	currentRule := levelForExperience(experience)
	level, name, current = currentRule.Level, currentRule.Name, currentRule.Exp
	if nextRule := nextLevelRule(level); nextRule != nil {
		next = nextRule.Exp
		span := next - current
		if span > 0 {
			progress = int(math.Round(float64(experience-current) * 100 / float64(span)))
		}
		if progress < 0 {
			progress = 0
		}
		if progress > 100 {
			progress = 100
		}
	} else {
		next = current
		progress = 100
	}
	return
}

func medalCurrentValueTx(ctx context.Context, tx *sql.Tx, userKey string, conditionType string) (int, error) {
	var value int
	var query string
	switch conditionType {
	case "first_checkin", "checkin_count":
		query = `SELECT COUNT(1) FROM task_records WHERE user_key=?`
	case "points_total":
		query = `SELECT COALESCE(SUM(CASE WHEN change_amount > 0 THEN change_amount ELSE 0 END),0) FROM point_ledger WHERE user_key=?`
	case "first_order", "order_count":
		query = `SELECT COUNT(1) FROM orders WHERE user_key=? AND status <> '已取消'`
	case "experience":
		query = `SELECT COALESCE(experience,0) FROM user_profiles WHERE user_key=?`
	default:
		return 0, nil
	}
	if err := tx.QueryRowContext(ctx, query, userKey).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return value, nil
}

func medalTarget(medal Medal) int {
	if medal.ConditionValue < 1 {
		return 1
	}
	return medal.ConditionValue
}

func listMedalsTx(ctx context.Context, tx *sql.Tx, userKey string) ([]Medal, error) {
	rows, err := tx.QueryContext(ctx, `SELECT m.id,m.name,m.subtitle,m.description,m.icon_url,m.cover_url,
		m.condition_type,m.condition_value,m.sort,m.active,um.unlocked_at
		FROM medals m LEFT JOIN user_medals um ON um.medal_id=m.id AND um.user_key=?
		WHERE m.active=1 ORDER BY m.sort ASC, m.id ASC`, userKey)
	if err != nil {
		return nil, err
	}
	items := []Medal{}
	for rows.Next() {
		var item Medal
		var active int
		var unlockedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.Subtitle, &item.Description, &item.IconURL, &item.CoverURL,
			&item.ConditionType, &item.ConditionValue, &item.Sort, &active, &unlockedAt); err != nil {
			_ = rows.Close()
			return nil, err
		}
		item.Active = active == 1
		item.Target = medalTarget(item)
		if unlockedAt.Valid {
			when := unlockedAt.Time
			item.UnlockedAt = &when
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for index := range items {
		value, err := medalCurrentValueTx(ctx, tx, userKey, items[index].ConditionType)
		if err != nil {
			return nil, err
		}
		items[index].Progress = value
		if items[index].Progress > items[index].Target {
			items[index].Progress = items[index].Target
		}
		items[index].Unlocked = items[index].UnlockedAt != nil || value >= items[index].Target
	}
	return items, nil
}

func unlockMedalsTx(ctx context.Context, tx *sql.Tx, userKey string) ([]Medal, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id,condition_type,condition_value FROM medals WHERE active=1 ORDER BY sort ASC,id ASC`)
	if err != nil {
		return nil, err
	}
	type medalCondition struct {
		id            string
		conditionType string
		target        int
	}
	conditions := []medalCondition{}
	for rows.Next() {
		var item medalCondition
		if err := rows.Scan(&item.id, &item.conditionType, &item.target); err != nil {
			_ = rows.Close()
			return nil, err
		}
		if item.target < 1 {
			item.target = 1
		}
		conditions = append(conditions, item)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for _, item := range conditions {
		value, err := medalCurrentValueTx(ctx, tx, userKey, item.conditionType)
		if err != nil {
			return nil, err
		}
		if value >= item.target {
			if _, err := tx.ExecContext(ctx, `INSERT IGNORE INTO user_medals (id,user_key,medal_id) VALUES (?,?,?)`, newID(), userKey, item.id); err != nil {
				return nil, err
			}
		}
	}
	return listMedalsTx(ctx, tx, userKey)
}

func profileFromRow(ctx context.Context, tx *sql.Tx, userKey string, points, experience, coupons, balance int) (UserProfile, error) {
	level, name, current, next, progress := profileValues(experience)
	medals, err := listMedalsTx(ctx, tx, userKey)
	if err != nil {
		return UserProfile{}, err
	}
	unlocked := make([]Medal, 0, len(medals))
	for _, medal := range medals {
		if medal.Unlocked {
			unlocked = append(unlocked, medal)
		}
	}
	return UserProfile{
		UserKey: userKey, Nickname: "湘村体验者", Points: points, Experience: experience,
		Level: level, LevelName: name, CurrentLevelExp: current, NextLevelExp: next,
		Progress: progress, Coupons: coupons, BalanceCents: balance,
		MedalCount: len(unlocked), Medals: unlocked,
	}, nil
}

func profileTx(ctx context.Context, tx *sql.Tx, userKey string) (UserProfile, error) {
	if err := ensureUserProfile(ctx, tx, userKey); err != nil {
		return UserProfile{}, err
	}
	var points, experience, coupons, balance int
	var nickname string
	if err := tx.QueryRowContext(ctx, `SELECT nickname,points,experience,coupons,balance_cents FROM user_profiles WHERE user_key=?`, userKey).
		Scan(&nickname, &points, &experience, &coupons, &balance); err != nil {
		return UserProfile{}, err
	}
	profile, err := profileFromRow(ctx, tx, userKey, points, experience, coupons, balance)
	if err != nil {
		return UserProfile{}, err
	}
	profile.Nickname = nickname
	return profile, nil
}

func (s *MySQLStore) GetUserProfile(ctx context.Context, userKey string) (UserProfile, error) {
	if err := ensureUserProfile(ctx, s.DB, userKey); err != nil {
		return UserProfile{}, err
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return UserProfile{}, err
	}
	defer tx.Rollback()
	if _, err := unlockMedalsTx(ctx, tx, userKey); err != nil {
		return UserProfile{}, err
	}
	profile, err := profileTx(ctx, tx, userKey)
	if err != nil {
		return UserProfile{}, err
	}
	if err := tx.Commit(); err != nil {
		return UserProfile{}, err
	}
	return profile, nil
}

func (s *MySQLStore) ListMedals(ctx context.Context, userKey string) ([]Medal, error) {
	if err := ensureUserProfile(ctx, s.DB, userKey); err != nil {
		return nil, err
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := unlockMedalsTx(ctx, tx, userKey); err != nil {
		return nil, err
	}
	items, err := listMedalsTx(ctx, tx, userKey)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *MySQLStore) GetMedal(ctx context.Context, userKey, id string) (Medal, error) {
	items, err := s.ListMedals(ctx, userKey)
	if err != nil {
		return Medal{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return Medal{}, ErrNotFound
}

func (s *MySQLStore) ListPointLedger(ctx context.Context, userKey string, limit int) ([]PointLedger, error) {
	if err := ensureUserProfile(ctx, s.DB, userKey); err != nil {
		return nil, err
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT id,change_amount,balance_after,event_type,event_id,description,created_at
		FROM point_ledger WHERE user_key=? ORDER BY created_at DESC LIMIT `+fmt.Sprint(limit), userKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PointLedger{}
	for rows.Next() {
		var item PointLedger
		if err := rows.Scan(&item.ID, &item.ChangeAmount, &item.BalanceAfter, &item.EventType, &item.EventID, &item.Description, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func ensureUserProfileTx(ctx context.Context, tx *sql.Tx, userKey string) error {
	return ensureUserProfile(ctx, tx, userKey)
}

func applyPointChangeTx(ctx context.Context, tx *sql.Tx, userKey string, change int, eventType, eventID, description string) error {
	if change == 0 {
		return nil
	}
	if err := ensureUserProfileTx(ctx, tx, userKey); err != nil {
		return err
	}
	var existing int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM point_ledger WHERE user_key=? AND event_type=? AND event_id=?`, userKey, eventType, eventID).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}
	var balance int
	if err := tx.QueryRowContext(ctx, `SELECT points FROM user_profiles WHERE user_key=? FOR UPDATE`, userKey).Scan(&balance); err != nil {
		return err
	}
	if balance+change < 0 {
		return errors.New("积分不足")
	}
	balance += change
	if _, err := tx.ExecContext(ctx, `UPDATE user_profiles SET points=? WHERE user_key=?`, balance, userKey); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO point_ledger
		(id,user_key,change_amount,balance_after,event_type,event_id,description) VALUES (?,?,?,?,?,?,?)`,
		newID(), userKey, change, balance, eventType, eventID, description)
	return err
}

func addExperienceTx(ctx context.Context, tx *sql.Tx, userKey string, delta int) error {
	if delta <= 0 {
		return nil
	}
	if err := ensureUserProfileTx(ctx, tx, userKey); err != nil {
		return err
	}
	var experience int
	if err := tx.QueryRowContext(ctx, `SELECT experience FROM user_profiles WHERE user_key=? FOR UPDATE`, userKey).Scan(&experience); err != nil {
		return err
	}
	experience += delta
	level, _, _, _, _ := profileValues(experience)
	_, err := tx.ExecContext(ctx, `UPDATE user_profiles SET experience=?,level=? WHERE user_key=?`, experience, level, userKey)
	return err
}

func profileAfterTx(ctx context.Context, tx *sql.Tx, userKey string) (UserProfile, error) {
	if _, err := unlockMedalsTx(ctx, tx, userKey); err != nil {
		return UserProfile{}, err
	}
	return profileTx(ctx, tx, userKey)
}

func sortMedals(items []Medal) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Unlocked != items[j].Unlocked {
			return items[i].Unlocked
		}
		if items[i].Sort != items[j].Sort {
			return items[i].Sort < items[j].Sort
		}
		return items[i].ID < items[j].ID
	})
}

// compile-time check keeps the time package used by generated API payloads and documents timestamp semantics.
var _ = time.Time{}
