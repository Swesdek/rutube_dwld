package interactions

import (
	"fmt"
	"maps"
	"slices"
)

func SuggestResolution(options map[string]string) string {
	fmt.Println("Available video resolutions:")
	resolutions := slices.Collect(maps.Keys(options))
	for i, resolution := range resolutions {
		fmt.Printf("%d: %s\n", i+1, resolution)
	}
	fmt.Print("Type number of resolution you want to download: ")
	var correctInput bool
	var num int
	for !correctInput {
		_, err := fmt.Scan(&num)
		if err != nil {
			fmt.Println("Incorrect number")
			continue
		}

		if num > len(options) || num < 1 {
			fmt.Println("Incorrect number")
			continue
		}

		correctInput = true
	}

	return options[resolutions[num-1]]
}
