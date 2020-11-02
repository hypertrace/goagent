// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/hypertrace/goagent/docker/internal/container"
)

func main() {
	containerID, err := container.GetID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(containerID)
}
