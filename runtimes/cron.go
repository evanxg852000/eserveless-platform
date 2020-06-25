package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	envs := make(map[string]string)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "ESERVELESS_") {
			parts := strings.Split(entry, "=")
			envs[parts[0]] = parts[1]
		}
	}
	Ticker(envs)
	fmt.Println(envs)
}
