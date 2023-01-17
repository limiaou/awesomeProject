package send

import (
	"fmt"
	"log"
	"testing"
)

func TestReadExcelMsg(t *testing.T) {
	msg, err := ReadExcelMsg()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range msg {
		fmt.Println(v)
	}
	fmt.Println(msg)
}
