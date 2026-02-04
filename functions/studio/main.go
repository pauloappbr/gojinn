package main

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func main() {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Gojinn Studio</title>
    <style>
        body { font-family: 'Courier New', monospace; background: #0d1117; color: #c9d1d9; padding: 20px; display: flex; flex-direction: column; align-items: center; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; width: 100%; max-width: 800px; margin-top: 20px; }
        .card { background: #161b22; border: 1px solid #30363d; padding: 20px; border-radius: 6px; }
        h1 { color: #58a6ff; width: 100%; max-width: 800px; border-bottom: 1px solid #30363d; }
        h2 { margin-top: 0; color: #79c0ff; }
        
        .btn { background: #238636; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; font-family: inherit; font-weight: bold; }
        .btn:hover { background: #2ea043; }
        .btn-warn { background: #da3633; }
        .input-group { display: flex; gap: 10px; align-items: center; margin-top: 10px; }
        input { background: #0d1117; border: 1px solid #30363d; color: white; padding: 5px; border-radius: 4px; width: 60px; text-align: center; }
        
        .stat-val { font-weight: bold; color: #fff; }
    </style>
</head>
<body>
    <h1>☁️ Gojinn Sovereign Studio</h1>
    
    <div class="grid">
        <div class="card">
            <h2>🕸️ Topology</h2>
            <div>Node ID: <span id="node-id">...</span></div>
            <div>Peers: <span id="peer-count" class="stat-val">0</span></div>
        </div>

        <div class="card">
            <h2>🔥 Hot Control</h2>
            <div>Current Pool Size: <span id="pool-size" class="stat-val">-</span></div>
            
            <div class="input-group">
                <label>Scale Workers:</label>
                <input type="number" id="target-pool" value="5" min="1" max="100">
                <button class="btn" onclick="patchPool()">Apply Patch</button>
            </div>
            <div style="margin-top: 10px; font-size: 0.8em; color: #8b949e;">
                * Updates variable in-memory immediately.
            </div>
        </div>
    </div>

    <script>
        async function update() {
            try {
                const res = await fetch('/_sys/status');
                const data = await res.json();
                document.getElementById('node-id').innerText = data.node_id;
                document.getElementById('peer-count').innerText = (data.active_peers || []).length;
                document.getElementById('pool-size').innerText = data.pool_size;
            } catch(e) {}
        }

        async function patchPool() {
            const val = parseInt(document.getElementById('target-pool').value);
            try {
                const res = await fetch('/_sys/patch', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({ pool_size: val })
                });
                if (res.ok) {
                    const data = await res.json();
                    alert("✅ " + data.msg);
                    update(); // Refresh UI immediately
                } else {
                    alert("❌ Patch Failed");
                }
            } catch(e) { alert("Error: " + e); }
        }

        update();
        setInterval(update, 2000);
    </script>
</body>
</html>`

	resp := Response{
		Status:  200,
		Headers: map[string][]string{"Content-Type": {"text/html"}},
		Body:    html,
	}
	jsonBytes, _ := json.Marshal(resp)
	fmt.Println(string(jsonBytes))
}
