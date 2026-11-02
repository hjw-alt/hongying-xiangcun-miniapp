package main

import "time"

type Artwork struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Theme       string    `json:"theme"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	ImageURL    string    `json:"imageUrl"`
	Likes       int       `json:"likes"`
	OwnerKey    string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        string    `json:"id"`
	ArtworkID string    `json:"artworkId"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Stats struct {
	ArtworkTotal int64 `json:"artworkTotal"`
	LikeTotal    int64 `json:"likeTotal"`
	PuppetTotal  int64 `json:"puppetTotal"`
	ImpressTotal int64 `json:"impressTotal"`
	StoryTotal   int64 `json:"storyTotal"`
}
