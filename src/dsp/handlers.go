package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	MAX_NUMBER_OF_CAMPAIGNS  = 10000
	MAX_NUMBER_OF_TARGETS    = 26
	MAX_NUMBER_OF_ATTRIBUTES = 100
)

func GetCampaigns(ctx *gin.Context) {

	maxAttributesNumber, err := strconv.Atoi(ctx.DefaultQuery("x", string(MAX_NUMBER_OF_ATTRIBUTES)))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, gin.H{"error": "Argument 'x' must be a number"})
	}

	maxTargetsNumber, err := c.DefaultQuery("y", string(MAX_NUMBER_OF_TARGETS))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, gin.H{"error": "Argument 'y' must be a number"})
	}

	campaignsNumber, err := c.DefaultQuery("z", string(MAX_NUMBER_OF_CAMPAIGNS))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, gin.H{"error": "Argument 'z' must be a number"})
	}

	if err := validateCampaignArguments(campaignsNumber, maxTargetsNumber, maxAttributesNumber); err != nil {
		return ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, campaign.GenerateMany(maxAttributesNumber, maxTargetsNumber, maxAttributesNumber))
}

func ImportCampaigns(ctx *gin.Context) {
	ctx.String(200, "import_camp")
}

func GetUser(ctx *gin.Context) {
	ctx.JSON(200, user.Generate())
}

func Search(ctx *gin.Context) {
	ctx.String(200, "search")
}

func SearchAuto(ctx *gin.Context) {
	ctx.String(200, "search_auto")
}

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
