package campaign

type Target struct {
	Name   string   `json:"target"`
	Values []string `json:"attr_list"`
}

type Campaign struct {
	Name       string   `json:"campaign_name"`
	Price      float64  `json:"price"`
	TargetList []Target `json:"target_list"`
}
