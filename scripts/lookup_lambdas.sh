#!/bin/sh

# Define the root directory to search under
services_dir="./services"
ws_dir="./services/websocket"

# Temporary file to store JSON entries
json_file=$(mktemp)

# Find directories named "lambda" under services_dir and create JSON entries
find "$services_dir" -type d -name "lambda" | while IFS= read -r lambda_path; do
    service_path=$(dirname "$lambda_path")
    service_name=$(basename "$service_path")

    # Create a JSON object and append it to the temporary file
    printf '"%s": "%s"\n' "$service_name" "$lambda_path" >> "$json_file"

done


# Find websocket lambdas
find "$ws_dir" -type d -name "*" -depth 1 | while IFS= read -r lambda_path; do
    service_path=$(dirname "$lambda_path")
    # As the main of these websocker are lambda
    service_name="$(basename "$ws_dir")/$(basename "$lambda_path")"

    # Create a JSON object and append it to the temporary file
    printf '"%s": "%s"\n' "$service_name" "$lambda_path" >> "$json_file"

done



# Concatenate all JSON entries into a single JSON object
json=$(cat "$json_file" | paste -sd "," -)
json="{ $json }"

echo $json

# Clean up temporary file
rm "$json_file"