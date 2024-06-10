package handlers

import (
	"go-skeleton/internal/application/domain/dummy"
	"go-skeleton/internal/application/providers/pagination"
	dummyCreate "go-skeleton/internal/application/services/dummy/CREATE"
	dummyDelete "go-skeleton/internal/application/services/dummy/DELETE"
	dummyEdit "go-skeleton/internal/application/services/dummy/EDIT"
	dummyGet "go-skeleton/internal/application/services/dummy/GET"
	dummyList "go-skeleton/internal/application/services/dummy/LIST"
	dummyRepository "go-skeleton/internal/repositories/dummy"
	"go-skeleton/pkg/idCreator"
	"go-skeleton/pkg/logger"
	"go-skeleton/pkg/registry"
	"go-skeleton/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DummyHandlers struct {
	DummyRepository *dummyRepository.DummyRepository
	reg             *registry.Registry
	logger          *logger.Logger
	idCreator       *idCreator.IdCreator
	validator       *validator.Validator
}

func NewDummyHandlers(reg *registry.Registry) *DummyHandlers {
	return &DummyHandlers{
		reg: reg,
	}
}

func (hs *DummyHandlers) initDeps() {
	hs.DummyRepository = hs.reg.Inject("dummyRepository.DummyRepository").(*dummyRepository.DummyRepository)
	hs.logger = hs.reg.Inject("logger.Logger").(*logger.Logger)
	hs.idCreator = hs.reg.Inject("idCreator.IdCreator").(*idCreator.IdCreator)
	hs.validator = hs.reg.Inject("validator.Validator").(*validator.Validator)
}

func (hs *DummyHandlers) HandleGetDummy(context echo.Context) error {
	hs.initDeps()
	s := dummyGet.NewService(hs.logger, hs.DummyRepository)
	data := new(dummyGet.Data)

	if errors := context.Bind(data); errors != nil {
		s.CustomError(http.StatusBadRequest, errors)
		return context.JSON(s.Error.Status, s.Error)
	}

	s.Execute(
		dummyGet.NewRequest(data),
	)

	response, err := s.GetResponse()
	if err != nil {
		return context.JSON(err.Status, err)
	}
	return context.JSON(http.StatusOK, response)
}

func (hs *DummyHandlers) HandleCreateDummy(context echo.Context) error {
	hs.initDeps()
	s := dummyCreate.NewService(hs.logger, hs.DummyRepository, hs.idCreator)
	data := new(dummyCreate.Data)

	if errors := context.Bind(data); errors != nil {
		s.CustomError(http.StatusBadRequest, errors)
		return context.JSON(s.Error.Status, s.Error)
	}

	s.Execute(
		dummyCreate.NewRequest(data, hs.validator),
	)

	response, err := s.GetResponse()
	if err != nil {
		return context.JSON(err.Status, err)
	}
	return context.JSON(http.StatusCreated, response)
}

func (hs *DummyHandlers) HandleEditDummy(context echo.Context) error {
	hs.initDeps()
	s := dummyEdit.NewService(hs.logger, hs.DummyRepository)
	data := new(dummyEdit.Data)

	if errors := context.Bind(data); errors != nil {
		s.CustomError(http.StatusBadRequest, errors)
		return context.JSON(s.Error.Status, s.Error)
	}

	s.Execute(
		dummyEdit.NewRequest(data, hs.validator),
	)

	response, err := s.GetResponse()
	if err != nil {
		return context.JSON(err.Status, err)
	}
	return context.JSON(http.StatusOK, response)
}

func (hs *DummyHandlers) HandleListDummy(context echo.Context) error {
	hs.initDeps()
	s := dummyList.NewService(
		hs.logger,
		hs.DummyRepository,
		pagination.NewPaginationProvider[dummy.Dummy](hs.DummyRepository),
	)

	data := new(dummyList.Data)
	bindErr := echo.QueryParamsBinder(context).
		Int("page", &data.Page).
		BindErrors()

	if bindErr != nil {
		s.CustomError(http.StatusBadRequest, bindErr)
		return context.JSON(http.StatusBadRequest, s.Error)
	}

	s.Execute(
		dummyList.NewRequest(data),
	)

	response, err := s.GetResponse()
	if err != nil {
		return context.JSON(err.Status, err)
	}
	return context.JSON(http.StatusOK, response)
}

func (hs *DummyHandlers) HandleDeleteDummy(context echo.Context) error {
	hs.initDeps()
	s := dummyDelete.NewService(hs.logger, hs.DummyRepository)
	data := new(dummyDelete.Data)

	if errors := context.Bind(data); errors != nil {
		s.CustomError(http.StatusBadRequest, errors)
		return context.JSON(s.Error.Status, s.Error)
	}

	s.Execute(
		dummyDelete.NewRequest(data),
	)

	response, err := s.GetResponse()
	if err != nil {
		return context.JSON(err.Status, err)
	}
	return context.JSON(http.StatusOK, response)
}
