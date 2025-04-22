# Первый этап - сборка приложения
FROM golang:alpine AS builder

WORKDIR /app

# Создаем директорию для статических файлов
RUN mkdir -p /app/static

# Копируем шаблон HTML в директорию static
COPY static/index.html /app/static/

# Копируем исходный код Go
COPY main.go go.mod go.sum* /app/

# Включаем CGO (если требуется, в большинстве случаев можно отключить)
ENV CGO_ENABLED=0

RUN --mount=type=secret,id=api_key \
    export API_KEY=$(cat /run/secrets/api_key.txt) && \
    # Компилируем приложение, флаги для максимальной оптимизации размера
    go build -ldflags="-s -w -X main.apiKey=$API_KEY" -o weather-app .


# Второй этап - создание минимального образа
FROM scratch

# Определяем аргумент сборки с значением по умолчанию
ARG PORT=3000

# Устанавливаем переменную окружения из аргумента
ENV PORT=${PORT}

# Копируем SSL-сертификаты из образа Go для HTTPS запросов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем исполняемый файл из предыдущего этапа
COPY --from=builder /app/weather-app /weather-app

# Устанавливаем рабочую директорию
WORKDIR /

# Открываем порт
EXPOSE ${PORT}

# Запускаем приложение
ENTRYPOINT ["/weather-app"]