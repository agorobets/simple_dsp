package bidding

import (
	"dsp/campaign"
	"fmt"
	"sync"
)

type Campaign struct {
	Name       string                     `json:"campaign_name"`
	Price      float64                    `json:"price"`
	TargetList map[string]map[string]bool `json:"target_list"`
}

type dataStruct struct {
	sync.RWMutex
	Campaigns    map[string]*Campaign   `json:"campaigns"`
	TargetLeaves map[string][]*Campaign `json:"leaves"`
}

// Creates new dataStruct with initialized fields
func newData() *dataStruct {
	return &dataStruct{
		Campaigns:    make(map[string]*Campaign),
		TargetLeaves: make(map[string][]*Campaign),
	}
}

// Adds campaign
// If campaign already added or not valid, will return error
func (d *dataStruct) AddCampaign(c campaign.Campaign) error {
	if _, ok := d.Campaigns[c.Name]; ok {
		return fmt.Errorf("Found more than one campaign '%s'", c.Name)
	}

	d.Campaigns[c.Name] = &Campaign{
		Name:       c.Name,
		Price:      c.Price,
		TargetList: buildTargetingList(c.TargetList),
	}

	if err := d.addTargetLeaf(d.Campaigns[c.Name]); err != nil {
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
func (d *dataStruct) GetTargetedCampaigns(profile map[string]string) []*Campaign {
	d.RLock()
	defer d.RUnlock()

	var targetLeaf string
	var matched bool

	targetedCampaigns := []*Campaign{}
	for target, _ := range profile {
		targetLeaf = targetLeaf + target
		if campaigns, ok := d.TargetLeaves[targetLeaf]; ok {
			for _, c := range campaigns {
				matched = true
				for name, vals := range c.TargetList {
					if !vals[profile[name]] {
						matched = false
						break
					}
				}

				if matched {
					targetedCampaigns = append(targetedCampaigns, c)
				}
			}
		}
	}
	return targetedCampaigns
}

// Adds leaf to targetLeaves map, where keys is concatenated string of campaign target names
// and value is array of campaign ids
func (d *dataStruct) addTargetLeaf(c *Campaign) error {
	var leaf string
	for name, _ := range c.TargetList {
		leaf = leaf + name
	}

	if _, ok := d.TargetLeaves[leaf]; !ok {
		d.TargetLeaves[leaf] = []*Campaign{c}
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

func buildTargetingList(tl []campaign.Target) map[string]map[string]bool {
	list := make(map[string]map[string]bool, len(tl))
	for _, t := range tl {
		list[t.Name] = buildAttributesHash(t.Values)
	}
	return list
}

func buildAttributesHash(vals []string) map[string]bool {
	attrs := make(map[string]bool, len(vals))
	for _, val := range vals {
		attrs[val] = true
	}
	return attrs
}
