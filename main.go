package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	ky "sigs.k8s.io/yaml"
)

var yaml_text = `services:
  # Node.js gives OS info about the node (Host)
  nodeinfo: &function
    image: functions/nodeinfo:latest
    labels:
      function: "true"
    depends_on:
      - gateway
    networks:
      - functions
    environment:
      no_proxy: "gateway"
      https_proxy: $https_proxy
    deploy:
      placement:
        constraints:
          - 'node.platform.os == linux'
  # Uses cat to echo back response, fastest function to execute.
  echoit:
    <<: *function
    image: functions/alpine:health
    environment:
      fprocess: "cat"
      no_proxy: "gateway"
      https_proxy: $https_proxy`

var yt2 = `a: &fa
  b: &fb x
  c: "y"
d: *fb`

type SA struct {
	A string `json:"a"`
	B string `json:"b"`
}

func main() {
	var mj yaml.Node
	err := yaml.Unmarshal([]byte(yt2), &mj)
	if err != nil {
		panic(err)
	}
	mj.Content[0].Content[1].Content[1].Value = "we"

	str, err := yaml.Marshal(&mj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))

	str, err = ky.YAMLToJSON([]byte(str))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))

	str, err = ky.YAMLToJSON([]byte(yaml_text))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))


	//m := map[string]SA{}
	//data, err := json.MarshalIndent(m, "", "  ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(data))
}
