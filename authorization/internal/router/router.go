package router

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/trennepohl/istio-auth-poc/authorization/internal"
)

type authorizationRouter struct {
	engine         *echo.Echo
	authorization  internal.AuthorizationService
	authentication internal.AuthenticationService
}

func (a *authorizationRouter) Serve() {
	if err := a.engine.Start("0.0.0.0:4000"); err != nil {
		log.Fatalln(err)
	}
}

func New(authorization internal.AuthorizationService, authentication internal.AuthenticationService, r *echo.Echo) internal.AuthorizationRouter {
	router := &authorizationRouter{
		engine:         r,
		authorization:  authorization,
		authentication: authentication,
	}

	router.engine.GET("/jwk", router.handleJWK)
	router.engine.GET("/", router.handleMain)

	admin := router.engine.Group("/admin")
	admin.GET("/users", router.listUsers)
	admin.GET("/roles", router.listRoles)
	admin.POST("/role", router.createRole)
	admin.POST("/role/assign", router.assignRole)
	admin.DELETE("/role/delete", router.removeRole)

	auth := router.engine.Group("/auth")
	auth.GET("/callback/google", router.handleGoogleCallback)
	auth.GET("/google", router.handleGoogleLogin)
	auth.POST("/login", router.basicAuthLogin)

	return &authorizationRouter{
		engine: r,
	}
}

func (a *authorizationRouter) handleMain(ctx echo.Context) error {
	var htmlIndex = `<html><body><a href="/auth/google">Google Log In</a></body></html>`
	return ctx.HTML(http.StatusOK, htmlIndex)
}

func (a *authorizationRouter) handleGoogleCallback(ctx echo.Context) error {
	state := ctx.FormValue("state")
	code := ctx.FormValue("code")

	//TODO: validate state with state string
	if state == "" || code == "" {
		return ctx.JSON(http.StatusUnauthorized, &map[string]interface{}{})
	}

	user, err := a.authentication.GetUserInfo(code, state)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "failed to fetch userInfo")
	}

	token, err := a.authorization.GetJWT(user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, JSON{"token": token})
}

func (a *authorizationRouter) handleGoogleLogin(ctx echo.Context) error {
	url := a.authentication.Login()
	return ctx.Redirect(http.StatusMovedPermanently, url)
}

func (a *authorizationRouter) handleJWK(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, a.authorization.GetJWK())
}

func (a *authorizationRouter) basicAuthLogin(context echo.Context) error {
	credentials := &Login{}
	if err := context.Bind(credentials); err != nil {
		return context.JSON(http.StatusBadRequest, JSON{"error": err})
	}

	token, err := a.authorization.BasicLogin(credentials.Email, credentials.Password)
	if err == nil {
		return context.JSON(http.StatusOK, JSON{"token": token})
	}

	if errors.Is(err, internal.ErrWrongAdminCredentials) {
		return context.JSON(http.StatusUnauthorized, JSON{"error": "wrong credentials"})
	}

	return context.JSON(http.StatusInternalServerError, JSON{"error": err})
}
