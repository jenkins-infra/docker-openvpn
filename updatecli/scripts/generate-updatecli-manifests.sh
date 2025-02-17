#!/bin/bash

# This script loops through all files in a the private certificates folder
# and for each file update the correspondings manifests depending on the vnets in comments

# Define the directory to loop through
directory="./cert/ccd/private"

# Loop through all files in the directory
for file in "$directory"/*; do
    # Check if the item is a file
    if [ -f "$file" ]; then
        file=./cert/ccd/private/smerle #debug
        echo "Processing file: $file"
        # for each couple comment/line check the corresponding manifest
        # Extract two consecutive lines
        # Skip the first line and process every two subsequent lines
        tail -n +2 "$file" | while read -r comment; do
            read -r command || break  # Read the second line, break if no more lines

            # Print the comment and execute the command
            echo "Comment: $comment"
            echo "Executing: $command"
            echo ""

            # For each comment extract the vnet name and 

        done
        break
    fi
done
