package user

import (
	"math/rand"
	"strconv"
	"sync"
)

const (
	MAX_USER_PROFILE_COUNTER = 27
	MAX_ATTRIBUTE_VALUE      = 200
)

var userCounter = 0
var userProfileCounter = 0
var mutex sync.Mutex

// Generates user with filled attributes
func Generate() *User {
	mutex.Lock()

	incrementCounters()
	user := &User{
		ID:      "u" + strconv.Itoa(userCounter),
		Profile: generateProfile(),
	}

	mutex.Unlock()
	return user
}

// Generates user profile with userProfileCounter attributes
func generateProfile() map[string]string {
	profile := make(map[string]string, userProfileCounter)
	for i := 0; i < userProfileCounter; i++ {
		character := string(i + 'A')
		profile["attr_"+character] = character + strconv.Itoa(rand.Intn(MAX_ATTRIBUTE_VALUE))
	}
	return profile
}

// Increments user counters on every Generate() call
func incrementCounters() {
	userCounter++
	userProfileCounter++

	if userProfileCounter >= MAX_USER_PROFILE_COUNTER {
		userProfileCounter = 1
	}
}
