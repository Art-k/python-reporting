package include

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func createHash() string {
	return uuid.New().String()
}

func GETLogoProcessing(cnt *gin.Context) {

	id := cnt.Param("id")

	var msg DBOutgoingMails
	db.Where("id = ?", id).Last(&msg)
	if msg.ID != "" {

		var msgHistory DBOutgoingMailHistory
		msgHistory.DBOutgoingMailsID = msg.ID
		msgHistory.RecType = "opened"
		msgHistory.HistoryMessage = cnt.Request.RemoteAddr
		db.Create(&msgHistory)

	}

	http.ServeFile(cnt.Writer, cnt.Request, "./assets/icon.png")

}

func setPaginationParameters(cnt *gin.Context) (page, perPage int) {

	var err error

	pageStr := cnt.Query("page")
	perPageStr := cnt.Query("per-page")

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	} else {
		page = 1
	}

	if perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			perPage = 15
		}
	} else {
		perPage = 15
	}

	return page, perPage

}
