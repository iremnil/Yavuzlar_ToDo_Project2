package config

import (
	"fmt"
	"log"
	"os"
	"todo-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
// Projenin her yerinden bu değişkene erişilmeli ,veritabanı işlemleri yapmaya yarayacak.
var DB *gorm.DB
func ConnectDatabase() {
	// Docker için Environment Variables 
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	// GORM ile PostgreSQLe bağlantısı
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı kurulamadı!: ", err)
	}
	// models/todo.go'da yazılan todo structinı veritabanında otomatik tabloya dönüştürür
	err = database.AutoMigrate(&models.Todo{}, &models.User{})
	if err != nil {
    	log.Fatal("Tablolar oluşturulurken hata çıktı!: ", err)
	}
	// Bağlantı başarılıysa global değişkenimize aktarıyoruz
	DB = database
	fmt.Println("PostgreSQL bağlantısı başarılı ve Todo tablosu hazır!")
}