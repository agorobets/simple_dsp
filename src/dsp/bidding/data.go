package bidding

import (
	"dsp/campaign"
	"fmt"
	"sync"
)

type dataStruct struct {
	sync.RWMutex
	campaigns       map[string]*campaign.Campaign
	targetLeaves    map[string][]string
	attributeLeaves map[string]map[string]bool
}

func newData() *dataStruct {
	return &dataStruct{
		campaigns:       make(map[string]*campaign.Campaign),
		targetLeaves:    make(map[string][]string),
		attributeLeaves: make(map[string]map[string]bool),
	}
}

func (d *dataStruct) addCampaign(c *campaign.Campaign) error {
	if err := validateCampaign(c); err != nil {
		return err
	}

	if _, ok := d.campaigns[c.Name]; ok {
		return fmt.Errorf("Found more than one campaign '%s'", c.Name)
	}

	d.campaigns[c.Name] = c
	return nil
}

func (d *dataStruct) addTargetLeaves(c *campaign.Campaign) error {
	return nil
}

func (d *dataStruct) addAttributeLeaves(c *campaign.Campaign) error {
	return nil
}

func (d *dataStruct) reload(newData *dataStruct) {
	d.Lock()
	defer d.Unlock()

	d.campaigns = newData.campaigns
	d.targetLeaves = newData.targetLeaves
	d.attributeLeaves = newData.attributeLeaves
}

func validateCampaign(c *campaign.Campaign) error {
	return nil
}
