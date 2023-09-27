#!/bin/sh

find . -type f \( -name "*.svelte" -o -name "*.go" \) -not -path "*front/node_modules*" | xargs wc -l
