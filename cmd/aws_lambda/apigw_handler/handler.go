package apigwhandler

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
	"go-skeleton/cmd/aws_lambda/apigw_handler/apiGwAdapter"
	"net/http"
)

type APIGWHandler struct {
	server *echo.Echo
}

func NewAPIGWHandler(s *echo.Echo) *APIGWHandler {
	return &APIGWHandler{
		server: s,
	}
}

func (h *APIGWHandler) Handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	requestHttp, err := AdaptAPIGatewayRequest(req)
	if err != nil {
		return nil, err
	}
	writer := &apiGwAdapter.ResponseWriter{}
	h.server.ServeHTTP(writer, requestHttp)
	return &events.APIGatewayV2HTTPResponse{
		StatusCode:      writer.Status,
		Headers:         req.Headers,
		Body:            string(writer.Body),
		IsBase64Encoded: false,
		Cookies:         req.Cookies,
	}, err
}

func AdaptAPIGatewayRequest(req events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	var body []byte
	if req.Body != "" {
		if req.IsBase64Encoded {
			decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				return nil, err
			}
			body = decodedBody
		} else {
			body = []byte(req.Body)
		}
	}

	// Create http.Request
	httpRequest, err := http.NewRequest(req.RequestContext.HTTP.Method, req.RequestContext.HTTP.Path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for key, value := range req.PathParameters {
		httpRequest.SetPathValue(key, value)
	}

	for key, value := range req.Headers {
		httpRequest.Header.Set(key, value)
	}

	query := httpRequest.URL.Query()
	for key, value := range req.QueryStringParameters {
		query.Add(key, value)
	}
	httpRequest.URL.RawQuery = query.Encode()
	query.Encode()
	return httpRequest, nil
}
