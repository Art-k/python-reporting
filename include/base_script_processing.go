package include

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func AddBaseScript(cnt *gin.Context) {

	var newBaseScript POSTBaseScript
	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
	}
	err = json.Unmarshal(jsonData, &newBaseScript)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	var baseScript DBBaseScript
	baseScript.POSTBaseScript = newBaseScript
	db.Create(&baseScript)

	cnt.JSON(http.StatusCreated, &baseScript)

}

func SetScriptParameters(cnt *gin.Context) {

	parametersFor := cnt.Param("script_hash")
	var baseScript DBBaseScript
	err := db.Where("id = ?", parametersFor).Last(&baseScript).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err, "value": parametersFor})
		return
	}

	var parameters []POSTScriptParameter
	jsonData, err := ioutil.ReadAll(cnt.Request.Body)
	if err != nil {
		Log.Error(err)
		return
	}
	err = json.Unmarshal(jsonData, &parameters)
	if err != nil {
		cnt.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var dbParameters []DBScriptParameter
	db.Where("db_base_script_id = ?", baseScript.ID).Delete(&dbParameters)

	dbParameters = nil
	for _, parameter := range parameters {
		var dbParameter DBScriptParameter
		dbParameter.DBBaseScriptID = baseScript.ID
		dbParameter.POSTScriptParameter = parameter
		db.Create(&dbParameter)
	}

	var scriptResponse DBBaseScript
	db.Preload(clause.Associations).
		Where("id = ?", baseScript.ID).
		Last(&scriptResponse)

	cnt.JSON(http.StatusOK, scriptResponse)
}

func UploadScriptFiles(cnt *gin.Context) {

	filesFor := cnt.Param("script_hash")
	var baseScript DBBaseScript
	err := db.Where("id = ?", filesFor).Last(&baseScript).Error
	if err != nil {
		cnt.JSON(http.StatusNotFound, gin.H{"error": err, "value": filesFor})
		return
	}

	var existingFiles []DBScriptFile
	db.Where("db_base_script_id = ?", baseScript.ID).Find(&existingFiles)
	for _, exFile := range existingFiles {
		fExt := filepath.Ext(exFile.FileName)
		err := os.Rename(exFile.PathToFile, "./trash/"+baseScript.ID+"_"+exFile.ID+fExt)
		if err != nil {
			cnt.JSON(http.StatusInternalServerError, gin.H{"error": err, "value": "move to trash"})
			return
		}
		exFile.PathToFile = "./trash/" + baseScript.ID + "_" + exFile.ID + fExt
		db.Save(&exFile)
		db.Delete(&exFile)
	}

	form, _ := cnt.MultipartForm()
	files := form.File["upload[]"]
	mainFile := form.Value
	//fmt.Println(mainFile)

	for _, file := range files {
		Log.Trace(file.Filename)
		err := cnt.SaveUploadedFile(file, "./scripts/"+file.Filename)
		if err != nil {
			Log.Error(err)
		}

		var mainFlag bool
		if file.Filename == mainFile["main_file"][0] {
			mainFlag = true
		}

		var newFile DBScriptFile
		newFile = DBScriptFile{
			ScriptFile:     mainFlag,
			DBBaseScriptID: baseScript.ID,
			PathToFile:     "./scripts/" + file.Filename,
			FileName:       file.Filename,
		}
		db.Create(&newFile)

	}

	var scriptResponse DBBaseScript
	db.Preload(clause.Associations).
		Where("id = ?", baseScript.ID).
		Last(&scriptResponse)

	cnt.JSON(http.StatusOK, scriptResponse)

}

func GetBaseScripts(cnt *gin.Context) {

	var baseScripts []DBBaseScript

	page, perPage := setPaginationParameters(cnt)

	format := cnt.Query("format")

	DB := db

	DB.Preload(clause.Associations).
		Order("created_at desc").
		Find(&baseScripts).
		Limit(perPage).
		Offset(page - 1*perPage)

	switch format {

	default:
		cnt.JSON(http.StatusOK, baseScripts)
	}

}

func UploadFiles(cnt *gin.Context) {
	form, _ := cnt.MultipartForm()
	files := form.File["upload[]"]
	for _, file := range files {
		Log.Trace(file.Filename)
		err := cnt.SaveUploadedFile(file, "./scripts/"+file.Filename)
		if err != nil {
			Log.Error(err)
		}
	}
}
