module github.com/trennepohl/istio-auth-poc/authorization

go 1.16

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lestrrat-go/jwx v1.2.1
	golang.org/x/oauth2 v0.0.0-20210622215436-a8dc77f794b6
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)
