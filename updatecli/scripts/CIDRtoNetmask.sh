#!/bin/bash

# This script converts a CIDR notation to a netmask.
# Function: cidr_to_netmask
# Description: This function takes a CIDR value as input and returns the corresponding netmask.
# Parameters:
#   $1 - The CIDR value (e.g., 14 from 10.0.0.0/14)

cidr_to_netmask() {
    # Calculate the netmask value from the CIDR
    value=$(( 0xffffffff ^ (1 << ( 32 - $1 )) - 1 ))
    # Print the netmask in the format xxx.xxx.xxx.xxx
    echo "$(( (value >> 24) & 0xff )).$(( (value >> 16) & 0xff )).$(( (value >> 8) & 0xff )).$(( value & 0xff ))"
}

# Check if a command line argument is provided
# If not, print an error message and exit the script
if [ $# -eq 0 ]; then
    echo "No arguments provided. Please provide a CIDR value."
    exit 1
fi

# Extract the CIDR value from the IP/CIDR string
# This allows the script to accept input in the format xxx.xxx.xxx.xxx/yy
cidr="${1#*/}"

# Call the cidr_to_netmask function with the extracted CIDR value
netmask=$(cidr_to_netmask "${cidr}")

echo "${1%/*} ${netmask}"
