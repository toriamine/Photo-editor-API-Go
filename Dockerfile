# Базовый образ для сборки
FROM golang:alpine AS builder

LABEL stage=gobuilder

# Отключаем CGO (нужно для совместимости)
ENV CGO_ENABLED = 0

# Устанавливаем зависимости
RUN apk update --no-cache && apk add --no-cache tzdata

# Рабочая директория
WORKDIR /build

# Загружаем зависимости заранее (ускоряет билд)
ADD go.mod .
ADD go.sum .
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем бинарник
RUN go build -ldflags="-s -w" -o /app/main main.go

# Финальный минимальный образ
FROM scratch

# Копируем сертификаты и таймзону
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ = Europe/Moscow

# Указываем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл
COPY --from=builder /app/main /app/main

# Копируем HTML-шаблоны и папку с загрузками
COPY html/template /app/html/template
COPY uploads /app/uploads

# Открываем порт (нужно для Kubernetes)
EXPOSE 8081

# Запуск приложения
CMD ["./main"]