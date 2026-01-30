package main

import (
	"github.com/pauloappbr/gojinn/sdk"
)

func main() {
	sdk.Parse()
	sdk.Log("Starting database test...")

	rows, err := sdk.DB.Query(
		"SELECT 1 as id, 'Gojinn Rocks SQLite' as message, datetime('now') as server_time",
	)

	if err != nil {
		sdk.Log("SQL error: %v", err)
		sdk.SendError(500, "Query failed: "+err.Error())
		return
	}

	sdk.Log("Query executed successfully! Returned rows: %d", len(rows))
	sdk.SendJSON(rows)
}
