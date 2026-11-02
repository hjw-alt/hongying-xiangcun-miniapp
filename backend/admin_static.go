package main

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func resolveAdminDir() string {
	candidates := []string{
		"admin/dist",
		filepath.Join("..", "admin", "dist"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return candidates[0]
}

func serveAdminApp(root string) http.Handler {
	rootPath, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relative := strings.TrimPrefix(r.URL.Path, "/admin")
		if relative == "" || relative == "/" {
			relative = "/index.html"
		}
		target, err := filepath.Abs(filepath.Join(rootPath, filepath.FromSlash(relative)))
		if err != nil || !strings.HasPrefix(target, rootPath) {
			http.NotFound(w, r)
			return
		}
		if info, statErr := os.Stat(target); statErr != nil || info.IsDir() {
			target = filepath.Join(rootPath, "index.html")
		}
		data, readErr := os.ReadFile(target)
		if readErr != nil {
			http.Error(w, "后台管理页面尚未构建，请在 admin 目录执行 npm install 与 npm run build", http.StatusNotFound)
			return
		}
		contentType := mime.TypeByExtension(filepath.Ext(target))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		if !strings.Contains(contentType, "charset") {
			contentType += "; charset=utf-8"
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
}
