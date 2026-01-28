#!/bin/bash
# Extract expiration date from an OpenVPN certificate
set -eux -o pipefail

cert_file="${1}"
cert_type="${2:-x509}"

if [ ! -f "${cert_file}" ]; then
    echo "ERROR: Certificate file ${cert_file} not found"
    exit 1
fi

case "${cert_type}" in
    "x509")
        # Extract the notAfter date from the certificate
        # Output format: notAfter=Jan 15 12:34:56 2026 GMT
        expiry_raw="$(openssl x509 -enddate -noout -in "${cert_file}" 2>/dev/null | cut -d= -f2)"
        ;;
    "crl")
        # Extract the nextUpdate date from the certificate
        # Output format: nextUpdate=Jun  1 10:20:48 2026 GMT
        expiry_raw="$(openssl crl --nextupdate -noout -in "${cert_file}" 2>/dev/null | cut -d= -f2)"
        ;;
    *)
        echo "ERROR: unsupported certificate type ${cert_type}".
        exit 1
        ;;
esac

if [ -z "${expiry_raw}" ]; then
    echo "ERROR: Could not extract expiration date from ${cert_file}"
    exit 1
fi

# Convert to ISO 8601 format for easier parsing
DATE_BIN='date'
if command -v gdate >/dev/null 2>&1; then
    DATE_BIN='gdate'
fi

command -v "${DATE_BIN}" >/dev/null 2>&1 || { echo "ERROR: ${DATE_BIN} command not found. Exiting."; exit 1; }

# Convert to ISO format: YYYY-MM-DDTHH:MM:SSZ
expiry_iso=$("${DATE_BIN}" --utc -d "${expiry_raw}" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null)

echo "${expiry_iso}"
