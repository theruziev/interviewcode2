module github.com/theruziev/oson_auth

go 1.20

require (
	github.com/Masterminds/squirrel v1.5.3
	github.com/alecthomas/kong v0.6.1
	github.com/georgysavva/scany/v2 v2.0.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-pkgz/repeater v1.1.3
	github.com/go-pkgz/requester v0.0.4
	github.com/go-playground/validator/v10 v10.11.1
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/uuid v1.3.0
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa
	github.com/jackc/pgx/v5 v5.2.0
	github.com/joho/godotenv v1.4.0
	github.com/json-iterator/go v1.1.12
	github.com/mailgun/mailgun-go/v4 v4.8.1
	github.com/pquerna/otp v1.3.0
	github.com/stretchr/testify v1.8.0
	github.com/theruziev/oson_auth/pkg/events v0.0.0-00010101000000-000000000000
	github.com/wagslane/go-rabbitmq v0.11.0
	go.uber.org/zap v1.23.0
	golang.org/x/crypto v0.0.0-20221012134737-56aed061732a
	golang.org/x/net v0.3.0
	golang.org/x/sync v0.0.0-20220923202941-7f9b1623fab7

)

require (
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-jet/jet/v2 v2.9.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/puddle/v2 v2.1.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rabbitmq/amqp091-go v1.3.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/theruziev/oson_auth/pkg/events => ./pkg/events
