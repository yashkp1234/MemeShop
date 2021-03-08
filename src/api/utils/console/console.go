package console

import (
	"encoding/json"
	"log"
)

//Pretty prints data in a pretty way
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}
