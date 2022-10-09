package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

// GetMsg Bot Api 接受来自cqhttp上报的信息
func GetMsg(c *gin.Context) {
	var jsonReq string
	if c.Request.Body != nil {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		jsonReq = string(data)
	}
	log.Println(jsonReq)
}
