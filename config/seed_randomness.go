package config

import (
	"math/rand"
	"time"
)

func init() {
	// By default, the math/rand package will have the same seed when starting up
	// That means unless something happens randomly, anything using the fast
	// random generation of math/rand will always generate the same values.
	// With this code, that won't happen - the seed will be based on the current time.
	rand.Seed(time.Now().UnixNano())
}
