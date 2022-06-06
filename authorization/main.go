package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/trennepohl/istio-auth-poc/authorization/internal"
	"github.com/trennepohl/istio-auth-poc/authorization/internal/authorization"
	"github.com/trennepohl/istio-auth-poc/authorization/internal/gauth"
	database "github.com/trennepohl/istio-auth-poc/authorization/internal/postgres"
	"github.com/trennepohl/istio-auth-poc/authorization/internal/router"
)

func main() {
	databaseConfig := &internal.DatabaseConfig{}
	serviceConfig := &internal.ServiceConfig{}

	if err := loadEnvironments(databaseConfig, serviceConfig); err != nil {
		log.Fatalln(err)
	}

	connection, err := database.New(databaseConfig)
	if err != nil {
		log.Fatalln(err)
	}

	authentication := gauth.NewGoogleAuth(serviceConfig)

	privateKey, publicKey := loadKeys()

	authorizationConfig := &authorization.Config{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		AdminEmail:    serviceConfig.DefaultAdminEmail,
		AdminPassword: serviceConfig.DefaultAdminPassword,
	}

	authorizationService, err := authorization.New(connection, authorizationConfig)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.Use()

	authRouter := router.New(authorizationService, authentication, e)
	authRouter.Serve()
}

func loadEnvironments(dbConfig *internal.DatabaseConfig, serviceConfig *internal.ServiceConfig) error {
	err := envconfig.Process("auth", dbConfig)
	if err != nil {
		return err
	}

	return envconfig.Process("auth", serviceConfig)
}

func loadKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	prk, err := ioutil.ReadFile("/data/private-key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prk)
	if err != nil {
		log.Fatalln(err)
	}

	pbk, err := ioutil.ReadFile("/data/public-key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pbk)
	if err != nil {
		log.Fatalln(err)
	}

	return privateKey, pubKey
}
