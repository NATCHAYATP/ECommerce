package servers

import (
	"github.com/NATCHAYATP/E-Commerce/modules/middlewares/middlewaresHandlers"
	"github.com/NATCHAYATP/E-Commerce/modules/middlewares/middlewaresRepositories"
	"github.com/NATCHAYATP/E-Commerce/modules/middlewares/middlewaresUsecases"
	"github.com/NATCHAYATP/E-Commerce/modules/users/usersHandlers"
	"github.com/NATCHAYATP/E-Commerce/modules/users/usersRepositories"
	"github.com/NATCHAYATP/E-Commerce/modules/users/usersUsecases"
	"github.com/NATCHAYATP/E-Commerce/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func initMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.Middlewaresrepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")
	router.Post("/signup", handler.SignUpCustomer)
}
