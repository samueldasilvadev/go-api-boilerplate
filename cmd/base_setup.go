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

	dbFunc := func() registry.Dependency {
		return db
	}

	idCFunc := func() registry.Dependency {
		return idCreator.NewIdCreator()
	}

	valFunc := func() registry.Dependency {
		val := validator.NewValidator()
		val.Boot()
		return val
	}

	confFunc := func() registry.Dependency {
		return conf
	}

	Reg = registry.NewRegistry()

	Reg.Provide("logger.Logger", lFunc)
	Reg.Provide("validator.Validator", valFunc)
	Reg.Provide("config.Config", confFunc)
	Reg.Provide("idCreator.IdCreator", idCFunc)
	Reg.Provide("database.MySql", dbFunc)

	dummyRepositoryFunc := func() registry.Dependency {
		return dummyRepository.NewDummyRepository(Reg.Inject("database.MySql").(*database.MySql).Db)
	}

	Reg.Provide("dummyRepository.DummyRepository", dummyRepositoryFunc)
	//{{codeGen6}}
}
