package apierror

import "github.com/gin-gonic/gin"

type Response struct {
	Message string   `json:"message"`
	Error   string   `json:"error"`
	Status  int      `json:"status"`
	Cause   []string `json:"cause"`
}

func Respond(c *gin.Context, status int, code string, err error) {
	c.AbortWithStatusJSON(status, Response{
		Message: err.Error(),
		Error:   code,
		Status:  status,
		Cause:   []string{},
	})
}
