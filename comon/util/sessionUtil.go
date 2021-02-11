package util

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
)

const secretKey = "stress_tool_session_secret_key_20210209180757"
const sessionName = "stress_tool_session_name"

var log *go_logger.Logger

func InitSession(engine *gin.Engine) {
	// 创建基于cookie的存储引擎，secretKey 参数是用于加密的密钥
	store := cookie.NewStore([]byte(secretKey))
	// 设置 session 中间件，参数 sessionName，指的是 session 的名字，也是 cookie 的名字
	// store 是前面创建的存储引擎，我们可以替换成其他存储引擎
	engine.Use(sessions.Sessions(sessionName, store))
	log = GetLogger()
}

func SetSession(c *gin.Context, key, value interface{}) error {
	session := sessions.Default(c)
	session.Set(key, value)
	err := session.Save()
	return WrapError("save session error", err)
}

func GetSession(c *gin.Context, key interface{}) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

func DeleteSession(c *gin.Context, key interface{}) error {
	session := sessions.Default(c)
	session.Delete(key)
	err := session.Save()
	return WrapError("save session error", err)
}
