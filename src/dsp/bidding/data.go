package bidding

import (
	"dsp/campaign"
	"fmt"
	"sync"
)

type dataStruct struct {
	sync.RWMutex
	Campaigns    map[string]*campaign.Campaign `json:"campaigns"`
	TargetLeaves map[string][]string           `json:"leaves"`
}

// Creates new dataStruct with initialized fields
func newData() *dataStruct {
	return &dataStruct{
		Campaigns:    make(map[string]*campaign.Campaign),
		TargetLeaves: make(map[string][]string),
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
func (d *dataStruct) GetCampaignsByUserProfile(profile map[string]string) []string {
	d.RLock()
	defer d.RUnlock()

	var targetLeaf string

	validCampaigns := []string{}
	for target, _ := range profile {
		targetLeaf = targetLeaf + target
		if campaigns, ok := d.TargetLeaves[targetLeaf]; ok {
			validCampaigns = append(validCampaigns, campaigns...)
		}
	}
	return validCampaigns
}

// Returns campaigns, which all target attributes matched  with user profile attributes:
// {'campaign_name': price, ...}
func (d *dataStruct) GetTargetedCampaigns(campaignNames []string, profile map[string]string) map[string]float64 {
	d.RLock()
	defer d.RUnlock()

	targetedCampaigns := make(map[string]float64)
	for _, cName := range campaignNames {
		matchesCount := 0
		for _, target := range d.Campaigns[cName].TargetList {
			matched := false
			for _, attribute := range target.Values {
				if attribute == profile[target.Name] {
					matched = true
					break
				}
			}
			if !matched {
				break
			}
			matchesCount++
		}

		// if all campaign target attributes matched, then add it
		if matchesCount == len(d.Campaigns[cName].TargetList) {
			targetedCampaigns[cName] = d.Campaigns[cName].Price
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
		d.TargetLeaves[leaf] = []string{c.Name}
	} else {
		d.TargetLeaves[leaf] = append(d.TargetLeaves[leaf], c.Name)
	}

	return nil
}
