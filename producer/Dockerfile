# Используем образ с Go для простоты
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модуля Go
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o producer ./cmd/

# Указываем команду по умолчанию с конфигурационным файлом
CMD ["./producer", "-config", "./configs/producer.yaml"]
