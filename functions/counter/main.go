package main

import (
	"fmt"
	"strconv"

	"github.com/pauloappbr/gojinn/sdk"
)

func main() {
	val, found := sdk.KV.Get("visit_count")

	count := 0
	if found {
		count, _ = strconv.Atoi(val)
	}

	count++
	newVal := fmt.Sprintf("%d", count)

	sdk.KV.Set("visit_count", newVal)

	sdk.Log("Visit number: %s", newVal)

	sdk.SendJSON(map[string]interface{}{
		"message": "Global Counter",
		"visits":  count,
	})
}
