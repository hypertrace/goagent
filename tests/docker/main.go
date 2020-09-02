// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/traceableai/goagent/docker/internal"
)

func main() {
	containerID, err := internal.GetContainerID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(containerID)
}
