package include

import (
	"github.com/gin-gonic/gin"
	"os"
)

func ApiProcessing() {

	r := gin.Default()

	r.POST("/files", UploadFiles)
	r.GET("/logo/:id", GETLogoProcessing)
	r.GET("/scripts", GetBaseScripts)

	r.MaxMultipartMemory = 8 << 20
	r.POST("/script", AddBaseScript)
	r.POST("/script/:script_hash/files", UploadScriptFiles)
	r.POST("/script/:script_hash/parameters", SetScriptParameters)
	r.POST("/script/:script_hash/task", SetTask)

	r.GET("/tasks", GetTask)
	r.PATCH("/task/:task_id", PatchTask)
	r.POST("/task/:task_id/parameters", PostTaskParameter)
	r.POST("/task/:task_id/recipients", PostTaskRecipients)
	r.GET("/task/:task_id/recipients", GETTaskRecipients)
	r.POST("/task/:task_id/schedule", PostTaskSchedule)

	r.POST("/task/:task_id/run", RunTask)
	r.POST("/task/:task_id/enable", EnableTask)

	r.POST("/job_done/:job_id", FinishingTask)

	//r.LoadHTMLGlob("results/*.html")
	r.GET("/report/:report_id", GetReport)                          //return file, open counter increased
	r.GET("/report-info/:report_id", GetReportInfo)                 //return json
	r.GET("/report/:report_id/:recipient_id", GetReportByRecipient) // return file, open counter increased, download statistic is saved for recipient <> report

	r.GET("/jobs", GetJobs)
	r.GET("/schedule", GetSchedule)

	//TODO add auto update for a server
	//TODO add error page, not found, bad gateway, etc

	r.Run(":" + os.Getenv("PORT"))
}
