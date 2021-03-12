package include

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetBaseScripts(cnt *gin.Context) {

	var baseScripts []DBBaseScript

	var page int
	var perPage int
	var err error

	pageStr := cnt.Query("page")
	perPageStr := cnt.Query("per-page")
	enabled := cnt.Query("enabled")
	format := cnt.Query("format")

	fmt.Println(pageStr)
	fmt.Println(perPageStr)
	fmt.Println(enabled)

	DB := db
	if enabled == "" {
		DB = db.Where("enabled <> ?", true)
	}

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

	DB.Find(&baseScripts).Limit(perPage).Offset(page - 1*perPage)

	switch format {

	default:
		cnt.JSON(http.StatusOK, baseScripts)
	}

}
