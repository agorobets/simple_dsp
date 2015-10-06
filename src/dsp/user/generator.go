package user

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	MAX_USER_PROFILE_COUNTER = 26
	MAX_ATTRIBUTE_VALUE      = 200
)

var counter uint64
var rndChan chan int

// Creates goroutine with random generator
func InitRandomGenerator() {
	rndChan = make(chan int, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
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
func Generate() *User {
	id := int(atomic.AddUint64(&counter, 1))
	return &User{
		ID:      "u" + strconv.Itoa(id),
		Profile: generateProfile(id % MAX_USER_PROFILE_COUNTER),
	}
}

// Generates user profile
func generateProfile(length int) map[string]string {
	if length == 0 {
		length = MAX_USER_PROFILE_COUNTER
	}

	profile := make(map[string]string, length)
	for i := 0; i < length; i++ {
		character := string(i + 'A')
		profile["attr_"+character] = character + strconv.Itoa(<-rndChan)
	}
	return profile
}
