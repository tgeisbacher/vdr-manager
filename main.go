package main

import (
	"fmt"

	"github.com/tgeisbacher/vdr-manager/svdrp"
)

func main() {

	channels, err := svdrp.ListAllChannels("10.0.0.198:6419")

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	fmt.Println(channels)
}
