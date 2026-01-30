package main

import (
	"time"

	"github.com/pauloappbr/gojinn/sdk"
)

func main() {
	req, err := sdk.Parse()
	if err != nil {
		sdk.SendError(500, "Failed to read request")
		return
	}

	sdk.Log("Received HTMX call on method: %s", req.Method)

	agora := time.Now().Format("15:04:05")

	html := `
		<div class="clock-widget" style="border: 2px solid #4CAF50; padding: 20px; border-radius: 10px; text-align: center;">
			<h2 style="color: #333;">Correct Time (Gojinn)</h2>
			<p style="font-size: 2em; font-weight: bold; color: #4CAF50;">` + agora + `</p>
			<small>Rendered on server via WASM</small>
		</div>
	`

	sdk.SendHTML(html)
}
