package main

import "github.com/gin-gonic/gin"
import "dsp/user"
import "dsp/campaign"

func main() {

	r := gin.Default()

	r.GET("/campaign", GetCampaigns)
	r.GET("/user", GetUser)
	r.POST("/import_camp", ImportCampaigns)
	r.POST("/search", Search)
	r.GET("/search_auto", SearchAuto)

	r.Run(":3000")
}
