# Użyj oficjalnego obrazu Golang jako bazowego obrazu
FROM golang:1.21.4

# Ustaw katalog roboczy w kontenerze
WORKDIR /app

# Skopiuj zawartość bieżącego katalogu do kontenera
COPY . .

# Pobierz zależności z Go mod
RUN go mod download

# Skompiluj aplikację Go
RUN go build -o result_checker .

# Ostateczny obraz, użyj obrazu minimalnego dla uruchomienia aplikacji
FROM alpine:latest

# Zainstaluj Redis-cli
RUN apk --no-cache add redis

# Skopiuj plik wykonywalny z poprzedniego etapu do obecnego obrazu
COPY --from=0 /app/result_checker /app/result_checker

# Ustaw zmienną środowiskową dla adresu Redis
ENV REDIS_ADDR=redis:6379

# Zdefiniuj port, na którym będzie działać mikroserwis result_checker
EXPOSE 8080

# Uruchom aplikację result_checker
CMD ["/app/result_checker"]


FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tw-employee-data-service ./cmd/api

EXPOSE 8080 

CMD ["/app/result_checker]