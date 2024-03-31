package utils

import "math/rand"

func RandomFloat() float64 {
	randomNum := rand.Float64()*90000 + 10000
	randomNum = float64(int(randomNum*100)) / 100
	return randomNum
}
