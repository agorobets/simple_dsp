package campaign

import (
	"math/rand"
	"strconv"
)

const (
	MAX_PRICE_INTEGER_PART = 5
)

func GenerateMany(number, maxTargetsNumber, maxAttributesNumber int) []*Campaign {
	campaigns := make([]*Campaign, number)
	for i := 0; i < number; i++ {
		campaigns[i] = generate(i+1, maxTargetsNumber, maxAttributesNumber)
	}
	return campaigns
}

func generate(id, maxTargetsNumber, maxAttributesNumber int) *Campaign {
	return &Campaign{
		Name:       "campaign" + strconv.Itoa(id),
		Price:      rand.Float64() * MAX_PRICE_INTEGER_PART,
		TargetList: generateTargetList(maxTargetsNumber, maxAttributesNumber),
	}
}

func generateTargetList(maxTargetsNumber, maxAttributesNumber int) []Target {
	targetsNumber := rand.Intn(maxTargetsNumber)
	attributesNumber := rand.Intn(maxAttributesNumber)

	targets := make([]Target, targetsNumber)
	for i := 0; i <= targetsNumber; i++ {

	}

	return targets
}
