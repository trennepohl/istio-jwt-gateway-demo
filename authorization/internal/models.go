package internal

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailIsEmpty          = errors.New("email can't be empty")
	ErrRolesAreEmpty         = errors.New("roles must not be empty")
	ErrWrongAdminCredentials = errors.New("admin login failed")
)

type ServiceConfig struct {
	DefaultAdminEmail    string `envconfig:"ADMIN_EMAIL" default:"admin@istio-auth-poc.io"`
	DefaultAdminPassword string `envconfig:"ADMIN_PASSWORD" default:"password"`
	GoogleStateCode      string `envconfig:"GOOGLE_STATE_CODE" default:"somethingunique"`
	GoogleCallbackURL    string `envconfig:"GOOGLE_CALLBACK_URL" default:"http://authorization.com:8080/auth/callback/google"`
	GoogleClientID       string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret   string `envconfig:"GOOGLE_CLIENT_SECRET"`
}

type DatabaseConfig struct {
	DatabaseUser     string `envconfig:"DATABASE_USER" default:"istio-poc"`
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"localhost"`
	DatabasePort     int    `envconfig:"DATABASE_PORT" default:"5432"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" default:"mysecretpassword"`
	DatabaseName     string `envconfig:"DATABASE_NAME" default:"authorization"`
}

type User struct {
	//ID column to be used in Database
	ID uint64 `json:"user_id,omitempty" gorm:"primaryKey,autoIncrement:true"`

	//Email: unique user identifier
	Email string `json:"email,omitempty" gorm:"unique"`

	// Name: The user's full name.
	Name string `json:"name,omitempty"`

	// Picture: URL of the user's picture image.
	Picture string `json:"picture,omitempty"`

	// VerifiedEmail: Boolean flag which is true if the email address is
	// verified. Always verified because we only return the user's primary
	// email address.
	//
	// default:" true
	VerifiedEmail *bool `json:"verified_email,omitempty"`

	//Password is only used by the admin user
	Password string `json:"-"`

	//Roles that the user holds
	Roles []Role `json:"roles" gorm:"many2many:user_roles"`
}

type Role struct {
	ID   uint64 `gorm:"primaryKey,autoIncrement:true" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type Claims struct {
	Roles []string
	jwt.StandardClaims
}
