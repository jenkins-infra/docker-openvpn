#!/bin/bash
# Return netmask for a given network and CIDR.
# Convert CIDR suffix to netmask
cidr_to_netmask() {
    local cidr=$1
    local mask=(0 0 0 0)
    for ((i=0; i<cidr; i++)); do
        mask[i/8]=$((mask[i/8] + (1 << (7 - i % 8))))
    done
    echo "${mask[0]}.${mask[1]}.${mask[2]}.${mask[3]}"
}
# Extract IP address and CIDR suffix
ip=$(cut -d'/' -f1 <<< "$1")
suffix=$(cut -d'/' -f2 <<< "$1")

# Convert CIDR suffix to netmask
netmask=$(cidr_to_netmask "${suffix}")

# echo result
echo "${ip} ${netmask}"
