package apigwhandler

import (
	"github.com/aws/aws-lambda-go/events"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
)

type APIGWHandler struct {
	server *echo.Echo
}

func NewAPIGWHandler(s *echo.Echo) *APIGWHandler {
	return &APIGWHandler{
		server: s,
	}
}

func (h *APIGWHandler) Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adapter := echoadapter.New(h.server)
	return adapter.Proxy(req)
}
