package bidding

import (
	"dsp/campaign"
	"fmt"
	"sync"
)

type dataStruct struct {
	sync.RWMutex
	Campaigns    map[string]*campaign.Campaign   `json:"campaigns"`
	TargetLeaves map[string][]*campaign.Campaign `json:"leaves"`
}

// Creates new dataStruct with initialized fields
func newData() *dataStruct {
	return &dataStruct{
		Campaigns:    make(map[string]*campaign.Campaign),
		TargetLeaves: make(map[string][]*campaign.Campaign),
	}
}

// Adds campaign
// If campaign already added or not valid, will return error
func (d *dataStruct) AddCampaign(c campaign.Campaign) error {
	if _, ok := d.Campaigns[c.Name]; ok {
		return fmt.Errorf("Found more than one campaign '%s'", c.Name)
	}

	d.Campaigns[c.Name] = &c

	if err := d.addTargetLeaf(&c); err != nil {
		return err
	}

	return nil
}

// Reloads data attributes with newData struct
func (d *dataStruct) Reload(newData *dataStruct) {
	d.Lock()
	defer d.Unlock()

	d.Campaigns = newData.Campaigns
	d.TargetLeaves = newData.TargetLeaves
}

// Returns array of campaigns, which all targeting also have in user profile
func (d *dataStruct) GetTargetedCampaigns(profile map[string]string) []*campaign.Campaign {
	d.RLock()
	defer d.RUnlock()

	var targetLeaf string

	targetedCampaigns := []*campaign.Campaign{}
	for target, _ := range profile {
		targetLeaf = targetLeaf + target
		if campaigns, ok := d.TargetLeaves[targetLeaf]; ok {
			for _, campaign := range campaigns {
				if IsTargetedCampaign(campaign, profile) {
					targetedCampaigns = append(targetedCampaigns, campaign)
				}
			}
		}
	}
	return targetedCampaigns
}

// Adds leaf to targetLeaves map, where keys is concatenated string of campaign target names
// and value is array of campaign ids
func (d *dataStruct) addTargetLeaf(c *campaign.Campaign) error {
	var leaf string
	for _, t := range c.TargetList {
		leaf = leaf + t.Name
	}

	if _, ok := d.TargetLeaves[leaf]; !ok {
		d.TargetLeaves[leaf] = []*campaign.Campaign{c}
	} else {
		d.TargetLeaves[leaf] = append(d.TargetLeaves[leaf], c)
	}

	return nil
}

// Returns campaigns, which all target attributes matched  with user profile attributes:
// {'campaign_name': price, ...}
func IsTargetedCampaign(campaign *campaign.Campaign, profile map[string]string) bool {
	for _, target := range campaign.TargetList {
		matched := false
		for _, attribute := range target.Values {
			if attribute == profile[target.Name] {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
