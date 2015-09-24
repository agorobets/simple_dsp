package main

import (
	"dsp/bidding"
	"dsp/campaign"
	"dsp/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strconv"
)

const (
	MAX_NUMBER_OF_CAMPAIGNS  = 10000
	MAX_NUMBER_OF_TARGETS    = 26
	MAX_NUMBER_OF_ATTRIBUTES = 100
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	router := gin.Default()

	router.GET("/campaign", GetCampaigns)
	router.GET("/user", GetUser)
	router.POST("/import_camp", ImportCampaigns)
	router.POST("/search", Search)
	router.GET("/search_auto", SearchAuto)
	router.GET("/dump", DumpData)

	router.Run(":3000")
}

// Handler generates JSON array of campaigns
// Has 3 not required arguments:
//   x - max number of attributes per target in campaign target list (max: 100)
//   y - max number of targets in campaign target list (max: 26)
//   z - number of campaigns to generate (max: 10000)
// Responses:
//   200 OK - return JSON array of campaigns
//   400 Bad Request
func GetCampaigns(ctx *gin.Context) {
	maxAttributesNumber, err := strconv.Atoi(ctx.DefaultQuery("x", strconv.Itoa(MAX_NUMBER_OF_ATTRIBUTES)))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Argument 'x' must be a number")
		return
	}

	maxTargetsNumber, err := strconv.Atoi(ctx.DefaultQuery("y", strconv.Itoa(MAX_NUMBER_OF_TARGETS)))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Argument 'y' must be a number")
		return
	}

	campaignsNumber, err := strconv.Atoi(ctx.DefaultQuery("z", strconv.Itoa(MAX_NUMBER_OF_CAMPAIGNS)))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Argument 'z' must be a number")
		return
	}

	if err := validateCampaignArguments(campaignsNumber, maxTargetsNumber, maxAttributesNumber); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(200, campaign.GenerateMany(campaignsNumber, maxTargetsNumber, maxAttributesNumber))
}

// Generates randomly filled User JSON object
// Responses:
//   200 OK - return JSON object of user
func GetUser(ctx *gin.Context) {
	ctx.JSON(200, user.Generate())
}

// Imports campaigns and updates bidding.data struct
// Json array of campaigns should be send in request body (POST)
// Responses:
//    200 OK - return nothing
//    400 Bad Request - return error string
func ImportCampaigns(ctx *gin.Context) {
	campaigns := make([]campaign.Campaign, 0)

	if ctx.BindJSON(&campaigns) != nil {
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	if err := bidding.ReloadData(campaigns); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.String(200, "")
}

// Returns won campaign by user profile and campaign max price
// Json array of campaigns should be send in request body (POST)
// Responses:
//    200 OK - return winner struct as Json object
//    400 Bad Request - return error string
func Search(ctx *gin.Context) {
	u := &user.User{}

	if ctx.BindJSON(u) != nil {
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(200, bidding.ProcessBid(u))
}

// Generates user and tries to find campaign by user profile and max price
// Responses:
//    200 OK - return winner struct as Json object
//    400 Bad Request - return error string
func SearchAuto(ctx *gin.Context) {
	ctx.JSON(200, bidding.ProcessBid(user.Generate()))
}

// Returns dump of bidding.data object as JSON
// Responses:
//    200 OK
func DumpData(ctx *gin.Context) {
	ctx.JSON(200, bidding.GetData())
}

// Validates arguments of GetCampaigns handler
func validateCampaignArguments(campaignsNumber, maxTargetsNumber, maxAttributesNumber int) error {
	if maxAttributesNumber > MAX_NUMBER_OF_ATTRIBUTES {
		return fmt.Errorf("Argument 'x' cannot be more than %d", MAX_NUMBER_OF_ATTRIBUTES)
	}

	if maxTargetsNumber > MAX_NUMBER_OF_TARGETS {
		return fmt.Errorf("Argument 'y' cannot be more than %d", MAX_NUMBER_OF_TARGETS)
	}

	if campaignsNumber > MAX_NUMBER_OF_CAMPAIGNS {
		return fmt.Errorf("Argument 'z' cannot be more than %d", MAX_NUMBER_OF_CAMPAIGNS)
	}

	return nil
}
