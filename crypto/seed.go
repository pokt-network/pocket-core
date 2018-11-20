package crypto

import (
	"math/rand"
	"time"
)

func GenerateSeed(){
	rand.Seed(time.Now().UTC().UnixNano())
}
