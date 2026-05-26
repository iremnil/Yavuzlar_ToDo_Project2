# Yavuzlar_ToDo_Project2
#Bu proje, ilk hafta geliştirilmiş olan HTML/CSS/JavaScript tabanlı To-Do uygulamasının; **GoLang (REST API)**, **PostgreSQL (Veritabanı)** ve **Docker** bileşenleri kullanılarak bir web uygulaması mimarisine dönüştürülmüş  ve dockerize edilmiş halidir.
Sistemde 3 farklı container ayağa kalkacaktır:
Database Servisi (PostgreSQL): Verilerin kalıcı olarak saklandığı katman.
Backend Servisi (GoLang): REST API endpoints ve JWT kimlik doğrulama mekanizması.
Frontend Servisi (Nginx): Frontend dosyalarının containerize edilmiş hali.

Uygulamayı test etmek, tüm fonksiyonları (CRUD ve Auth) eksiksiz kullanmak için http://localhost:8080 ziyaret ediniz.
Projede JWT tabanlı kimlik doğrulama (HttpOnly Cookie yönetimi) ve tarayıcıların CORS (Cross-Origin Resource Sharing) politikalarının en kararlı şekilde çalışması amacıyla; frontend statik dosyaları (HTML/CSS/JS) doğrudan Go servisi üzerinden http.ServeFile sunulmaktadır.
Tarayıcıdan http://localhost:80 adresine gidildiğinde, Nginx imajı frontend arayüzünü render eder.
Backend: GoLang net/http kütüphanesi ile REST API mimarisi.
Base Path: Tüm API endpoint'leri /api/v1 formatındadır.
ORM & DB: Veritabanı işlemleri için GORM kullanılmış, veritabanı olarak PostgreSQL seçilmiştir.
Veri Tipi: To-Do nesnelerindeki id alanları UUID veri tipindedir.
Authentication: JWT mekanizması kurulmuştur. Login gerçekleşmeden ve geçerli token sağlanmadan CRUD işlemlerine izin verilmez.
Veri Kalıcılığı: Docker container yeniden başlatılsa dahi verilerin kaybolmaması için pgdata adında bir docker volume yapısı kurulmustur.
