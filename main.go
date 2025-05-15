package main

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Input  string `json:"input"`
	Output string `json:"reply"`
}

func main() {
	fmt.Println("My Crypto Wallet!!!")
	in := Request{Text: "rudithu"}
	out := Response{Input: in.Text, Output: "welcome"}

	m, err := json.Marshal(in)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	n, err1 := json.Marshal(out)
	if err1 != nil {
		fmt.Println("Error marshaling JSON:", err1)
		return
	}
	fmt.Printf(string(m))
	fmt.Printf(string(n))

}
