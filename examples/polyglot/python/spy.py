import sys
import json
import os

def main():
    _ = sys.stdin.read()
    response = {"status": 200, "headers": {"Content-Type": "application/json"}, "body": ""}
    
    report = {}

    try:
        report["root_ls"] = os.listdir('/')
    except Exception as e:
        report["root_ls"] = f"BLOCKED/EMPTY ({str(e)})"

    try:
        report["cwd_ls"] = os.listdir('/app')
    except Exception as e:
        report["cwd_ls"] = str(e)

    try:
        with open('/mnt/secret.txt', 'r') as f:
            report["secret"] = f.read().strip()
    except Exception as e:
        report["secret"] = f"ACCESS DENIED ({str(e)})"

    try:
        with open('/etc/passwd', 'r') as f:
            report["host_attack"] = "CRITICAL FAIL: I can read /etc/passwd"
    except:
        report["host_attack"] = "SUCCESS: Host file blocked"

    response["body"] = json.dumps(report, indent=2)
    print(json.dumps(response))

if __name__ == "__main__":
    main()