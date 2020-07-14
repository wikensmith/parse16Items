package modules

import (
	"fmt"
	"testing"
)

func TestConvertTimeToBeijing(t *testing.T) {
	// YVR CKG MEL
	result, err := ConvertTimeToBeijing("YVR", "2020-07-13 00:00:00")
	fmt.Println(result, err)

}