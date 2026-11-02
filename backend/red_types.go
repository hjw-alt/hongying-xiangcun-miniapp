package main

import "time"

// Banner 首页轮播。
type Banner struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"imageUrl"`
	LinkType string `json:"linkType"`
	LinkID   string `json:"linkId"`
	Sort     int    `json:"sort"`
	Active   bool   `json:"active"`
}

// Village 红色村落。
type Village struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Region       string   `json:"region"`
	Summary      string   `json:"summary"`
	History      string   `json:"history"`
	CoverURL     string   `json:"coverUrl"`
	Lat          float64  `json:"lat"`
	Lng          float64  `json:"lng"`
	Honors       []string `json:"honors"`
	Themes       []string `json:"themes"`
	SpotCount    int      `json:"spotCount"`
	StoryCount   int      `json:"storyCount"`
	ProductCount int      `json:"productCount"`
	Sort         int      `json:"sort"`
}

// Spot 村落内的红色景点。
type Spot struct {
	ID        string  `json:"id"`
	VillageID string  `json:"villageId"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Summary   string  `json:"summary"`
	ImageURL  string  `json:"imageUrl"`
	AudioURL  string  `json:"audioUrl"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Sort      int     `json:"sort"`
}

// Route 红色主题线路。
type Route struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Theme      string   `json:"theme"`
	Summary    string   `json:"summary"`
	CoverURL   string   `json:"coverUrl"`
	Days       int      `json:"days"`
	DistanceKM float64  `json:"distanceKm"`
	SpotIDs    []string `json:"spotIds"`
	SpotNames  []string `json:"spotNames"`
	Sort       int      `json:"sort"`
}

// Story 红迹数字故事。
type Story struct {
	ID        string    `json:"id"`
	VillageID string    `json:"villageId"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Era       string    `json:"era"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	CoverURL  string    `json:"coverUrl"`
	AudioURL  string    `json:"audioUrl"`
	ARImage   string    `json:"arImage"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"createdAt"`
}

// Product 文创与助农商品。
type Product struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	PriceCents   int      `json:"priceCents"`
	OriginalCent int      `json:"originalCents"`
	Stock        int      `json:"stock"`
	Sales        int      `json:"sales"`
	CoverURL     string   `json:"coverUrl"`
	Images       []string `json:"images"`
	Summary      string   `json:"summary"`
	Detail       string   `json:"detail"`
	StoryCard    string   `json:"storyCard"`
	Tags         []string `json:"tags"`
	Customizable bool     `json:"customizable"`
	PointsPrice  int      `json:"pointsPrice"`
	RewardPoints int      `json:"rewardPoints"`
	RewardExp    int      `json:"rewardExp"`
	FarmerID     string   `json:"farmerId"`
	FarmerName   string   `json:"farmerName"`
	VillageID    string   `json:"villageId"`
	Status       string   `json:"status"`
	Sort         int      `json:"sort"`
}

// Farmer 农户或合作社。
type Farmer struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	VillageID     string `json:"villageId"`
	VillageName   string `json:"villageName"`
	Intro         string `json:"intro"`
	AvatarURL     string `json:"avatarUrl"`
	Project       string `json:"project"`
	SupportTarget int    `json:"supportTarget"`
	SupportCount  int    `json:"supportCount"`
}

// Task 研学打卡任务。
type Task struct {
	ID        string `json:"id"`
	VillageID string `json:"villageId"`
	SpotID    string `json:"spotId"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Summary   string `json:"summary"`
	Points    int    `json:"points"`
	RewardExp int    `json:"rewardExp"`
	CoverURL  string `json:"coverUrl"`
	Sort      int    `json:"sort"`
}

// TaskRecord 打卡记录。
type TaskRecord struct {
	ID        string    `json:"id"`
	UserKey   string    `json:"-"`
	TaskID    string    `json:"taskId"`
	TaskTitle string    `json:"taskTitle"`
	VillageID string    `json:"villageId"`
	PhotoURL  string    `json:"photoUrl"`
	Note      string    `json:"note"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"createdAt"`
}

// OrderItem 订单明细。
type OrderItem struct {
	ProductID  string `json:"productId"`
	Name       string `json:"name"`
	PriceCents  int    `json:"priceCents"`
	PointsPrice int    `json:"pointsPrice,omitempty"`
	Quantity    int    `json:"quantity"`
	CustomText string `json:"customText"`
}

// Order 订单。
type Order struct {
	ID          string      `json:"id"`
	OrderNo     string      `json:"orderNo"`
	UserKey     string      `json:"-"`
	Items       []OrderItem `json:"items"`
	TotalCents  int         `json:"totalCents"`
	Status      string      `json:"status"`
	Receiver    string      `json:"receiver"`
	Phone       string      `json:"phone"`
	Address     string      `json:"address"`
	Remark      string      `json:"remark"`
	SupportFarm bool        `json:"supportFarm"`
	PaymentType string      `json:"paymentType"`
	PointsSpent int         `json:"pointsSpent"`
	RewardPoints int        `json:"rewardPoints"`
	RewardExp   int         `json:"rewardExp"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// UserProfile 用户成长账户。
type UserProfile struct {
	UserKey         string    `json:"userKey"`
	Nickname        string    `json:"nickname"`
	Points          int       `json:"points"`
	Experience      int       `json:"experience"`
	Level           int       `json:"level"`
	LevelName       string    `json:"levelName"`
	CurrentLevelExp int       `json:"currentLevelExp"`
	NextLevelExp    int       `json:"nextLevelExp"`
	Progress        int       `json:"progress"`
	Coupons         int       `json:"coupons"`
	BalanceCents    int       `json:"balanceCents"`
	MedalCount      int       `json:"medalCount"`
	Medals          []Medal   `json:"medals"`
}

// Medal 徽章及用户解锁状态。
type Medal struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Subtitle       string     `json:"subtitle"`
	Description    string     `json:"description"`
	IconURL        string     `json:"iconUrl"`
	CoverURL       string     `json:"coverUrl"`
	ConditionType  string     `json:"conditionType"`
	ConditionValue int        `json:"conditionValue"`
	Sort           int        `json:"sort"`
	Active         bool       `json:"active"`
	Unlocked       bool       `json:"unlocked"`
	UnlockedAt     *time.Time `json:"unlockedAt,omitempty"`
	Progress       int        `json:"progress"`
	Target         int        `json:"target"`
}

// PointLedger 积分流水。
type PointLedger struct {
	ID           string    `json:"id"`
	UserKey      string    `json:"-"`
	ChangeAmount int       `json:"changeAmount"`
	BalanceAfter int       `json:"balanceAfter"`
	EventType    string    `json:"eventType"`
	EventID      string    `json:"eventId"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
}
