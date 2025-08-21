# FUNCTION.md

이 파일은 프로젝트의 각 폴더와 파일별 함수들과 동작 방식을 상세히 설명합니다.

## 📁 main.go
**진입점 파일**
- `main()` - 애플리케이션 시작점, 서비스 초기화 및 실행

## 📁 config/
**설정 관련 파일들**

### config/config.go
- `Init(dir string) ModelConfig` - model.json 파일을 읽어 설정을 초기화
- `GetPubspec() string` - pubspec.yml 파일에서 프로젝트명 추출

## 📁 router/
**라우팅 관련 파일들**

### router/router.go
- `SetRouter(app *fiber.App)` - 모든 REST API 라우트를 자동 생성하여 등록
  - CRUD 작업을 위한 GET, POST, PUT, DELETE 엔드포인트 생성
  - JWT 인증이 필요한 `/api/*` 경로와 인증 없는 `/api/jwt` 경로 설정

### router/auth.go  
**인증 및 JWT 관련**
- `JwtAuthRequired` - JWT 토큰 검증 미들웨어 함수
- `JwtAuth(c *fiber.Ctx, loginid string, passwd string) map[string]interface{}` - 로그인 처리 및 JWT 토큰 발급
  - 일반 로그인, 소셜 로그인(카카오, 네이버, 구글, 애플) 지원
  - IP 차단 정책 검사
  - 로그인 로그 기록
- `validateAppleToken(idToken string) (*AppleResponse, error)` - 애플 토큰 검증
- `getApplePublicKey(kid string) (*rsa.PublicKey, error)` - 애플 공개키 가져오기
- `parseRSAPublicKey(n, e string) (*rsa.PublicKey, error)` - RSA 공개키 파싱

## 📁 controllers/
**컨트롤러 관련 파일들**

### controllers/controllers.go
**베이스 컨트롤러**
- `NewController(g *fiber.Ctx) *Controller` - 새 컨트롤러 생성
- `(c *Controller) Init(g *fiber.Ctx)` - 컨트롤러 초기화
- `(c *Controller) NewConnection() *Connection` - DB 연결 생성
- `(c *Controller) Close()` - 리소스 정리
- `(c *Controller) Lock()` / `(c *Controller) Unlock()` - 뮤텍스 잠금/해제
- `(c *Controller) Error(err error)` - 에러 설정
- `(c *Controller) Set(name string, value interface{})` - 결과값 설정
- `(c *Controller) Get(name string) string` - 쿼리 파라미터 가져오기
- `(c *Controller) Geti(name string) int` - 정수형 쿼리 파라미터 가져오기
- `(c *Controller) Bind(obj interface{}) error` - 요청 본문을 객체에 바인딩
- `(c *Controller) Paging(page int, totalRows int, pageSize int)` - 페이징 정보 계산
- `(c *Controller) GetUpload(uploadPath string, name string) (string, string)` - 파일 업로드 처리
- `(c *Controller) Download(filename string, downloadFilename string)` - 파일 다운로드

### controllers/rest/*.go
**각 모델별 REST 컨트롤러 (auto-generated)**
각 모델(User, Gym, Payment 등)마다 다음 메소드들을 제공:
- `Read(id int64)` - 단일 항목 조회
- `Index(page int, pagesize int)` - 목록 조회 (페이징, 검색 기능 포함)
- `Count()` - 전체 개수 조회  
- `Insert(item *Model)` - 새 항목 생성
- `Update(item *Model)` - 기존 항목 수정
- `Delete(item *Model)` - 항목 삭제
- `Insertbatch(item *[]Model)` - 일괄 생성
- `Deletebatch(item *[]Model)` - 일괄 삭제
- 모델별 특화 메소드들 (예: GetByLoginid, FindByLevel 등)

## 📁 models/
**데이터 모델 및 데이터베이스 관련**

### models/db.go
**데이터베이스 연결 관리**
- `GetConnection() *Connection` - DB 연결 생성
- `NewConnection() *Connection` - 재시도 로직을 포함한 DB 연결 생성
- `(c *Connection) Close()` - 연결 종료
- `(c *Connection) Begin()` - 트랜잭션 시작
- `(c *Connection) Commit() error` - 트랜잭션 커밋
- `(c *Connection) Rollback()` - 트랜잭션 롤백
- `(c *Connection) Exec(query string, params ...interface{}) (sql.Result, error)` - SQL 실행
- `(c *Connection) Query(query string, params ...interface{}) (*sql.Rows, error)` - SQL 쿼리
- `InitDate() string` - 초기 날짜 반환
- `Paging(page int, pagesize int) PagingType` - 페이징 객체 생성
- `Ordering(order string) OrderingType` - 정렬 객체 생성
- `Limit(limit int) LimitType` - 제한 객체 생성

### models/cache.go
- `InitCache()` - 캐시 초기화

### models/{model}.go & models/{model}/{model}.go
**각 모델별 데이터베이스 관리 클래스 (auto-generated)**
각 모델마다 다음과 같은 매니저 클래스와 메소드들 제공:
- `New{Model}Manager(conn *Connection) *{Model}Manager` - 매니저 생성
- `(p *{Model}Manager) Close()` - 매니저 종료
- `(p *{Model}Manager) Insert(item *{Model}) error` - 데이터 삽입
- `(p *{Model}Manager) Update(item *{Model}) error` - 데이터 수정
- `(p *{Model}Manager) Delete(id int64) error` - 데이터 삭제
- `(p *{Model}Manager) Get(id int64) *{Model}` - ID로 단일 조회
- `(p *{Model}Manager) Find(args []interface{}) []{Model}` - 조건별 다중 조회
- `(p *{Model}Manager) Count(args []interface{}) int` - 개수 조회
- `(p *{Model}Manager) GroupBy(name string, args []interface{}) []Groupby` - 그룹별 집계
- `(p *{Model}Manager) UpdateWhere(columns []Params, args []interface{}) error` - 조건별 업데이트
- `(p *{Model}Manager) DeleteWhere(args []interface{}) error` - 조건별 삭제
- 컬럼별 개별 업데이트 메소드들 (예: UpdateName, UpdateEmail 등)
- 특화 검색 메소드들 (예: GetByLoginid, FindByLevel 등)

### 모델별 상수 파일 (models/{model}/{model}.go)
각 모델의 enum 값들과 변환 함수들:
- Column 상수들 (ColumnId, ColumnName 등)
- Status/Type/Level 등의 enum 상수들
- Get{EnumName}() - enum을 문자열로 변환
- Find{EnumName}() - 문자열을 enum으로 변환
- Convert{EnumName}() - 배열 변환

## 📁 global/
**전역 유틸리티 및 공통 기능**

### global/global.go
**유틸리티 함수들**
- `ToMap(slice []string) map[string]int` - 슬라이스를 맵으로 변환
- `ReverseMap(inmap map[int]string) map[string]int` - 맵 역변환
- `Atoi(value string) int` - 문자열을 정수로 변환 (콤마 제거)
- `Atol(value string) int64` - 문자열을 int64로 변환
- `Atof(value string) float64` - 문자열을 float64로 변환
- `UUID() string` - UUID 생성
- `RandomString(n int) string` - 랜덤 문자열 생성
- `GetSha256(str string) string` - SHA256 해시 생성
- `StripTags(content string) string` - HTML 태그 제거
- `FindImages(htm string) []string` - HTML에서 이미지 URL 추출
- `IsEmptyDate(date string) bool` - 빈 날짜 확인
- `SendSMS(tel string, content string) bool` - SMS 전송
- `WriteFile(filename string, content string) error` - 파일 쓰기
- `ReadFile(filename string) string` - 파일 읽기
- `Substr(str string, start int, end int) string` - 문자열 자르기 (유니코드 지원)
- `Strlen(s string) int` - 문자열 길이 (유니코드 지원)
- `DownloadImage(url string, filename string) int64` - 이미지 다운로드
- `MakeUniqueSlice[T](arr []T) []T` - 중복 제거
- `JsonEncode(item interface{}) string` - JSON 인코딩
- `JsonDecode(str string, item interface{}) error` - JSON 디코딩
- `XmlEncode(item interface{}) string` - XML 인코딩
- `XmlDecode(str string, item interface{}) error` - XML 디코딩
- `MakeSearchKeyword(str string) []string` - 검색 키워드 생성
- `Reverse[T any](original []T) (reversed []T)` - 슬라이스 역순 정렬
- `SendMail(email string, title string, content string) error` - 이메일 전송

### global/config/config.go
**설정 관리**
- `Init()` - 환경설정 초기화 (.env.yml 파일 로드)
- `Get(name string) interface{}` - 설정값 가져오기
- `GetString(name string) string` - 문자열 설정값 가져오기
- `GetInt(name string) int` - 정수 설정값 가져오기

### global/jwt/jwt.go
**JWT 토큰 관리**
- `Check(str string) (*AuthTokenClaims, error)` - JWT 토큰 검증
- `MakeToken(item models.User) string` - JWT 토큰 생성
- `CheckPasswd(dbPasswd string, inputPasswd string) bool` - 비밀번호 검증
- `GeneratePasswd(passwd string) (string, error)` - 비밀번호 해싱

### global/log/log.go
**로깅 관리**
- `Rotate()` - 로그 파일 로테이션
- `init()` - 로깅 시스템 초기화 (zerolog 설정)
- `Println(a ...any)` - 로그 출력
- `Get() *zerolog.Logger` - 로거 인스턴스 가져오기
- `Debug() *zerolog.Event` - 디버그 로그
- `Info() *zerolog.Event` - 정보 로그
- `Warn() *zerolog.Event` - 경고 로그
- `Error() *zerolog.Event` - 에러 로그

### global/time/time.go
**시간 관련 유틸리티**
- `Now() *Time` - 현재 시간
- `Parse(str string) *Time` - 문자열 파싱
- `(c *Time) Clone() *Time` - 시간 복사
- `(c *Time) Timestamp() int64` - 타임스탬프
- `(c *Time) String() string` - 문자열 변환
- `(c *Time) Datetime() string` - 날짜시간 문자열
- `(c *Time) Date() string` - 날짜 문자열
- `(c *Time) Time() string` - 시간 문자열
- `(c *Time) Year() int` - 년도 추출
- `(c *Time) Month() int` - 월 추출
- `(c *Time) Day() int` - 일 추출
- `(c *Time) Add(duration time.Duration) *Time` - 시간 더하기
- `(c *Time) AddDate(years int, months int, days int) *Time` - 날짜 더하기
- `(c *Time) After(value *Time) bool` - 이후 시간 확인
- `(c *Time) Before(value *Time) bool` - 이전 시간 확인
- `Sleep(d time.Duration)` - 대기
- `After(d time.Duration) <-chan time.Time` - 타이머

### global/excel.go
**Excel 파일 처리**
- `NewExcelReader(filename string) *Excel` - Excel 리더 생성
- `OpenExcel(filename string, title string, fontSize float64, header []string, width []int, align []string) *Excel` - Excel 파일 열기
- `New() *Excel` - 새 Excel 객체 생성
- `NewExcel(title string, sheet string, fontSize float64, header []string, width []int, align []string) *Excel` - 새 Excel 파일 생성
- `(p *Excel) GetCell(col string, row int) string` - 셀 값 읽기
- `(p *Excel) SetSheet(str string)` - 시트 설정
- `(p *Excel) NewSheet(sheet string, header []string, width []int, align []string)` - 새 시트 추가
- `(p *Excel) Save(filename string) string` - 파일 저장
- `(p *Excel) Cell(str string) string` - 셀에 값 쓰기
- `(p *Excel) CellInt(value int)` - 정수값 쓰기
- `(p *Excel) CellPrice(value int)` - 가격 형식으로 쓰기
- `(p *Excel) CellImage(filename string)` - 이미지 삽입
- `(p *Excel) Close()` - Excel 객체 정리

### global/pdf.go & global/pdf/*.go
**PDF 생성**
- `NewPdf(title string, header []string, width []int, align []string, headerFont float64, bodyFont float64) *Pdf` - PDF 객체 생성
- `(p *Pdf) Save(filename string) string` - PDF 파일 저장
- PDF/관련 함수들:
  - `NewPdf() *Pdf` - 새 PDF 생성
  - `(c *Pdf) SetFont(font string)` - 폰트 설정
  - `(c *Pdf) SetFontSize(size float64)` - 폰트 크기 설정
  - `(c *Pdf) AddPage()` - 새 페이지 추가
  - `(c *Pdf) TextOut(x, y, width, height float64, str string, align int)` - 텍스트 출력
  - `(c *Pdf) FillRect(x, y, width, height float64, color Color)` - 사각형 그리기
  - `(c *Pdf) Save(filename string)` - PDF 저장

### global/image.go
**이미지 처리**
- `MakeThumbnail(w int, h int, filename string, targetFilename string)` - 썸네일 생성

### global/notify.go
**알림 시스템**
- `init()` - 알림 채널 초기화
- `SendNotify(id int64, UUID string, message MessageType)` - 단일 알림 전송
- `SendNotifys(ids []int64, message MessageType)` - 다중 알림 전송
- `GetChannel() chan Notify` - 알림 채널 가져오기
- `GetMessage(message MessageType) string` - 메시지 내용 생성

### global/fcm.go
**Firebase Cloud Messaging**
- `init()` - FCM 채널 초기화
- `SendFcm(item Fcm)` - FCM 메시지 전송
- `GetFcm() chan Fcm` - FCM 채널 가져오기

### global/cron.go
**크론 작업 관리**
- `init()` - 크론 채널 초기화
- `RestartCron()` - 크론 재시작
- `GetCronChannel() chan bool` - 크론 채널 가져오기

### global/setting/instance.go
**설정 인스턴스 관리**
- `GetInstance() *Instance` - 설정 인스턴스 가져오기 (싱글톤)
- `(c *Instance) InitSetting()` - 설정 초기화
- `(c *Instance) Init()` - 인스턴스 초기화
- `(c *Instance) Setting(key string) string` - 설정값 조회
- `(c *Instance) SettingInt(key string) int` - 정수 설정값 조회

### global/setting/ip.go
**IP 주소 관리**
- `NewIP(str string) (*IP, error)` - IP 객체 생성
- `(c *IP) Match(dest *IP) bool` - IP 매칭 확인
- `(c *IP) Contains(dest *IP) bool` - IP 포함 확인
- `(c *IP) Equal(dest *IP) bool` - IP 동일성 확인
- `MatchIP(src IP, dest IP) bool` - IP 매칭 함수

## 📁 services/
**백그라운드 서비스들**

### services/http.go
**HTTP 서버**
- `Http()` - HTTP 서버 시작 (Fiber 웹서버 구동)

### services/cron.go  
**크론 작업 스케줄러**
- `Cron()` - 크론 서비스 시작
- `InitCron()` - 크론 작업 초기화
  - 로그 로테이션 (매일 자정)
  - 게시글 조회수 증가 (30분마다)
- `RestartCron()` - 크론 재시작
- `boardCounter()` - 게시글 카운터 증가 작업

### services/chat.go
**채팅/WebSocket 서비스**
- `Chat()` - 채팅 서비스 시작 (WebSocket 서버)
- `(p *ChatService) SendTo(notify global.Notify)` - 특정 사용자에게 메시지 전송

### services/notify.go
**알림 서비스**
- `Notify()` - 알림 서비스 시작 (실시간 알림 처리)

### services/fcm.go
**Firebase Cloud Messaging 서비스**
- `Fcm()` - FCM 서비스 시작 (푸시 알림 전송)

## 🔧 자동 생성 파일 특징

### REST 컨트롤러들 (controllers/rest/*.go)
- buildtool-router 도구에 의해 자동 생성됨
- 각 데이터베이스 테이블마다 표준 CRUD 작업 제공
- 검색, 페이징, 정렬 기능 자동 포함
- 날짜 범위 검색, LIKE 검색 등 고급 검색 기능

### 데이터 모델들 (models/*.go)
- buildtool-model 도구에 의해 자동 생성됨  
- ORM 스타일의 데이터베이스 액세스 레이어
- 트랜잭션 지원, 연결 풀링
- 타입 안전한 쿼리 빌더
- 자동 페이징 및 정렬

### 라우터 (router/router.go)
- buildtool-router 도구에 의해 자동 생성됨
- RESTful API 엔드포인트 자동 생성
- JWT 인증 자동 적용
- 표준 HTTP 메소드 매핑 (GET, POST, PUT, DELETE)

## 📊 데이터 흐름

1. **요청 처리**: HTTP 요청 → router/router.go → JWT 검증 → 해당 Controller
2. **데이터 처리**: Controller → Model Manager → Database Connection
3. **응답 생성**: 처리 결과 → Controller.Result → JSON 응답
4. **백그라운드**: Services (Cron, Chat, Notify, FCM) 독립 실행

이 구조는 모든 기본적인 CRUD 작업을 자동화하고, 확장 가능한 아키텍처를 제공합니다.