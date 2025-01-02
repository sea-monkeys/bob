#!/bin/bash
#echo "🎉Fetching $1"
content=$(curl -s "$1")
echo "$content"
