package include

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func EnableTask(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Preload(clause.Associations).Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	task.Enabled = true
	db.Save(&task)

	RunScheduler(&task)

	cnt.JSON(http.StatusAccepted, task)
}

func RunTask(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Preload(clause.Associations).Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	job := StartJob(&task, "api")

	cnt.JSON(http.StatusCreated, job)

}

func StartJob(task *DBTask, source string) DBJob {

	Log.Trace("Start Job '", source, "'")

	var scriptFile DBScriptFile
	db.Where("db_base_script_id = ?", task.DBBaseScriptID).Where("script_file = ?", true).Find(&scriptFile)

	var parameters []DBTaskParameter
	db.Where("db_task_id = ?", task.ID).Find(&parameters)

	var job DBJob
	job = DBJob{
		DBBaseScriptID: task.DBBaseScriptID,
		DBTaskID:       task.ID,
		DBScriptFileID: scriptFile.ID,
		Source:         source,
		CommandString:  "",
		CommandOutput:  "",
		Error:          "",
	}
	db.Create(&job)

	// "/home/art-k/PROJECT/MY/python-reporting" +
	//args := strings.Replace(scriptFile.PathToFile, "./", "", 1) +
	//	" --task_id \"" + job.ID + "\"" +
	//	" --call_back_url \"http://127.0.0.1:49999/done\""

	cmd := exec.Command("python3")
	cmd.Args = append(cmd.Args, strings.Replace(scriptFile.PathToFile, "./", "", 1))
	cmd.Args = append(cmd.Args, "--task_id")
	cmd.Args = append(cmd.Args, job.ID)

	for _, param := range parameters {
		cmd.Args = append(cmd.Args, "--"+param.ParameterName)
		cmd.Args = append(cmd.Args, param.ParameterValue)
	}

	job.CommandString = strings.Join(cmd.Args[:], " ")

	db.Save(&job)

	go func(j *DBJob) {
		out, err := cmd.Output()
		if err != nil {
			j.Error = err.Error()
		}
		j.CommandOutput = string(out)
		db.Save(&j)
	}(&job)

	return job
}

func FinishingTask(cnt *gin.Context) {

	jobId := cnt.Param("job_id")
	var job DBJob
	err := db.Where("id = ?", jobId).Find(&job).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
	}
	var postJobDone POSTJobDone
	err = json.Unmarshal(jsonData, &postJobDone)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	job.DurationMs = postJobDone.DurationMs
	db.Save(&job)

	var dbFiles []DBReport
	for _, file := range postJobDone.Files {
		dbFiles = append(dbFiles, DBReport{
			DBJobID:  job.ID,
			FileName: file,
		})
	}
	db.Create(&dbFiles)

	cnt.JSON(http.StatusAccepted, job)

	var task DBTask
	db.Where("id = ?", job.DBTaskID).Find(&task)
	var recipients []DBRecipient
	db.Where("db_task_id = ?", task.ID).Find(&recipients)

	for _, recipient := range recipients {

		fileBlock := "<p>Here is a list of reports :<ul>"
		for _, file := range dbFiles {
			fExt := filepath.Ext(file.FileName)
			fileBlock += "<li><a href='" + os.Getenv("DOMAIN") + "/report/" + file.ID + "/" + recipient.ID + "'>report" + fExt + "</a></li>"
		}
		fileBlock += "</ul>"

		msg := strings.Replace(task.Message, "[[RECIPIENT_NAME]]", recipient.Name, 1)
		msg = strings.Replace(msg, "[[REPORTS]]", fileBlock, 1)

		_, msgId, _ := SendEmailOAUTH2(recipient.Email, task.Subject, msg)

		var outMsg DBOutgoingMails
		db.Where("id = ?", msgId).Find(&outMsg)
		outMsg.DBJobID = job.ID
		db.Save(&outMsg)
	}

}

func GetSchedule(cnt *gin.Context) {

	type jobStatus struct {
		LastRun       time.Time
		NextRun       time.Time
		RunCount      int
		ScheduledTime time.Time
		Tag           []string
		//Error         string
	}

	var jobs []jobStatus
	for _, jb := range Sch.Jobs() {
		jobs = append(jobs, jobStatus{
			LastRun:       jb.LastRun(),
			NextRun:       jb.NextRun(),
			RunCount:      jb.RunCount(),
			ScheduledTime: jb.ScheduledTime(),
			Tag:           jb.Tags(),
			//Error:         jb.Error().Error(),
		})
	}

	cnt.JSON(http.StatusOK, gin.H{"job_count": len(Sch.Jobs()), "jobs": jobs})

}

func GetJobs(cnt *gin.Context) {

	var jobs []DBJob

	page, perPage := setPaginationParameters(cnt)

	format := cnt.Query("format")

	DB := db

	DB.Preload(clause.Associations).
		Order("created_at desc").
		Find(&jobs).
		Limit(perPage).
		Offset(page - 1*perPage)

	switch format {

	default:
		cnt.JSON(http.StatusOK, jobs)
	}

}

func GetReportByRecipient(cnt *gin.Context) {

	reportId := cnt.Param("report_id")
	recipientId := cnt.Param("recipient_id")

	var report DBReport
	err := db.Where("id = ?", reportId).Find(&report).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	var recipient DBRecipient
	err = db.Where("id = ?", recipientId).Find(&recipient).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	http.ServeFile(cnt.Writer, cnt.Request, "results/"+report.FileName)

	var reportDownloadHistory DBReportDownloadHistory
	reportDownloadHistory = DBReportDownloadHistory{
		DBReportID:    reportId,
		DBRecipientID: recipientId,
	}
	db.Create(&reportDownloadHistory)

	report.OpenCount += 1
	db.Save(report)
}
func GetReport(cnt *gin.Context) {

	reportId := cnt.Param("report_id")
	var report DBReport
	err := db.Where("id = ?", reportId).Find(&report).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	fExt := filepath.Ext(report.FileName)
	switch fExt {
	case ".csv":
		http.ServeFile(cnt.Writer, cnt.Request, "results/"+report.FileName)
		break
	case ".html":
		http.ServeFile(cnt.Writer, cnt.Request, "results/"+report.FileName)
		break
	}

	report.OpenCount += 1
	db.Save(report)

}
