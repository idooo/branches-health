#!/usr/bin/env bash

DATA=$(cat ./assets/index.html)

echo "package core" > ./core/template.go
echo "const IndexTemplate = \`$DATA\`" >> ./core/template.go