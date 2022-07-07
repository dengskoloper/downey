package util

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func ReadHeadersFromJSON(r *http.Header) {
	jsonBytesFromFile, err := ioutil.ReadFile("headers.json")
	if err != nil {
		panic(err)
	}
	var jsonHeaderMap map[string]string

	err = json.Unmarshal(jsonBytesFromFile, &jsonHeaderMap)
	if err != nil {
		panic(err)
	}

	for k, v := range jsonHeaderMap {
		fmt.Printf("Setting HTTP Header \"%s : %s\"\n", k, v)
		r.Set(k, v)
	}
}