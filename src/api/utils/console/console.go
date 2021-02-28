package console

import (
	"encoding/json"
	"fmt"
	"log"
)

//Pretty prints data in a pretty way
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}
