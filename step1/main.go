package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Data string `json:"data"`
}

func main() {
	ginEngine := gin.Default()
	ginEngine.GET("/myapi", myApiHandler)
	ginEngine.Run()
}

func myApiHandler(context *gin.Context) {
	data := context.Request.URL.Query().Get("data")
	apiResponse := ApiResponse{}
	apiResponse.Data = data
	context.JSON(http.StatusOK, apiResponse)
}
