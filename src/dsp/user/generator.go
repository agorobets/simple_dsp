package user

import (
	"math/rand"
	"strconv"
	"sync"
)

const (
	MAX_USER_PROFILE_COUNTER = 26
	MAX_ATTRIBUTE_VALUE      = 200
)

var userGlobalCounter = 0
var userProfileGlobalCounter = 0
var counterMutex sync.Mutex

// Generates user with filled attributes
func Generate() *User {
	userCounter, userProfileCounter := incrementCounters()
	return &User{
		ID:      "u" + strconv.Itoa(userCounter),
		Profile: generateProfile(userProfileCounter),
	}
}

// Generates user profile with userProfileCounter attributes
func generateProfile(userProfileCounter int) map[string]string {
	profile := make(map[string]string, userProfileCounter)
	for i := 0; i < userProfileCounter; i++ {
		character := string(i + 'A')
		profile["attr_"+character] = character + strconv.Itoa(rand.Intn(MAX_ATTRIBUTE_VALUE))
	}
	return profile
}

// Increments user counters on every Generate() call
func incrementCounters() (int, int) {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	userGlobalCounter++
	userProfileGlobalCounter++

	if userProfileGlobalCounter > MAX_USER_PROFILE_COUNTER {
		userProfileGlobalCounter = 1
	}

	return userGlobalCounter, userProfileGlobalCounter
}
