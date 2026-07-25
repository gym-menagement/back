# ---- build ----
FROM --platform=linux/amd64 golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ---- runtime ----
FROM --platform=linux/amd64 alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Seoul /etc/localtime && echo "Asia/Seoul" > /etc/timezone
ENV TZ=Asia/Seoul
WORKDIR /app
COPY --from=build /app/server /app/server
EXPOSE 8004
CMD ["/app/server"]
