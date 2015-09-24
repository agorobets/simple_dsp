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

// Creates new dataStruct with initialized fields
func newData() *dataStruct {
	return &dataStruct{
		campaigns:       make(map[string]*campaign.Campaign),
		targetLeaves:    make(map[string][]string),
		attributeLeaves: make(map[string]map[string]bool),
	}
}

// Adds campaign
// If campaign already added or not valid, will return error
func (d *dataStruct) addCampaign(c *campaign.Campaign) error {
	d.Lock()
	defer d.Unlock()

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
	d.Lock()
	defer d.Unlock()

	if _, ok := d.campaigns[c.Name]; !ok {
		return fmt.Errorf("Campaign '%s' doesn't exist", c.Name)
	}

	targetLeaves := make([]string, len(c.TargetList))
	createTargetLeaves(c, targetLeaves)

	for _, leave := range targetLeaves {

	}

	return nil
}

func (d *dataStruct) addAttributeLeaves(c *campaign.Campaign) error {
	d.Lock()
	defer d.Unlock()

	if _, ok := d.campaigns[c.Name]; !ok {
		return fmt.Errorf("Campaign '%s' doesn't exist", c.Name)
	}

	return nil
}

// Reloads data attributes with newData struct
func (d *dataStruct) reload(newData *dataStruct) {
	d.Lock()
	defer d.Unlock()

	d.campaigns = newData.campaigns
	d.targetLeaves = newData.targetLeaves
	d.attributeLeaves = newData.attributeLeaves
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

func createTargetLeaves(c *campaign.Campaign, leaves []string) {
	leaves := make([]string, len(c.TargetList))

	return leaves
}
