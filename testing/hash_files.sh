#!/bin/bash

HASH_MOD="sha384"

# Directories containing assets
DIRS=(
  "/Users/main/Projects/Go/wastetags/cmd/wastetags/assets/css"
  "/Users/main/Projects/Go/wastetags/cmd/wastetags/assets/js"
  "/Users/main/Projects/Go/wastetags/cmd/wastetags/assets/images"
)

# Directory containing HTML files
TEMPLATE_DIR="/Users/main/Projects/Go/wastetags/cmd/wastetags/templates"

# Temporary file to store hashes
HASH_FILE=$(mktemp)

# Step 1: Generate hashes for asset files
for dir in "${DIRS[@]}"; do
  for file in "$dir"/*; do
    if [ -f "$file" ]; then
      # Compute the hash
      hash=$(cat "$file" | openssl dgst -$HASH_MOD -binary | openssl base64 -A)
      filename=$(basename "$file")
      # Save the filename and hash
      echo "$filename $HASH_MOD-$hash" >> "$HASH_FILE"
    fi
  done
done

# Step 2: Update all HTML files in the template directory
for html_file in "$TEMPLATE_DIR"/*.html; do
  while read -r line; do
    # Extract filename and hash
    filename=$(echo "$line" | awk '{print $1}')
    hash=$(echo "$line" | awk '{print $2}')

    # Update <link> and <script> tags
    sed -i '' "s|\\(href=.*$filename\\)|\\1 integrity=\"$escaped_hash\" crossorigin=\"anonymous\"|" "$html_file"
    sed -i '' "s|\\(src=.*$filename\\)|\\1 integrity=\"$escaped_hash\" crossorigin=\"anonymous\"|" "$html_file"
  done < "$HASH_FILE"
done

# Clean up temporary hash file
rm "$HASH_FILE"

echo "Integrity hashes added to all HTML files in $TEMPLATE_DIR."
