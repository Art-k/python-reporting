package include

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func createHash() string {
	return uuid.New().String()
}

func setPaginationParameters(cnt *gin.Context) (page, perPage int) {

	var err error

	pageStr := cnt.Query("page")
	perPageStr := cnt.Query("per-page")

	if pageStr == "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}

	if perPageStr == "" {
		perPage, err = strconv.Atoi(pageStr)
		if err != nil {
			perPage = 15
		}
	}

	return page, perPage

}
