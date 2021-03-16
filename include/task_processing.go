package include

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"net/http"
)

func GetTask(cnt *gin.Context) {

	var tasks []DBTask
	page, perPage := setPaginationParameters(cnt)

	format := cnt.Query("format")

	db.Preload(clause.Associations).
		Order("created_at desc").
		Find(&tasks).
		Limit(perPage).
		Offset(page - 1*perPage)

	switch format {

	default:
		cnt.JSON(http.StatusOK, tasks)
	}

}

func PostTaskSchedule(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}

	var incSchedule POSTSchedule
	err = json.Unmarshal(jsonData, &incSchedule)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	task.POSTSchedule = incSchedule
	db.Save(&task)

	for _, job := range Sch.Jobs() {
		for _, tag := range job.Tags() {
			if tag == task.ID {
				err := Sch.RemoveByTag(task.ID)
				if err != nil {
					Log.Error(err)
				}
				RunScheduler(&task)
			}
		}
	}

	cnt.JSON(http.StatusCreated, task)

}

func PatchTask(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}

	var patchedTask POSTTask
	err = json.Unmarshal(jsonData, &patchedTask)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	task.POSTTask = patchedTask

	db.Save(&task)

	cnt.JSON(http.StatusAccepted, task)

}

func PostTaskRecipients(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}

	var iRecipients []POSTRecipient
	err = json.Unmarshal(jsonData, &iRecipients)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var recipients []DBRecipient

	db.Where("db_task_id = ?", task.ID).Delete(&recipients)
	recipients = nil

	for _, iR := range iRecipients {
		recipients = append(recipients, DBRecipient{
			DBTaskID:      task.ID,
			POSTRecipient: iR,
		})
	}

	db.Create(&recipients)

	cnt.JSON(http.StatusCreated, recipients)

}

func PostTaskParameter(cnt *gin.Context) {

	taskId := cnt.Param("task_id")

	var task DBTask
	err := db.Where("id =?", taskId).Find(&task).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}

	var incomingParameters []POSTTaskParameter
	err = json.Unmarshal(jsonData, &incomingParameters)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var parameters []DBTaskParameter

	db.Where("db_task_id = ?", task.ID).Delete(&parameters)
	parameters = nil

	for _, incParam := range incomingParameters {
		parameters = append(parameters, DBTaskParameter{
			DBTaskID:          task.ID,
			POSTTaskParameter: incParam,
		})
	}

	db.Create(&parameters)

	cnt.JSON(http.StatusCreated, parameters)

}

func SetTask(cnt *gin.Context) {

	baseScriptHash := cnt.Param("script_hash")

	var baseScript DBBaseScript
	err := db.Where("id =?", baseScriptHash).Find(&baseScript).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}

	var newTask POSTTask
	err = json.Unmarshal(jsonData, &newTask)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var dbNewTask DBTask
	dbNewTask.DBBaseScriptID = baseScript.ID
	dbNewTask.POSTTask = newTask
	db.Create(&dbNewTask)

	db.Preload(clause.Associations).Where("id = ?", dbNewTask.ID).Find(&dbNewTask)

	cnt.JSON(http.StatusCreated, dbNewTask)

}
