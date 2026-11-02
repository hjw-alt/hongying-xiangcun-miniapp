package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root@tcp(127.0.0.1:3306)/shadowart_xiangtan?charset=utf8mb4&parseTime=true&loc=Local"
	}

	if err := CreateDatabase(dsn); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	store, err := NewMySQLStore(dsn)
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	defer store.DB.Close()

	api := &API{Store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", api.handleHealth)
	mux.HandleFunc("GET /api/banners", api.handleListBanners)
	mux.HandleFunc("GET /api/villages", api.handleListVillages)
	mux.HandleFunc("GET /api/villages/{id}", api.handleGetVillage)
	mux.HandleFunc("GET /api/spots", api.handleListSpots)
	mux.HandleFunc("GET /api/routes", api.handleListRoutes)
	mux.HandleFunc("GET /api/stories", api.handleListStories)
	mux.HandleFunc("GET /api/stories/{id}", api.handleGetStory)
	mux.HandleFunc("GET /api/products", api.handleListProducts)
	mux.HandleFunc("GET /api/products/{id}", api.handleGetProduct)
	mux.HandleFunc("GET /api/farmers", api.handleListFarmers)
	mux.HandleFunc("GET /api/farm-summary", api.handleFarmSummary)
	mux.HandleFunc("GET /api/tasks", api.handleListTasks)
	mux.HandleFunc("POST /api/tasks/{id}/checkin", api.handleCheckinTask)
	mux.HandleFunc("GET /api/task-records", api.handleListTaskRecords)
	mux.HandleFunc("GET /api/profile", api.handleGetProfile)
	mux.HandleFunc("GET /api/medals", api.handleListMedals)
	mux.HandleFunc("GET /api/medals/{id}", api.handleGetMedal)
	mux.HandleFunc("GET /api/points/ledger", api.handleListPointLedger)
	mux.HandleFunc("POST /api/orders", api.handleCreateOrder)
	mux.HandleFunc("GET /api/orders", api.handleListOrders)
	mux.HandleFunc("GET /api/addresses", api.handleListAddresses)
	mux.HandleFunc("POST /api/addresses", api.handleCreateAddress)
	mux.HandleFunc("PUT /api/addresses/{id}", api.handleUpdateAddress)
	mux.HandleFunc("DELETE /api/addresses/{id}", api.handleDeleteAddress)
	mux.HandleFunc("POST /api/addresses/{id}/default", api.handleSetDefaultAddress)
	mux.HandleFunc("GET /api/cart", api.handleListCart)
	mux.HandleFunc("GET /api/cart/count", api.handleCartCount)
	mux.HandleFunc("POST /api/cart", api.handleAddCartItem)
	mux.HandleFunc("PUT /api/cart/{id}", api.handleUpdateCartItem)
	mux.HandleFunc("DELETE /api/cart/{id}", api.handleDeleteCartItem)
	mux.HandleFunc("GET /api/favorites/item/{type}/{id}", api.handleFavoriteItemStatus)
	mux.HandleFunc("POST /api/favorites/item/{type}/{id}/toggle", api.handleToggleFavoriteItem)
	mux.HandleFunc("GET /api/admin/red-summary", api.handleAdminRedSummary)
	mux.HandleFunc("GET /api/admin/orders", api.handleAdminOrders)
	mux.HandleFunc("POST /api/admin/orders/{id}/status", api.handleAdminOrderStatus)
	mux.HandleFunc("GET /api/admin/data/{resource}", api.handleAdminResourceList)
	mux.HandleFunc("POST /api/admin/data/{resource}", api.handleAdminResourceCreate)
	mux.HandleFunc("PUT /api/admin/data/{resource}/{id}", api.handleAdminResourceUpdate)
	mux.HandleFunc("DELETE /api/admin/data/{resource}/{id}", api.handleAdminResourceDelete)
	mux.HandleFunc("POST /api/admin/login", api.handleAdminLogin)
	mux.HandleFunc("GET /api/admin/session", api.handleAdminSession)
	mux.HandleFunc("GET /api/admin/summary", api.handleAdminSummary)
	mux.HandleFunc("POST /api/admin/artworks", api.handleAdminCreateArtwork)
	mux.HandleFunc("PUT /api/admin/artworks/{id}", api.handleAdminUpdateArtwork)
	mux.HandleFunc("DELETE /api/admin/artworks/{id}", api.handleAdminDeleteArtwork)
	mux.HandleFunc("POST /api/admin/clips", api.handleAdminCreateClip)
	mux.HandleFunc("PUT /api/admin/clips/{id}", api.handleAdminUpdateClip)
	mux.HandleFunc("DELETE /api/admin/clips/{id}", api.handleAdminDeleteClip)
	mux.HandleFunc("POST /api/admin/activities", api.handleAdminCreateActivity)
	mux.HandleFunc("PUT /api/admin/activities/{id}", api.handleAdminUpdateActivity)
	mux.HandleFunc("DELETE /api/admin/activities/{id}", api.handleAdminDeleteActivity)
	mux.HandleFunc("GET /api/stats", api.handleStats)
	mux.HandleFunc("GET /api/artworks", api.handleListArtworks)
	mux.HandleFunc("GET /api/artworks/mine", api.handleMyArtworks)
	mux.HandleFunc("POST /api/artworks", api.handleCreateArtwork)
	mux.HandleFunc("GET /api/artworks/{id}", api.handleGetArtwork)
	mux.HandleFunc("PUT /api/artworks/{id}", api.handleUpdateArtwork)
	mux.HandleFunc("DELETE /api/artworks/{id}", api.handleDeleteArtwork)
	mux.HandleFunc("POST /api/artworks/{id}/like", api.handleLikeArtwork)
	mux.HandleFunc("POST /api/artworks/{id}/comments", api.handleCreateComment)
	mux.HandleFunc("DELETE /api/artworks/{id}/comments/{commentID}", api.handleDeleteComment)
	mux.HandleFunc("POST /api/upload", api.handleUpload)
	mux.HandleFunc("POST /api/upload/video", api.handleUploadVideo)
	mux.HandleFunc("GET /api/shadow-plays", api.handleListShadowPlays)
	mux.HandleFunc("POST /api/shadow-plays", api.handleCreateShadowPlay)
	mux.HandleFunc("GET /api/shadow-plays/{id}", api.handleGetShadowPlay)
	mux.HandleFunc("DELETE /api/shadow-plays/{id}", api.handleDeleteShadowPlay)
	mux.HandleFunc("GET /api/clips", api.handleListClips)
	mux.HandleFunc("POST /api/clips", api.handleCreateClip)
	mux.HandleFunc("POST /api/clips/{id}/like", api.handleLikeClip)
	mux.HandleFunc("GET /api/favorites/{type}/{id}", api.handleFavoriteStatus)
	mux.HandleFunc("POST /api/favorites/{type}/{id}/toggle", api.handleFavorite)
	mux.HandleFunc("GET /api/favorites/artworks", api.handleListFavoriteArtworks)
	mux.HandleFunc("GET /api/favorites/clips", api.handleListFavoriteClips)
	mux.HandleFunc("GET /api/activities", api.handleListActivities)
	mux.HandleFunc("GET /api/activities/{id}", api.handleGetActivity)
	mux.HandleFunc("POST /api/bookings", api.handleCreateBooking)
	mux.HandleFunc("GET /api/bookings", api.handleListBookings)
	mux.HandleFunc("POST /api/bookings/{id}/cancel", api.handleCancelBooking)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/admin/", serveAdminApp(resolveAdminDir()))
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/", http.StatusFound)
	})

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Println("影绘湘潭服务已启动: http://127.0.0.1" + addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}
