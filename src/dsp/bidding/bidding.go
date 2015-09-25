package bidding

import (
	"dsp/campaign"
	"dsp/user"
	"fmt"
	"math/rand"
	"sync/atomic"
)

type Winner struct {
	CampaignName string `json:"winner"`
	Counter      uint64 `json:"counter"`
}

var data = newData()

var counter uint64

// Creates new data struct and fill it with campaigns
// when new data created, current data will be replaced with new
func ReloadData(campaigns []campaign.Campaign) error {
	newData := newData()

	for _, c := range campaigns {
		if err := validateCampaign(&c); err != nil {
			return err
		}
		if err := newData.AddCampaign(c); err != nil {
			return err
		}
	}

	data.Reload(newData)
	return nil
}

// Searches targeted campaigns, and determines winner by user
func ProcessBid(u *user.User) Winner {
	winner := Winner{
		CampaignName: "none",
		Counter:      atomic.AddUint64(&counter, 1),
	}

	validCampaigns := data.GetCampaignsByUserProfile(u.Profile)
	if len(validCampaigns) > 0 {
		targetedCampaigns := data.GetTargetedCampaigns(validCampaigns, u.Profile)
		if len(targetedCampaigns) > 0 {
			winner.CampaignName = runAuction(targetedCampaigns)
		}
	}

	return winner
}

// Returns data object
func GetData() *dataStruct {
	data.RLock()
	defer data.RUnlock()

	return data
}

// Determines the winner of auction by max price and returns its name
func runAuction(campaigns map[string]float64) string {
	var maxBid float64
	var maxBidders []string

	for name, price := range campaigns {
		if maxBid < price {
			maxBidders = []string{name}
			maxBid = price
		} else if maxBid == price {
			maxBidders = append(maxBidders, name)
		}
	}

	n := len(maxBidders)
	if n == 1 {
		return maxBidders[0]
	}

	// If determined more than one winner with max price
	// randomly choose one
	return maxBidders[rand.Intn(n)]
}

// Validates campaign attributes
func validateCampaign(c *campaign.Campaign) error {
	if c.Name == "" {
		return fmt.Errorf("Campaign name is empty.")
	}

	if len(c.TargetList) == 0 {
		return fmt.Errorf("Campaign '%s' hasn't any target.", c.Name)
	}

	for _, target := range c.TargetList {
		if target.Name == "" {
			return fmt.Errorf("Target name is empty in campaing '%s'.", c.Name)
		}
		if len(target.Values) == 0 {
			return fmt.Errorf("Target attribute list is empty in campaign '%s'.", c.Name)
		}
	}

	return nil
}
