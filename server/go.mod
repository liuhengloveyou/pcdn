module pcdn-server

go 1.24.3

require (
	github.com/bytedance/sonic v1.13.2
	github.com/go-playground/validator/v10 v10.26.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/liuhengloveyou/go-common v0.0.0-20250319112824-c28f82e5a12b
	github.com/liuhengloveyou/go-errors v0.0.0-20211025085721-e6717f3d23d1
	github.com/liuhengloveyou/passport v1.1.0
	github.com/redis/go-redis/v9 v9.8.0
	go.uber.org/zap v1.27.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.26.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Blank-Xu/sql-adapter v1.0.0 // indirect
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/casbin/casbin/v2 v2.100.0 // indirect
	github.com/casbin/govaluate v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1005 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.1005 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.10.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	xorm.io/builder v0.3.13 // indirect
)

replace github.com/liuhengloveyou/passport v1.1.0 => ../../dev/passport/
