package router

type JSON map[string]interface{}

type RoleAssociation struct {
	UserID uint64
	RoleID uint64
}

type Login struct {
	Email    string
	Password string
}
