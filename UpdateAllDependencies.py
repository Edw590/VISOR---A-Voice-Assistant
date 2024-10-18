import os
import subprocess

for root, dirs, files in os.walk('.'):
    if "go.mod" in files:
        print(f"go.mod in directory: {root}")
        subprocess.run(["go", "get", "-u", "./..."], cwd=root)
        subprocess.run(["go", "mod", "tidy"], cwd=root)
