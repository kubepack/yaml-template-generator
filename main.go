package main

import (
	"encoding/json"
	"fmt"
)

type SA struct {
	A string `json:"a"`
	B string `json:"b"`
}

func main() {
	m := map[string]SA{}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
