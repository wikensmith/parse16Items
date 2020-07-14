package utils

import (
	"fmt"
	"testing"
)

func TestQueryAirLine(t *testing.T) {
	// YVR 加拿大
	result, err := QueryAirport("CKG")
	fmt.Println(result, err)



}