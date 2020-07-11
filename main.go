package main

type SA struct {
	A string
	B string
}
func main() {
	m := map[SA]string{}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
