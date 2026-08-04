package jwt

import (
	"errors"
	"fmt"
	"gym/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokenClaims struct {
	User               models.User `json:"user"`
	jwt.StandardClaims             // 표준 토큰 Claims
}

// 서명 키는 JWT_SECRET 환경변수로만 받는다. 예전엔 실제 시크릿이 소스에
// 상수로 박혀 있었는데, 이 저장소는 공개라 그대로 유출됐다 — 그 키로 서명한
// 토큰은 누구나 위조할 수 있다. 결제 등 돈이 걸린 API 를 붙이기 전에 제거했다.
//
// production 에서 JWT_SECRET 이 없으면 기동 자체를 막는다(fail-fast).
// 개발 모드에서만 고정 기본값을 쓴다 — 유출된 옛 키와 다른 값이라 옛 토큰은
// 어디서도 더 이상 유효하지 않다. 키 교체로 기존 로그인은 전부 무효화된다.
var _secretCode string = loadSecret()

func loadSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	if os.Getenv("APP_MODE") == "production" {
		fmt.Fprintln(os.Stderr, "JWT_SECRET is required in production mode")
		os.Exit(1)
	}
	return "gym-dev-only-secret"
}

// Secret 은 파생 키(도어 QR 서명 등) 용도로 시크릿을 노출한다.
// 그대로 다른 서명에 쓰지 말고 반드시 용도 접두사를 붙여 파생할 것.
func Secret() string {
	return _secretCode
}

func Check(str string) (*AuthTokenClaims, error) {
	if len(str) < 7 || str[:7] != "Bearer " {
		err := errors.New("tokek is broken")
		return nil, err
	}

	token := str[7:]

	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(_secretCode), nil
	}

	_, err := jwt.ParseWithClaims(token, &claims, key)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func MakeToken(item models.User) string {
	now := time.Now()
	at := AuthTokenClaims{
		User: item,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.At(now),
			ExpiresAt: jwt.At(now.Add(time.Hour * 24 * 365 * 10)),
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, _ := atoken.SignedString([]byte(_secretCode))

	return signedAuthToken
}

func CheckPasswd(dbPasswd string, inputPasswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPasswd), []byte(inputPasswd))
	return err == nil
}

func GeneratePasswd(passwd string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hashedPasswd), err
	}
}
