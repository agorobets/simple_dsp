package bidding

import (
	"dsp/campaign"
	"dsp/user"
	"math/rand"
	"sync"
)

type Winner struct {
	CampaignName string `json:"winner"`
	Counter      int    `json:"counter"`
}

var data = newData()

var counter struct {
	sync.Mutex
	n int
}

// Creates new data struct and fill it with campaigns
// when new data created, current data will be replaced with new
func ReloadData(campaigns []campaign.Campaign) error {
	newData := newData()

	for _, c := range campaigns {
		if err := newData.AddCampaign(&c); err != nil {
			return err
		}
	}

	data.Reload(newData)
	return nil
}

func ProcessBid(u *user.User) (*Winner, error) {

	var targetLeaf string

	winner := &Winner{
		CampaignName: "none",
		Counter:      incrementCounter(),
	}

	validCampaigns := data.getCampaignsByUserProfile(u.Profile)
	if len(validCampaigns) > 0 {
		targetedCampaigns := data.targetCampaigns(validCampaigns, u.Profile)
		if len(targetedCampaigns) > 0 {
			winner.CampaignName = runAuction(targetedCampaigns)
		}
	}

	return winner, nil
}

func runAuction(campaigns map[string]float64) string {
	var maxBid float64

	maxBidders := []string{}
	for name, price := range campaigns {
		if maxBid < price {
			maxBidders = []string{name}
			maxBid = price
		} else if maxBid == price {
			maxBidders = append(maxBidders, price)
		}
	}

	n := len(maxBidders)
	if n == 1 {
		return maxBidders[0]
	}

	return maxBidders[rand.Intn(n)]
}

func incrementCounter() int {
	c.Lock()
	defer c.Unlock()

	countrer.n++
	return countrer.n
}
