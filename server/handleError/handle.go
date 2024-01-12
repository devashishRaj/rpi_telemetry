package handleError

import (
	"fmt"
	"log"
)

func CheckError(context string, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", context)
		log.Fatalln(err)

	}
}
