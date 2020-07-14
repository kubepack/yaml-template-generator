package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	ky "sigs.k8s.io/yaml"
	"strings"
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
d: *fa
e: z`

var yt_2 = `x:
  - a
  - b`

var yt__2 = `a`

type SA struct {
	A string `json:"a"`
	B string `json:"b"`
}

func main() {
	data, err := ioutil.ReadFile("/home/tamal/go/src/github.com/appscodelabs/tasty-kube/busy-dep.yaml")
	if err != nil {
		panic(err)
	}

	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		panic(err)
	}

	indent := 0
	column := 0
	var buf bytes.Buffer
	err = templatize(&node, &buf, indent, column, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

func templatize(node *yaml.Node, buf *bytes.Buffer, shift, column int, path []string) error {
	switch node.Kind {
	case yaml.DocumentNode:
		return templatize(node.Content[0], buf, shift, node.Column, path)
	case yaml.SequenceNode:
		// end it here
		buf.WriteString(fmt.Sprintf(" {{ toJson .Values.%s }}\n", strings.Join(path, ".")))
		return nil
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i = i + 2 {
			buf.WriteString(fmt.Sprintf("%s%s:", strings.Repeat(" ", shift+node.Content[i].Column-1), node.Content[i].Value))

			nextMap := node.Content[i+1].Kind == yaml.MappingNode ||
				(node.Content[i+1].Kind == yaml.AliasNode && node.Content[i+1].Alias.Kind == yaml.MappingNode)
			if nextMap {
				buf.WriteString("\n")
			}

			err := templatize(node.Content[i+1], buf, shift, node.Content[i].Column, append(path, node.Content[i].Value))
			if err != nil {
				return err
			}
		}
		return nil
	case yaml.ScalarNode:
		// end it here
		buf.WriteString(fmt.Sprintf(" {{ .Values.%s }}\n", strings.Join(path, ".")))
		return nil
	case yaml.AliasNode:
		return templatize(node.Alias, buf, shift, node.Alias.Column, path)
	}
	return nil
}

func traverse(node *yaml.Node, path []string) error {
	switch node.Kind {
	case yaml.DocumentNode:
		return traverse(node.Content[0], path)
	case yaml.SequenceNode:
		// end it here
		fmt.Println("> ", strings.Join(path, "."))
		return nil
	case yaml.MappingNode:
		// even number nodes
		for i := 0; i < len(node.Content); i = i + 2 {
			err := traverse(node.Content[i+1], append(path, node.Content[i].Value))
			if err != nil {
				return err
			}
		}
		return nil
	case yaml.ScalarNode:
		fmt.Println("> ", strings.Join(path, "."))
		return nil
	case yaml.AliasNode:
		return traverse(node.Alias, path)
	}
	return nil
}

func main22() {
	var mj yaml.Node
	err := yaml.Unmarshal([]byte(yt2), &mj)
	if err != nil {
		panic(err)
	}
	mj.Content[0].Content[1].Content[1].Value = "{{ .Values.xyz }}"

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
