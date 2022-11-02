package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	info, err := ParseApk(os.Args[1], true)
	if err != nil {
		fmt.Println(err)
	} else {
		bytes, _ := json.MarshalIndent(info, "", "\t")
		fmt.Println(string(bytes))
	}

}
