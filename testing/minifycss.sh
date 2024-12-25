#!/bin/bash
for file in ./input/*.css; do
  filename=$(basename "$file" .css)
  cleancss -o "./output/${filename}.min.css" "$file"
done