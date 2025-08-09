package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sidiik/moonpay/api_gateway/internal/dto"
)

func SendResponse(c *gin.Context, statusCode int, message string, data any, err error) {
	status := "success"

	if statusCode >= http.StatusBadRequest {
		status = "error"
	}

	resp := dto.CommonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	c.JSON(statusCode, resp)
}
