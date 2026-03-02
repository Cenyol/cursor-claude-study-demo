package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 JSON：{"code": int, "data": any, "msg": string}
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

const (
	CodeSuccess  = 0
	CodeBadReq   = 400
	CodeUnauth   = 401
	CodeConflict = 409
	CodeServer   = 500
)

func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Data: data, Msg: "success"})
}

func JSONBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Code: CodeBadReq, Msg: msg})
}

func JSONUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{Code: CodeUnauth, Msg: msg})
}

func JSONConflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, Response{Code: CodeConflict, Msg: msg})
}

func JSONServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{Code: CodeServer, Msg: msg})
}
