#!/bin/bash
go run ../main.go --create ../demo
echo "Who is Jean-Luc Picard?" > ../demo/prompt.md
go run ../main.go \
--prompt ../demo/prompt.md \
--settings ../demo/.bob \
--output ../demo/report.md

rm -rf ../demo
