# ---- build ----
# 1.26 최신 패치로 빌드해 표준 라이브러리 CVE(net/http·crypto/tls·x509 등)를 포함한다.
# govulncheck 로 검증하며, 버전을 내리면 stdlib 취약점이 되살아난다.
FROM --platform=linux/amd64 golang:1.26 AS build
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
