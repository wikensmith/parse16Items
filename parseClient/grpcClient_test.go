package parseClient

import (
	"fmt"
	"strconv"
	"testing"
)



func Test_main2(t *testing.T) {
	//l := []string{
	//}
	res, err := main2("", "")
	fmt.Println(err, res)
}

func TestDoLstTest(t *testing.T) {
	DoLstTest("GS")
}

func TestAA(t *testing.T)  {
	a := 2080.00123
	fmt.Println(strconv.FormatFloat(a,'f', 2, 64))
}
