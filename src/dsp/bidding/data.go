package bidding

import (
	"dsp/campaign"
	"fmt"
	"sync"
)

type dataStruct struct {
	sync.RWMutex
	campaigns    map[string]*campaign.Campaign
	targetLeaves map[string][]string
}

// Creates new dataStruct with initialized fields
func newData() *dataStruct {
	return &dataStruct{
		campaigns:    make(map[string]*campaign.Campaign),
		targetLeaves: make(map[string][]string),
	}
}

// Adds campaign
// If campaign already added or not valid, will return error
func (d *dataStruct) AddCampaign(c *campaign.Campaign) error {
	d.Lock()
	defer d.Unlock()

	if err := validateCampaign(c); err != nil {
		return err
	}

	if _, ok := d.campaigns[c.Name]; ok {
		return fmt.Errorf("Found more than one campaign '%s'", c.Name)
	}

	d.campaigns[c.Name] = c

	if err := d.addTargetLeaf(c); err != nil {
		return err
	}

	return nil
}

// Adds leaf to targetLeaves map, where keys is concatenated string of campaign target names
// and value is array of campaign ids
func (d *dataStruct) addTargetLeaf(c *campaign.Campaign) error {
	var leaf string
	for _, t := range c.TargetList {
		leaf = leaf + t.Name
	}

	if _, ok := d.targetLeaves[leaf]; !ok {
		d.targetLeaves[leaf] = []string{c.Name}
	} else {
		d.targetLeaves[leaf] = append(d.targetLeaves[leaf], c.Name)
	}

	return nil
}

// Returns array of campaigns, which all targeting also have in user profile
func (d *dataStruct) getCampaignsByUserProfile(profile map[string]string) []string {
	var targetLeaf string

	validCampaigns := make(map[string]string)
	for target, attribute := range profile {
		targetLeaf = targetLeaf + target

		if campaigns, ok := d.targetLeaves[targetLeaf]; ok {
			validCampaigns[targetLeaf] = append(validCampaigns[targetLeaf], d.targetLeaves[targetLeaf]...)
		}
	}
	return validCampaigns
}

// Reloads data attributes with newData struct
func (d *dataStruct) Reload(newData *dataStruct) {
	d.Lock()
	defer d.Unlock()

	d.campaigns = newData.campaigns
	d.targetLeaves = newData.targetLeaves
}

// Validates campaign attributes
func validateCampaign(c *campaign.Campaign) error {
	if c.Name == "" {
		return fmt.Errorf("Campaign name is empty.")
	}

	if len(c.TargetList) == 0 {
		return fmt.Errorf("Campaign '%s' hasn't any target.", c.Name)
	}

	if c.Price <= 0 {
		return fmt.Errorf("Price <= 0 for campaign '%s'.", c.Name)
	}

	for _, target := range c.TargetList {
		if target.Attribute == "" {
			return fmt.Errorf("Target name is empty in campaing '%s'.", c.Name)
		}
		if len(target.Values) == 0 {
			return fmt.Errorf("Target attribute list is empty in campaign '%s'.", c.Name)
		}
	}

	return nil
}
