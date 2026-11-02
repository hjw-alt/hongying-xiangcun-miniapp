package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

const adminTokenTTL = 12 * time.Hour

type adminClaims struct {
	Username string `json:"username"`
	Expires  int64  `json:"expires"`
}

func adminCredentials() (string, string) {
	username := strings.TrimSpace(os.Getenv("ADMIN_USER"))
	if username == "" {
		username = "hjw"
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "123456"
	}
	return username, password
}

func adminSecret() []byte {
	secret := os.Getenv("ADMIN_SECRET")
	if secret == "" {
		secret = "shadowart-xiangtan-admin-secret"
	}
	return []byte(secret)
}

func signAdminToken(username string) (string, int64) {
	expires := time.Now().Add(adminTokenTTL).Unix()
	payload, err := json.Marshal(adminClaims{Username: username, Expires: expires})
	if err != nil {
		panic(err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, adminSecret())
	mac.Write([]byte(encoded))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encoded + "." + signature, expires
}

func verifyAdminToken(token string) (adminClaims, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return adminClaims{}, false
	}
	mac := hmac.New(sha256.New, adminSecret())
	mac.Write([]byte(parts[0]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return adminClaims{}, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return adminClaims{}, false
	}
	var claims adminClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return adminClaims{}, false
	}
	if claims.Username == "" || time.Now().Unix() > claims.Expires {
		return adminClaims{}, false
	}
	return claims, true
}

func adminTokenFromRequest(r *http.Request) string {
	if token := strings.TrimSpace(r.Header.Get("X-Admin-Token")); token != "" {
		return token
	}
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authorization, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	}
	return ""
}

func requireAdmin(w http.ResponseWriter, r *http.Request) (adminClaims, bool) {
	claims, ok := verifyAdminToken(adminTokenFromRequest(r))
	if !ok {
		writeError(w, http.StatusUnauthorized, "管理员登录已失效，请重新登录")
		return adminClaims{}, false
	}
	return claims, true
}

func (a *API) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		return
	}
	username, password := adminCredentials()
	input.Username = strings.TrimSpace(input.Username)
	if input.Username != username || input.Password != password {
		writeError(w, http.StatusUnauthorized, "账号或密码不正确")
		return
	}
	token, expires := signAdminToken(username)
	writeJSON(w, http.StatusOK, map[string]any{
		"token":     token,
		"username":  username,
		"expiresAt": expires,
	})
}

func (a *API) handleAdminSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"username":  claims.Username,
		"expiresAt": claims.Expires,
	})
}

func (a *API) handleAdminSummary(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAdmin(w, r); !ok {
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	summary, err := a.Store.AdminSummary(ctx)
	if err != nil {
		handleStoreError(w, err, "读取后台统计")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
