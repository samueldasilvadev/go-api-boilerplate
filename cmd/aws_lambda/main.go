package main

import (
	apigwhandler "go-skeleton/cmd/aws_lambda/apigw_handler"
	"go-skeleton/cmd/http/server"
	dummyRepository "go-skeleton/internal/repositories/dummy"
	"go-skeleton/pkg/config"
	"go-skeleton/pkg/database"
	"go-skeleton/pkg/idCreator"
	"go-skeleton/pkg/logger"
	"go-skeleton/pkg/registry"
	"go-skeleton/pkg/validator"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	reg *registry.Registry
)

func main() {
	setup()
	echoserver := server.NewServer(reg)
	handler := apigwhandler.NewAPIGWHandler(echoserver.Setup())
	lambda.Start(handler.Handler)
}

func setup() {
	conf := config.NewConfig()
	err := conf.LoadEnvs()
	if err != nil {
		panic(err)
	}

	l := logger.NewLogger(
		conf.ReadConfig("ENVIRONMENT"),
		conf.ReadConfig("APP"),
		conf.ReadConfig("VERSION"),
	)

	l.Boot()

	db := database.NewMysql(
		l,
		conf.ReadConfig("DB_USER"),
		conf.ReadConfig("DB_PASS"),
		conf.ReadConfig("DB_URL"),
		conf.ReadConfig("DB_PORT"),
		conf.ReadConfig("DB_DATABASE"),
	)

	idC := idCreator.NewIdCreator()
	val := validator.NewValidator()

	db.Connect()
	val.Boot()

	reg = registry.NewRegistry()
	reg.Provide("logger", l)
	reg.Provide("validator", val)
	reg.Provide("config", conf)
	reg.Provide("idCreator", idC)

	reg.Provide("dummyRepository", dummyRepository.NewDummyRepository(db.Db))
	//{{codeGen6}}
}
