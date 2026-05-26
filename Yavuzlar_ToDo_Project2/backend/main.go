package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"todo-backend/config"
	"todo-backend/middleware"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" 
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendError(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, "Şifre işlenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Bu kullanıcı adı zaten alınmış!"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mesaj": "Kayıt başarılı! Şimdi giriş yapabilirsin."})
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	//cors ve optıons ıznı ---
	enableCors(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	//contexten UUID formatındaki kullanıcı idsini guvenlı almak
	userId, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Yetkisiz erişim!"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//Todo lısteleme GET 
    if r.Method == http.MethodGet {
        var todos []models.Todo
        // Sıralama komutu 
        config.DB.Where("user_id = ?", userId).Order("created_at asc").Find(&todos)
        json.NewEncoder(w).Encode(todos)
        return
    }
	//Yenı todo olusturma(POST) 
	if r.Method == http.MethodPost {
		var yeniGorev models.Todo
		// Frontendden gelen JSON'ı title ve completed çözüyoruz
		if err := json.NewDecoder(r.Body).Decode(&yeniGorev); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"mesaj": "Geçersiz veri formatı"})
			return
		}
		yeniGorev.UserID = userId // Görevi giriş yapan UUID'li kullanıcıya bağlıyoruz
		
		if err := config.DB.Create(&yeniGorev).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"mesaj": "Veritabanına kaydedilirken hata oluştu"})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(yeniGorev)
		return
	}
	//todo silme DELETE 
	if r.Method == http.MethodDelete {
		todoId := r.URL.Query().Get("id")
		if todoId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Güvenli silme:Sadece kendi görevini silebilir
		config.DB.Where("id = ? AND user_id = ?", todoId, userId).Delete(&models.Todo{})
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Silindi"})
		return
	}

	// Todo guncelleme (PUT)
	if r.Method == http.MethodPut {
		todoId := r.URL.Query().Get("id")
		if todoId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var guncelVeri models.Todo
		json.NewDecoder(r.Body).Decode(&guncelVeri)
		// Güvenli güncelleme:Sadece kendi görevini güncelleyebilir
		// GÜVENLİ GÜNCELLEME BLOĞU ŞU ŞEKİLDE OLACAK:
		config.DB.Model(&models.Todo{}).Where("id = ? AND user_id = ?", todoId, userId).Updates(map[string]interface{}{
			"title":  guncelVeri.Title,
			"status": guncelVeri.Status, // completed yerine status oldu
		})
		
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Güncellendi"})
		return
	}
	// GET, POST, PUT veya DELETE dışında bir istek gelirse engelle
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"mesaj": "Metot izinli değil"})
}
//KULLANICI BİLGİSİNİ DÖNEN ENDPOINT
func MeHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var user models.User
	// Veritabanından sadece bu UUID'ye ait kullanıcının adını çek
	config.DB.Select("username").Where("id = ?", userId).First(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"username": user.Username})
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Zaten giriş yapmışsınız!"}) //Kullanıcı gırıs yapmısssa tekrar logın olamasın
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Hatalı istek"})
		return
	}
	var user models.User
	if err := config.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Kullanıcı adı veya şifre hatalı!"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"mesaj": "Kullanıcı adı veya şifre hatalı!"})
		return
	}
	//uuidyi string olarak JWT içine koyuyoruz
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID.String(), 
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), 
	})

	tokenString, _ := token.SignedString(middleware.JwtKey)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, 
		Path:     "/",  
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mesaj": "Giriş başarılı!"})
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"mesaj": message})
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Uyarı: .env dosyası bulunamadı.")
	}
	middleware.JwtKey = []byte(os.Getenv("JWT_SECRET"))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), 
		HttpOnly: true,
		Path:     "/",
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mesaj": "Çıkış başarılı!"})
}
func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Todo{}, &models.User{})
	//Statk dosyaları tanımla (cakısmaları onlemek için FileServer yerine ServeFile)
	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./frontend/login.html") })
	http.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./frontend/register.html") })
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./frontend/style.css") })
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./frontend/script.js") })
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	//korumalı rotaları tanımla
	http.HandleFunc("/", middleware.WebMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Sadece ana dizine girilmesine izin ver, saçma URL'leri engelle
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./frontend/index.html")
	}))

	// API rotaları
	http.HandleFunc("/api/v1/todos", middleware.APIMiddleware(TodoHandler))
	http.HandleFunc("/api/v1/me", middleware.APIMiddleware(MeHandler))
	// Auth gerektirmeyen API'ler
	http.HandleFunc("/api/v1/register", RegisterHandler)
	http.HandleFunc("/api/v1/login", LoginHandler)
	http.HandleFunc("/api/v1/logout", LogoutHandler)

	fmt.Println("Backend 8080 portunda çalışıyor...")
	http.ListenAndServe(":8080", nil)
}