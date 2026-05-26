# Go'nun resmi imajı
FROM golang:1.25-alpine

# Çalışma klasörü
WORKDIR /app

# go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./
RUN go mod download

# Tüm kodları kopyala
COPY . .

# Uygulamayı derle
RUN go build -o main .

# Uygulamayı çalıştır
CMD ["./main"]
