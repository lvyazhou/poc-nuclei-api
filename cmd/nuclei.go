package main

import (
	"context"
	"fmt"
	"nuclei/pkg/nuclei"

	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

func main() {
	urls := []string{"http://218.77.53.28:8087"}

	options := &types.Options{}
	if wrapper, err := nuclei.NewWrapper(options); err != nil {
		panic(err)
	} else {
		results, scanErr := wrapper.RunEnumeration(context.Background(), urls)
		if scanErr != nil {
			panic(scanErr)
		}

		if len(results) > 0 {
			for _, r := range results {
				fmt.Println(r)
			}
		}
	}
}
