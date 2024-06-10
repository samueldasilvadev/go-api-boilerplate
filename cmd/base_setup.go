package cmd

import (
	dummyRepository "go-skeleton/internal/repositories/dummy"
	//{{codeGen5}}
	"go-skeleton/pkg/config"
	"go-skeleton/pkg/database"
	"go-skeleton/pkg/idCreator"
	"go-skeleton/pkg/logger"
	"go-skeleton/pkg/registry"
	"go-skeleton/pkg/validator"
)

var (
	Reg       *registry.Registry
	ApiPrefix string
)

func Setup() {
	conf := config.NewConfig()
	err := conf.LoadEnvs()
	if err != nil {
		panic(err)
	}

	db := database.NewMysql(
		Reg.Inject("logger.Logger").(*logger.Logger),
		conf.ReadConfig("DB_USER"),
		conf.ReadConfig("DB_PASS"),
		conf.ReadConfig("DB_URL"),
		conf.ReadConfig("DB_PORT"),
		conf.ReadConfig("DB_DATABASE"),
	)
	db.Connect()

	ApiPrefix = conf.ReadConfig("API_PREFIX")

	lFunc := func() registry.Dependency {
		l := logger.NewLogger(
			conf.ReadConfig("ENVIRONMENT"),
			conf.ReadConfig("APP"),
			conf.ReadConfig("VERSION"),
		)
		l.Boot()
		return l
	}

	idCFunc := func() registry.Dependency {
		return idCreator.NewIdCreator()
	}

	valFunc := func() registry.Dependency {
		val := validator.NewValidator()
		val.Boot()
		return val
	}

	Reg = registry.NewRegistry()

	Reg.OnDemandProvide("logger.Logger", lFunc)
	Reg.OnDemandProvide("validator.Validator", valFunc)
	Reg.Provide("config.Config", conf)
	Reg.OnDemandProvide("idCreator.IdCreator", idCFunc)
	Reg.Provide("database.MySql", db)

	dummyRepositoryFunc := func() registry.Dependency {
		return dummyRepository.NewDummyRepository(db.Db)
	}

	Reg.OnDemandProvide("dummyRepository.DummyRepository", dummyRepositoryFunc)
	//{{codeGen6}}
}
