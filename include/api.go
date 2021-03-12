package include

import "github.com/gin-gonic/gin"

func ApiProcessing() {

	r := gin.Default()

	r.GET("/scripts", GetBaseScripts)

	r.Run(":49999")
}
