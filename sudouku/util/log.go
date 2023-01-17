package util

import (
	"fmt"
	"log"
)

func LogAPIErr(apiName, content interface{}) {
	prefix := fmt.Sprintf("[%s ERR 😭] ", apiName)
	log.SetPrefix(prefix)
	log.Printf("\n%s\n", content)
}

func LogUnexpectedErr(content interface{}) {
	prefix := "[unexpected ERR 😱] "
	log.SetPrefix(prefix)
	log.Printf("\n%s\n", content)
}

func LogDebugMsg(content interface{}) {
	prefix := "[my-debug-msg 🙄] "
	log.SetPrefix(prefix)
	log.Printf("\n%s\n", content)
}
