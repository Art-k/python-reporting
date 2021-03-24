package include

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func ApiProcessing() {

	r := gin.Default()

	//r.LoadHTMLGlob("results/*.html")
	r.GET("/report/:report_id/:recipient_id", GetReportByRecipient) // return file, open counter increased, download statistic is saved for recipient <> report
	r.POST("/job_done/:job_id", FinishingTask)
	r.GET("/logo/:id", GETLogoProcessing)

	auth := r.Group("/")
	auth.Use(TokenAuthMiddleware())
	{

		auth.POST("/files", UploadFiles)
		auth.GET("/scripts", GetBaseScripts)

		//auth.MaxMultipartMemory = 8 << 20
		auth.POST("/script", AddBaseScript)
		auth.POST("/script/:script_hash/files", UploadScriptFiles)
		auth.POST("/script/:script_hash/parameters", SetScriptParameters)
		auth.POST("/script/:script_hash/task", SetTask)

		auth.GET("/tasks", GetTask)
		auth.PATCH("/task/:task_id", PatchTask)
		auth.POST("/task/:task_id/parameters", PostTaskParameter)
		auth.GET("/task/:task_id/parameters", GetTaskParameter)
		auth.POST("/task/:task_id/recipients", PostTaskRecipients)
		auth.GET("/task/:task_id/recipients", GETTaskRecipients)
		auth.POST("/task/:task_id/schedule", PostTaskSchedule)

		auth.POST("/task/:task_id/run", RunTask)
		auth.POST("/task/:task_id/enable", EnableTask)

		auth.GET("/report/:report_id", GetReport)          //return file, open counter increased
		auth.GET("/report-info/:report_id", GetReportInfo) //return json

		auth.GET("/jobs", GetJobs)
		auth.GET("/schedule", GetSchedule)
	}

	//TODO add auto update for a server
	//TODO add error page, not found, bad gateway, etc

	r.Run(":" + os.Getenv("PORT"))
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

//
//func Auth(c *gin.Context) gin.HandlerFunc{
//
//	requiredToken := os.Getenv("API_TOKEN")
//
//	// We want to make sure the token is set, bail if not
//	if requiredToken == "" {
//		Log.Fatal("Please set API_TOKEN environment variable")
//	}
//
//	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
//
//	if token == "" {
//		respondWithError(c, 401, "API token required")
//		return
//	}
//
//	if token != requiredToken {
//		respondWithError(c, 401, "Invalid API token")
//		return
//	}
//
//	c.Next()
//}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	// We want to make sure the token is set, bail if not
	if requiredToken == "" {
		Log.Fatal("Please set API_TOKEN environment variable")
	}

	return func(c *gin.Context) {
		//token := c.Request.FormValue("api_token")
		token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]

		if token == "" {
			respondWithError(c, 401, "API token required")
			return
		}

		if token != requiredToken {
			respondWithError(c, 401, "Invalid API token")
			return
		}

		c.Next()
	}
}
