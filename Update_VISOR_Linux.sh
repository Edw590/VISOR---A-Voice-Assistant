#!/bin/bash

git checkout -f main
git reset --hard origin/main
git pull origin main
python./tools/download_models.py
python./tools/update_rope.py

echo
echo --------Update completed--------
echo