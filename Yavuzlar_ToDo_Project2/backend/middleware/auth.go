package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" //unıque
)
var JwtKey []byte
type contextKey string
const UserIDKey contextKey = "userId"

// Token uretırekn uuid.UUID 
func GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID.String(), // UUID'yi string olarak JWT'nin içine gömüyoruz
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(JwtKey)
}
// APIMiddleware
func APIMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifyToken(w, r, next)
	}
}
// WebMiddleware: index.html gibi sayfalar için
func WebMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifyToken(w, r, next)
	}
}
// Ortak token doğrulayıcı
func verifyToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return 
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return 
	}
	//UUID donusumu
	userIDStr, ok := claims["userId"].(string)
	if !ok {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}
	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	// BAŞARILIYSA BURAYA GELİR
	ctx := context.WithValue(r.Context(), UserIDKey, parsedUserID)
	next.ServeHTTP(w, r.WithContext(ctx)) 
}