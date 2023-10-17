package handlerror

import (
	"fmt"
	"log"
)

func IsNil(context string, v interface{}) {
	if v == nil {
		fmt.Printf("Error: %s\n", context)
		log.Fatal("null value ")
	}
}

func CheckError(context string, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", context)
		log.Fatalln(err)

	}
}

func ReutrnMinus(context string, err error) {

	fmt.Printf("Error: %s\n", context)
}
