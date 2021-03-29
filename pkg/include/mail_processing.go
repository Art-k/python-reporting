package include

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
)

func GetOutgoingEmails(cnt *gin.Context) {
	page, perPage := setPaginationParameters(cnt)
	format := cnt.Query("format")
	dbJobId := cnt.Query("db_job_id")

	var emails []DBOutgoingMails

	DB := db

	if dbJobId != "" {
		DB = DB.Model(&DBOutgoingMails{}).Where("db_job_id = ?", dbJobId)
	}

	var resp getResponse
	DB.Model(&DBOutgoingMails{}).Count(&resp.Total)
	DB.Preload(clause.Associations).
		Order("created_at desc").
		Limit(perPage).
		Offset(page - 1*perPage).
		Find(&emails)
	resp.Entities = emails
	resp.Current = len(emails)

	switch format {

	default:
		cnt.JSON(http.StatusOK, resp)
	}

}

func GetEmailHistory(cnt *gin.Context) {
	page, perPage := setPaginationParameters(cnt)
	format := cnt.Query("format")
	dbOutgoingMailsId := cnt.Query("db_outgoing_mails_id")

	var emailsHistory []DBOutgoingMailHistory

	DB := db

	if dbOutgoingMailsId != "" {
		DB = DB.Model(&DBOutgoingMailHistory{}).Where("db_outgoing_mails_id = ?", dbOutgoingMailsId)
	}

	var resp getResponse
	DB.Model(&DBOutgoingMailHistory{}).Count(&resp.Total)
	DB.Preload(clause.Associations).
		Order("created_at desc").
		Limit(perPage).
		Offset(page - 1*perPage).
		Find(&emailsHistory)
	resp.Entities = emailsHistory
	resp.Current = len(emailsHistory)

	switch format {

	default:
		cnt.JSON(http.StatusOK, resp)
	}
}
