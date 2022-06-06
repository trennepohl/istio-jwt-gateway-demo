package authorization

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/trennepohl/istio-auth-poc/authorization/internal"
)

type Config struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	AdminEmail    string
	AdminPassword string
}

type authorizationService struct {
	db           internal.Database
	jwks         jwk.Set
	pubKey       *rsa.PublicKey
	privateKey   *rsa.PrivateKey
	jwkPublicKey jwk.RSAPublicKey
}

func New(db internal.Database, config *Config) (internal.AuthorizationService, error) {
	svc := &authorizationService{
		db:           db,
		privateKey:   config.PrivateKey,
		pubKey:       config.PublicKey,
		jwks:         jwk.NewSet(),
		jwkPublicKey: jwk.NewRSAPublicKey(),
	}

	if err := svc.jwkPublicKey.FromRaw(svc.pubKey); err != nil {
		return svc, err
	}

	if err := jwk.AssignKeyID(svc.jwkPublicKey); err != nil {
		return svc, err
	}

	svc.jwks.Add(svc.jwkPublicKey)

	if err := db.Migrate(internal.User{}, internal.Role{}); err != nil {
		return svc, err
	}

	err := svc.createAdminUser(config.AdminEmail, config.AdminPassword)
	return svc, err
}

func (a *authorizationService) createAdminUser(email, password string) error {
	pH := md5.Sum([]byte(password))
	password = hex.EncodeToString(pH[:])
	_ = a.db.CreateUser(&internal.User{Email: email, Password: password, Roles: []internal.Role{{Name: "Admin"}}})
	return nil
}

func (a *authorizationService) BasicLogin(email, password string) (string, error) {
	user, err := a.db.GetUser(email)
	if err != nil {
		return "", err
	}

	pH := md5.Sum([]byte(password))
	if user.Password != hex.EncodeToString(pH[:]) {
		return "", internal.ErrWrongAdminCredentials
	}

	return a.signToken(user.Roles)
}

func (a *authorizationService) AddRole(role string, userEmail string) error {
	user, err := a.db.GetUser(userEmail)
	if err != nil {
		return err
	}

	user.Roles = append(user.Roles, internal.Role{Name: role})
	return a.db.UpdateUser(user)
}

func (a *authorizationService) ListUsers() ([]internal.User, error) {
	return a.db.GetUsers()
}

func (a *authorizationService) ListRoles() ([]internal.Role, error) {
	return a.db.GetRoles()
}

func (a *authorizationService) RemoveRole(userID uint64, roleID uint64) error {
	return a.db.RemoveRole(roleID, userID)
}

func (a *authorizationService) AssignRole(userID uint64, roleID uint64) error {
	return a.db.AssignRole(roleID, userID)
}

func (a *authorizationService) CreateRole(roleName string) error {
	return a.db.CreateRole(internal.Role{Name: roleName})
}

func (a *authorizationService) GetJWT(user internal.User) (string, error) {
	if user.Email == "" {
		return "error", internal.ErrEmailIsEmpty
	}

	dUser, err := a.db.GetUser(user.Email)
	if err == nil {
		return a.signToken(dUser.Roles)
	}

	if errors.Is(err, internal.ErrUserNotFound) {
		return a.createAndSign(user)
	}

	return "", err
}

func (a *authorizationService) GetJWK() map[string]interface{} {
	res := make(map[string]interface{})
	b, _ := json.Marshal(a.jwks)
	_ = json.Unmarshal(b, &res)
	return res
}

func (a *authorizationService) createAndSign(user internal.User) (token string, err error) {
	err = a.db.CreateUser(&user, internal.WithDefaultRoles)
	if err != nil {
		return token, err
	}

	return a.signToken(user.Roles)
}

func (a *authorizationService) signToken(roles []internal.Role) (string, error) {
	if len(roles) == 0 {
		return "", internal.ErrRolesAreEmpty
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	dt := time.Now().UTC()
	claims := internal.Claims{
		Roles: []string{},
		StandardClaims: jwt.StandardClaims{
			Audience:  "poc-ecosystem",
			ExpiresAt: dt.Add(24 * time.Hour).Unix(),
			IssuedAt:  dt.Unix(),
			Issuer:    "istio-auth-poc",
			Subject:   "poc",
		},
	}

	for _, role := range roles {
		claims.Roles = append(claims.Roles, role.Name)
	}

	token.Claims = claims

	token.Header["kid"] = a.jwkPublicKey.KeyID()
	return token.SignedString(a.privateKey)
}
