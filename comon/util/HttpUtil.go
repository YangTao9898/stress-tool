package util

import "github.com/gin-gonic/gin"

const (
	RESULT_OK  = "0000"
	RESULT_ERR = "500"
)

func ResponsePack(resultCode, errMsg string, responseData interface{}) interface{} {
	return gin.H{
		"resultCode": resultCode,
		"errMsg": errMsg,
		"data": responseData,
	}
}

func ResponseSuccPack(responseData interface{}) interface{} {
	return gin.H{
		"resultCode": RESULT_OK,
		"errMsg": "",
		"data": responseData,
	}
}

func ResponseFailPack(errMsg string) interface{} {
	return gin.H{
		"resultCode": RESULT_ERR,
		"errMsg": errMsg,
		"data": nil,
	}
}
