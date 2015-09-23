package bidding

import (
	"dsp/campaign"
)

type Winner struct {
}

var data = newData()

// Creates new data struct and fill it with campaigns
// when new data created, current data will be replaced with new
func ReloadData(campaigns []campaign.Campaign) error {
	newData := newData()

	for _, c := range campaigns {
		if err := newData.addCampaign(&c); err != nil {
			return err
		}
		if err := newData.addTargetLeaves(&c); err != nil {
			return err
		}
		if err := newData.addAttributeLeaves(&c); err != nil {
			return err
		}
	}

	data.reload(newData)
	return nil
}
