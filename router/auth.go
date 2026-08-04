package router

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gym/global"
	"gym/global/jwt"
	"gym/global/log"
	"gym/global/setting"
	"gym/global/time"
	"gym/models"
	"gym/models/ipblock"
	"gym/models/systemlog"
	"gym/models/user"
	"math/big"
	"strconv"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"resty.dev/v3"
)

type KakaoResponse struct {
	ID          int64  `json:"id"`
	ConnectedAt string `json:"connected_at"`
	Properties  struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNicknameNeedsAgreement bool `json:"profile_nickname_needs_agreement"`
		Profile                       struct {
			Nickname          string `json:"nickname"`
			IsDefaultNickname bool   `json:"is_default_nickname"`
		} `json:"profile"`
		HasEmail            bool   `json:"has_email"`
		EmailNeedsAgreement bool   `json:"email_needs_agreement"`
		IsEmailValid        bool   `json:"is_email_valid"`
		IsEmailVerified     bool   `json:"is_email_verified"`
		Email               string `json:"email"`
	} `json:"kakao_account"`
}

type NaverResponse struct {
	Response struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		Name         string `json:"name"`
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
	} `json:"response"`
}

type GoogleResponse struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Nonce         string `json:"nonce"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Type          string `json:"type"`
}

type AppleResponse struct {
	Iss            string `json:"iss"`
	Aud            string `json:"aud"`
	Exp            int64  `json:"exp"`
	Iat            int64  `json:"iat"`
	Sub            string `json:"sub"`
	CHash          string `json:"c_hash"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	AuthTime       int64  `json:"auth_time"`
	NonceSupported bool   `json:"nonce_supported"`
}

type ApplePublicKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type AppleKeys struct {
	Keys []ApplePublicKey `json:"keys"`
}

var JwtAuthRequired = func(c *fiber.Ctx) error {
	if values := c.Get("Authorization"); len(values) > 0 {
		str := values

		claims, err := jwt.Check(str)
		if err == nil {
			user := (*claims).User
			user.Passwd = ""
			c.Locals("user", &user)
			return c.Next()
		}
	}

	if isPublicPath(c.Method(), c.Path()) {
		return c.Next()
	}

	log.Info().Str("method", c.Method()).Str("path", c.Path()).Msg("Jwt header is broken")

	// 예전엔 nil 을 반환했는데, fiber 에서 c.Next() 없이 nil 을 돌려주면
	// 빈 본문에 200 이 나간다. 인증 실패는 401 로 명확히 끊는다.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
}

// isPublicPath 는 토큰 없이 허용할 경로를 화이트리스트로 정의한다.
//
// 로그인 전에 반드시 필요한 것만 연다:
//   - GET  /api/jwt                     로그인
//   - POST /api/user                    회원가입
//   - GET  /api/user/search/loginid     가입 시 아이디 중복확인
//
// 예전 구현은 "/api/user" 로 시작하면 전부 통과시켜서
// `GET /api/user` 로 전체 회원 명단(이메일·전화번호·주소·생년월일)이
// 무인증 조회됐다. 또 u.Path[:9] 는 9자 미만 경로(예: "/api/gym")에서
// 슬라이스 범위 초과로 패닉했다.
//
// 경로는 접두사가 아니라 완전 일치로 본다. 애매하면 막는 쪽(fail-closed).
func isPublicPath(method, path string) bool {
	switch path {
	case "/api/jwt":
		return method == fiber.MethodGet
	case "/api/user":
		return method == fiber.MethodPost
	case "/api/user/search/loginid":
		return method == fiber.MethodGet
	// 아래 셋은 JWT 대신 자체 서명을 검증하는 기계 대 기계 경로다.
	// PG/서명서비스 웹훅은 HMAC 웹훅 시크릿(billing.VerifyWebhookSignature),
	// 도어 verify 는 도어별 시크릿(X-Door-Sign) — 핸들러가 서명 불일치면 403.
	case "/api/billing/webhook":
		return method == fiber.MethodPost
	case "/api/contract/webhook":
		return method == fiber.MethodPost
	case "/api/door/verify":
		return method == fiber.MethodPost
	}
	return false
}

// 무인증 회원가입이 스스로에게 줄 수 있는 권한의 상한.
//
// UserUpdate 는 level·role 을 그대로 JSON 바인딩하고 Insert 가 DB 에 쓰기 때문에,
// 가입 요청만으로 최고관리자 계정을 만들 수 있었다.
//
// 상한은 실제 사용 형태에 맞췄다. 두 값이 다른 건 Level 과 Role 의 라벨
// 순서가 어긋나 있어서다 — 같은 숫자가 같은 권한을 뜻하지 않는다.
//
//	Level 3 = 헬스장관리자, 4 = 플랫폼관리자, 5 = 최고관리자
//	Role  4 = 헬스장관리자, 5 = 플랫폼관리자
//
// 웹 가입 화면이 회원·트레이너·헬스장관리자 중 고르게 돼 있어 헬스장 사장의
// 직접 가입이 정상 흐름이다. 그 최대치가 Level 3 / Role 4 다. 실 DB 102건에서도
// level 4 는 0명, level 5 와 role 5 는 최고관리자 1명뿐이라 이 상한은 지금
// 쓰는 기능을 하나도 막지 않으면서 권한 상승만 끊는다.
const (
	maxSignupLevel = 3
	maxSignupRole  = 4
)

// BlockPrivilegedSignup 은 무인증 가입 요청이 상한을 넘는 level·role 을
// 요구하면 거부한다. 값을 조용히 깎지 않고 403 으로 끊는 건, 정상 클라이언트는
// 애초에 이 범위를 넘기지 않으므로 넘어온다는 건 공격이거나 버그라서다.
//
// JwtAuthRequired 다음에 등록해야 한다 — 인증을 통과한 관리자가 만드는
// 계정은 c.Locals("user") 로 구분해 그대로 통과시킨다.
var BlockPrivilegedSignup = func(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost || c.Path() != "/api/user" {
		return c.Next()
	}
	if c.Locals("user") != nil {
		return c.Next()
	}

	level, role := signupPrivilege(c)
	if level > maxSignupLevel || role > maxSignupRole {
		log.Info().Int("level", level).Int("role", role).Msg("privileged signup blocked")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
	}
	return c.Next()
}

// MaskPublicUserLookup 은 무인증 아이디 중복확인 응답에서 개인정보를 지운다.
//
// /api/user/search/loginid 는 가입 흐름의 중복확인 때문에 무인증으로 열어둬야
// 하는데, 핸들러는 사용자 레코드 전문을 돌려준다. 아이디만 알면 그 사람의
// 이메일·전화번호·주소·생년월일·권한이 그대로 나갔다. 추측하기 쉬운
// 아이디(admin 등)면 관리자 정보가 공개되는 셈이다.
//
// 존재 확인에 필요한 필드만 남기고 나머지를 버린다. 응답 형태({code, item})는
// 유지한다 — 클라이언트가 item 유무로 중복을 판단하므로 구조를 바꾸면 가입
// 흐름이 깨진다.
//
// JwtAuthRequired 다음에 등록해야 c.Locals("user") 로 인증 여부를 구분할 수 있다.
var MaskPublicUserLookup = func(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	if c.Method() != fiber.MethodGet || c.Path() != "/api/user/search/loginid" {
		return nil
	}
	// 인증된 요청은 원래 응답을 그대로 쓴다(관리자 화면에서 필요).
	if c.Locals("user") != nil {
		return nil
	}

	var payload map[string]any
	if json.Unmarshal(c.Response().Body(), &payload) != nil {
		return nil
	}
	// item 이 객체가 아니면(미존재·에러 응답) 개인정보가 없으므로 손대지 않는다.
	item, ok := payload["item"].(map[string]any)
	if !ok {
		return nil
	}

	masked := map[string]any{}
	for _, k := range publicUserFields {
		if v, exists := item[k]; exists {
			masked[k] = v
		}
	}
	payload["item"] = masked

	body, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	c.Response().SetBody(body)
	return nil
}

// 무인증 조회에서 내보낼 필드. 중복 여부 판단에 필요한 최소한만 둔다.
var publicUserFields = []string{"id", "loginid"}

// signupPrivilege 는 요청 본문에서 level·role 을 읽는다.
// JSON 과 form(urlencoded·multipart)을 모두 본다. 한쪽만 검사하면 다른 쪽
// 인코딩으로 그대로 우회된다.
func signupPrivilege(c *fiber.Ctx) (level, role int) {
	var payload map[string]any
	if json.Unmarshal(c.Body(), &payload) == nil {
		level, role = privilegeValue(payload["level"]), privilegeValue(payload["role"])
	}
	if v, err := strconv.Atoi(c.FormValue("level")); err == nil && v > level {
		level = v
	}
	if v, err := strconv.Atoi(c.FormValue("role")); err == nil && v > role {
		role = v
	}
	return level, role
}

// privilegeValue 는 JSON 값을 정수로 바꾼다. 숫자는 float64 로 들어오고,
// 따옴표로 감싸 보내는 클라이언트도 있어 문자열도 함께 받는다.
func privilegeValue(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

func JwtAuth(c *fiber.Ctx, loginid string, passwd string) map[string]interface{} {
	loginType := c.Query("type")
	loginToken := c.Query("token")

	conn := models.NewConnection()
	defer conn.Close()

	userManager := models.NewUserManager(conn)
	loginlogManager := models.NewLoginlogManager(conn)
	ipblockManager := models.NewIpblockManager(conn)

	var item *models.User

	if loginType == "kakao" {
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", loginToken),
			"Content-Type":  "application/json",
		}

		url := "https://kapi.kakao.com/v2/user/me"
		resp, err := client.R().
			EnableTrace().
			SetHeaders(headers).
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		response := KakaoResponse{}
		global.JsonDecode(resp.String(), &response)

		connectid := fmt.Sprintf("%v", response.ID)
		if connectid == "" {
			log.Println(response)
			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}
	} else if loginType == "naver" {
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", loginToken),
			"Content-Type":  "application/json",
		}

		url := "https://openapi.naver.com/v1/nid/me"
		resp, err := client.R().
			EnableTrace().
			SetHeaders(headers).
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		response := NaverResponse{}
		global.JsonDecode(resp.String(), &response)

		connectid := fmt.Sprintf("%v", response.Response.ID)
		if connectid == "" {
			log.Println(response)
			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}
	} else if loginType == "google" {
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

		url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%v", loginToken)
		resp, err := client.R().
			EnableTrace().
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		response := GoogleResponse{}
		global.JsonDecode(resp.String(), &response)

		connectid := fmt.Sprintf("%v", response.Sub)
		if connectid == "" {
			log.Println(response)

			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}

	} else if loginType == "apple" {
		response, err := validateAppleToken(loginToken)
		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "apple token validation failed",
			}
		}

		connectid := fmt.Sprintf("%v", response.Sub)
		if connectid == "" {
			log.Println(response)

			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}

	} else {
		item = userManager.GetByLoginid(loginid)

		if item == nil {
			return map[string]interface{}{
				"code":    "error",
				"message": "user not found",
			}
		}

		if !jwt.CheckPasswd(item.Passwd, passwd) {
			return map[string]interface{}{
				"code":    "error",
				"message": "wrong password",
			}
		}

		if item.Use != user.UseUse {
			return map[string]interface{}{
				"code":    "error",
				"message": "status error",
			}
		}

		// instance := setting.GetInstance()
		// days := instance.SettingInt("user.passwd.changeday")

		// if days > 0 {
		// 	now := time.Now()
		// 	now = now.AddDate(0, 0, -1*days)

		// 	lastlogin := time.Parse(item.Lastchangepasswddate)
		// 	if now.After(lastlogin) {
		// 		return map[string]interface{}{
		// 			"code":    "error",
		// 			"message": "expired passwd",
		// 		}
		// 	}
		// }
	}

	typeid := ipblock.TypeAdmin
	if item.Level < user.LevelAdmin {
		typeid = ipblock.TypeNormal
	}

	ipblocks := ipblockManager.Find([]interface{}{
		models.Where{Column: "type", Value: typeid, Compare: "="},
		models.Where{Column: "use", Value: ipblock.UseUse, Compare: "="},
		models.Ordering("ib_order,ib_id"),
	})

	access := ipblock.PolicyGrant

	for _, v := range ipblocks {
		check, err := setting.NewIP(v.Address)
		if err != nil {
			continue
		}

		ip, _ := setting.NewIP(c.IP())

		if check.Contains(ip) {
			access = v.Policy
		}
	}

	if access == ipblock.PolicyDeny {
		systemlogManager := models.NewSystemlogManager(conn)
		item := models.Systemlog{
			Type:    systemlog.TypeLogin,
			Content: fmt.Sprintf("%v IP에서 %v 계정 접근이 거부되었습니다", c.IP(), loginid),
		}

		err := systemlogManager.Insert(&item)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		return map[string]interface{}{
			"code":    "error",
			"message": "access denied",
		}
	}

	now := time.Now()
	ip, _ := setting.NewIP(c.IP())
	loginlog := models.Loginlog{User: item.Id, Ip: c.IP(), Ipvalue: ip.StartInt, Date: now.Datetime()}
	err := loginlogManager.Insert(&loginlog)
	if err != nil {
		log.Error().Msg(err.Error())
		return map[string]interface{}{
			"code":    "error",
			"message": "db error",
		}
	}

	err = userManager.UpdateLogindateById(now.Datetime(), item.Id)
	if err != nil {
		log.Error().Msg(err.Error())
		return map[string]interface{}{
			"code":    "error",
			"message": "db error",
		}
	}

	item.Passwd = ""

	token := jwt.MakeToken(*item)

	return map[string]interface{}{
		"code":  "ok",
		"token": token,
		"user":  item,
	}
}

func validateAppleToken(idToken string) (*AppleResponse, error) {
	claims := jwtgo.MapClaims{}

	// 토큰을 파싱하되 audience 에러는 무시
	token, err := jwtgo.ParseWithClaims(idToken, claims, func(token *jwtgo.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		publicKey, err := getApplePublicKey(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	// audience 관련 에러는 무시
	if err != nil && !strings.Contains(err.Error(), "audience") {
		return nil, err
	}

	if token != nil && !token.Valid && err != nil && !strings.Contains(err.Error(), "audience") {
		return nil, fmt.Errorf("invalid token")
	}

	response := &AppleResponse{}

	if iss, exists := claims["iss"]; exists {
		response.Iss = iss.(string)
	}

	if aud, exists := claims["aud"]; exists {
		response.Aud = aud.(string)
	}

	if sub, exists := claims["sub"]; exists {
		response.Sub = sub.(string)
	}

	if exp, exists := claims["exp"]; exists {
		response.Exp = int64(exp.(float64))
	}

	if iat, exists := claims["iat"]; exists {
		response.Iat = int64(iat.(float64))
	}

	if authTime, exists := claims["auth_time"]; exists {
		response.AuthTime = int64(authTime.(float64))
	}

	if email, exists := claims["email"]; exists {
		response.Email = email.(string)
	}

	if emailVerified, exists := claims["email_verified"]; exists {
		if verified, ok := emailVerified.(string); ok {
			response.EmailVerified = verified
		} else if verified, ok := emailVerified.(bool); ok {
			if verified {
				response.EmailVerified = "true"
			} else {
				response.EmailVerified = "false"
			}
		}
	}

	if cHash, exists := claims["c_hash"]; exists {
		response.CHash = cHash.(string)
	}

	return response, nil
}

func getApplePublicKey(kid string) (*rsa.PublicKey, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	resp, err := client.R().
		Get("https://appleid.apple.com/auth/keys")

	if err != nil {
		return nil, err
	}

	var appleKeys AppleKeys
	global.JsonDecode(resp.String(), &appleKeys)

	for _, key := range appleKeys.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(key.N, key.E)
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func parseRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}

	var eInt int
	for _, b := range eBytes {
		eInt = eInt*256 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}
