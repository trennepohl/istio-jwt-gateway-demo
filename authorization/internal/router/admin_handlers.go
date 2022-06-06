package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/trennepohl/istio-auth-poc/authorization/internal"
)

func (a *authorizationRouter) removeRole(context echo.Context) error {
	reqBody := &RoleAssociation{}
	if err := context.Bind(reqBody); err != nil {
		return context.JSON(http.StatusBadRequest, JSON{"error": err})
	}

	if err := a.authorization.RemoveRole(reqBody.UserID, reqBody.RoleID); err != nil {
		return context.JSON(http.StatusInternalServerError, JSON{"error": err})
	}

	response := JSON{"message": fmt.Sprintf("Role %d removed from user %d", reqBody.RoleID, reqBody.UserID)}
	return context.JSON(http.StatusOK, response)
}

func (a *authorizationRouter) assignRole(context echo.Context) error {
	reqBody := &RoleAssociation{}
	if err := context.Bind(reqBody); err != nil {
		return context.JSON(http.StatusBadRequest, JSON{"error": err})
	}

	if err := a.authorization.AssignRole(reqBody.UserID, reqBody.RoleID); err != nil {
		return context.JSON(http.StatusInternalServerError, JSON{"error": err})
	}

	response := JSON{"message": fmt.Sprintf("Role %d added to user %d", reqBody.RoleID, reqBody.UserID)}
	return context.JSON(http.StatusOK, response)
}

func (a *authorizationRouter) createRole(context echo.Context) error {
	role := &internal.Role{}
	err := context.Bind(role)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	err = a.authorization.CreateRole(role.Name)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}

	response := JSON{"message": fmt.Sprintf("Role %s was created", role.Name)}
	return context.JSON(http.StatusOK, response)
}

func (a *authorizationRouter) listRoles(context echo.Context) error {
	roles, err := a.authorization.ListRoles()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, JSON{"error": err})
	}

	response := JSON{"roles": roles}
	return context.JSON(http.StatusOK, response)
}

func (a *authorizationRouter) listUsers(context echo.Context) error {
	users, err := a.authorization.ListUsers()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, JSON{"error": err})
	}

	fmt.Println(users)
	response := JSON{"users": users}
	return context.JSON(http.StatusOK, response)
}
