package internal

type AuthorizationRouter interface {
	Serve()
}

type UserAuthorizationService interface {
	ListUsers() ([]User, error)
}

type RolesAuthorizationService interface {
	ListRoles() ([]Role, error)
	AddRole(role string, userEmail string) error
	RemoveRole(userID uint64, roleID uint64) error
	AssignRole(userID uint64, roleID uint64) error
}

type AuthorizationService interface {
	UserAuthorizationService
	RolesAuthorizationService
	BasicLogin(email, password string) (string, error)
	CreateRole(role string) error
	GetJWT(user User) (string, error)
	GetJWK() map[string]interface{}
}

type AuthenticationService interface {
	Login() string
	GetUserInfo(code string, state string) (User, error)
	GetRedirectUri() string
}

type UserStorage interface {
	CreateUser(user *User, options ...UserOption) error
	GetUsers() ([]User, error)
	GetUser(email string) (User, error)
	UpdateUser(user User) error
}

type RolesStorage interface {
	GetRoles() ([]Role, error)
	CreateRole(role Role) error
	RemoveRole(roleID uint64, userID uint64) error
	AssignRole(roleID uint64, userID uint64) error
}

type Database interface {
	UserStorage
	RolesStorage
	Migrate(...interface{}) error
}
type UserOption func(*User)

func WithDefaultRoles(user *User) {
	user.Roles = append(user.Roles, Role{Name: "ReadOnly"})
}
