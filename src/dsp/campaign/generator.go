package campaign

import (
	"math"
	"math/rand"
	"strconv"
)

const (
	PRICE_MULTIPLIER     = 5
	PRICE_ROUND_ON       = .5
	PRICE_DECIMAL_PLACES = 2
)

// Generates list of campaigns
func GenerateMany(number, maxTargetsNumber, maxAttributesNumber int) []*Campaign {
	campaigns := make([]*Campaign, number)
	for i := 0; i < number; i++ {
		campaigns[i] = generate(i+1, maxTargetsNumber, maxAttributesNumber)
	}
	return campaigns
}

// Generates one filled campaign
func generate(id, maxTargetsNumber, maxAttributesNumber int) *Campaign {
	return &Campaign{
		Name:       "campaign" + strconv.Itoa(id),
		Price:      round(rand.Float64()*PRICE_MULTIPLIER, PRICE_ROUND_ON, PRICE_DECIMAL_PLACES),
		TargetList: generateTargetList(maxTargetsNumber, maxAttributesNumber),
	}
}

// Generates list of campaign targets
func generateTargetList(maxTargetsNumber, maxAttributesNumber int) []Target {
	targetsNumber := rand.Intn(maxTargetsNumber)
	targets := make([]Target, targetsNumber+1)

	for i := 0; i < targetsNumber; i++ {
		character := string(i + 'A')
		targets[i] = Target{
			Attribute: "attr_" + character,
			Values:    generateAttributes(character, maxAttributesNumber),
		}
	}

	return targets
}

// Generates list of available attributes for campaign target
func generateAttributes(prefix string, maxAttributesNumber int) []string {
	attributesNumber := rand.Intn(maxAttributesNumber)
	attributes := make([]string, attributesNumber)
	for i := 0; i < attributesNumber; i++ {
		attributes[i] = prefix + strconv.Itoa(i)
	}
	return attributes
}

// Rounds fraction of float values to places digits
func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val

	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow
	return
}
