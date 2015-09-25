package user

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	MAX_USER_PROFILE_COUNTER = 26
	MAX_ATTRIBUTE_VALUE      = 200
)

var userGlobalCounter uint64
var userProfileGlobalCounter uint64

var rnd *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(source)
}

// Generates user with filled attributes
func Generate() *User {
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
		profile["attr_"+character] = character + strconv.Itoa(int(rnd.Float32()*MAX_ATTRIBUTE_VALUE)+1)
	}
	return profile
}
