import sys
import json
import sqlite3
import os

def main():
    _ = sys.stdin.read()
    
    response = {"status": 200, "headers": {"Content-Type": "application/json"}, "body": ""}
    
    db_path = "/mnt/data/state.db"
    
    try:
        os.makedirs("/mnt/data", exist_ok=True)
        
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        cursor.execute("CREATE TABLE IF NOT EXISTS visitors (count INTEGER)")
        
        cursor.execute("SELECT count FROM visitors")
        row = cursor.fetchone()
        
        if row is None:
            new_count = 1
            cursor.execute("INSERT INTO visitors (count) VALUES (?)", (new_count,))
        else:
            new_count = row[0] + 1
            cursor.execute("UPDATE visitors SET count = ?", (new_count,))
            
        conn.commit()
        conn.close()
        
        body = {
            "message": "Persistence Test",
            "visitor_count": new_count,
            "storage_path": db_path
        }
        response["body"] = json.dumps(body)

    except Exception as e:
        response["status"] = 500
        response["body"] = json.dumps({"error": str(e), "path_exists": os.path.exists("/mnt/data")})

    print(json.dumps(response))

if __name__ == "__main__":
    main()