package main

import (
	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"io/ioutil"
	"net/http"
	"stress-tool/comon/util"
	"stress-tool/src/controller"
)

const errMsgKey  = "A mistake seems to have happened, cause"
var log *go_logger.Logger

// 加载包中的方法
var loadHttpHandleMethodMap map[string]func([]byte)interface{}
func init() {
	log = util.GetLogger()
	loadHttpHandleMethodMap = make(map[string]func([]byte)interface{}, 10)
	// 注册 Controller 中的方法
	util.MapMerge(loadHttpHandleMethodMap, controller.TcpControllerMethodHandleMap)
}

func loadHttpInterface(router *gin.Engine) {
	for k, method := range loadHttpHandleMethodMap {
		v := method // method 会变，所以在下面的匿名函数中不能直接使用method，否则所有路径都会指向同一个方法
		router.POST(k, func(c *gin.Context) {
			closer := c.Request.Body
			body, e := ioutil.ReadAll(closer)
			if e != nil {
				log.Errorf(e.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					errMsgKey: "s",
				})
			}
			// 执行 handle 方法
			response := v(body)
			c.JSON(http.StatusOK, response)
			log.Debugf("path:[%s]\nrequest param:[%s]\nresponse:[%+v]", k, string(body), response)
		})
	}
}

func main() {
	router := gin.Default()
	// 静态资源
	router.StaticFS("/web-template/lib/", http.Dir("web-template/lib/"))
	router.StaticFS("/web-template/css/", http.Dir("web-template/css/"))
	router.StaticFS("/web-template/js/",  http.Dir("web-template/js/"))

	// 界面
	router.LoadHTMLGlob("web-template/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	router.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	// 批量监听 Controller 中的方法
	loadHttpInterface(router)

	router.Run(":8080")
}