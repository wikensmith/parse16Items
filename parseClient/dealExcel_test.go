package parseClient

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	ticketLst, err := Read("./解析客票.xlsx")
	fmt.Println(err, ticketLst)

}

