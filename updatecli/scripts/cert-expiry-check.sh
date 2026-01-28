#!/bin/bash
# Check if certificate expires within 30 days
set -eux -o pipefail

currentexpirydate="${1}"

DATE_BIN='date'
## non GNU operating system
if command -v gdate >/dev/null 2>&1; then
    DATE_BIN='gdate'
fi

command -v "${DATE_BIN}" >/dev/null 2>&1 || { echo "ERROR: ${DATE_BIN} command not found. Exiting."; exit 1; }

currentdateepoch=$("${DATE_BIN}" --utc "+%s" 2>/dev/null)
expirydateepoch=$("${DATE_BIN}" "+%s" -d "${currentexpirydate}")
datediff=$(((expirydateepoch-currentdateepoch)/(60*60*24))) # diff per days

echo "Certificate expires in ${datediff} days"

if [ "${datediff}" -lt 30 ]; then # Alert 30 days before expiration
    echo "Certificate expiring soon - action required"
    exit 0
else
    echo "Certificate not expiring soon"
    exit 1
fi
