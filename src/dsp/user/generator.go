package user

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	MAX_USER_PROFILE_COUNTER   = 26
	MAX_ATTRIBUTE_VALUE        = 200
	RANDOM_GENERATOR_CHAN_SIZE = 2
	USER_GENERATOR_CHAN_SIZE   = 2
)

var userGlobalCounter uint64
var userProfileGlobalCounter uint64

var rndChan chan int

// Initialize user generator
func InitGenerator() chan *User {
	initRandomGenerator()

	userChan := make(chan *User, USER_GENERATOR_CHAN_SIZE)
	for i := 0; i < USER_GENERATOR_CHAN_SIZE; i++ {
		go func() {
			for {
				userChan <- generate()
			}
		}()
	}
	return userChan
}

// Creates goroutine with random generator
func initRandomGenerator() {
	rndChan = make(chan int, RANDOM_GENERATOR_CHAN_SIZE*USER_GENERATOR_CHAN_SIZE)
	for i := 0; i < RANDOM_GENERATOR_CHAN_SIZE*USER_GENERATOR_CHAN_SIZE; i++ {
		go func() {
			source := rand.NewSource(time.Now().UnixNano())
			rnd := rand.New(source)
			for {
				rndChan <- rnd.Intn(MAX_ATTRIBUTE_VALUE-1) + 1
			}
		}()
	}
}

// Generates user with filled attributes
func generate() *User {
	return &User{
		ID:      "u" + strconv.Itoa(int(atomic.AddUint64(&userGlobalCounter, 1))),
		Profile: generateProfile(),
	}
}

// Generates user profile with userProfileCounter attributes
func generateProfile() map[string]string {
	userProfileCounter := int(atomic.AddUint64(&userProfileGlobalCounter, 1))
	if userProfileGlobalCounter > MAX_USER_PROFILE_COUNTER {
		userProfileCounter = 1
		atomic.StoreUint64(&userProfileGlobalCounter, 1)
	}

	profile := make(map[string]string, userProfileCounter)
	for i := 0; i < userProfileCounter; i++ {
		character := string(i + 'A')
		profile["attr_"+character] = character + strconv.Itoa(<-rndChan)
	}
	return profile
}
