#!/bin/bash

set -e

OPENVPN_CONF_DIR='/etc/openvpn/server'

: "${AUTH_LDAP_PASSWORD:? AUTH_LDAP_PASSWORD required}"
: "${AUTH_LDAP_URL:? AUTH_LDAP_URL required}"
: "${AUTH_LDAP_BINDDN:? AUTH_LDAP_BINDDN required}"
: "${AUTH_LDAP_GROUPS_MEMBER:? AUTH_LDAP_GROUPS_MEMBER required}"

# Use ~ in order to avoid wrong interpration with / in sed command.
# Sed should be replaced by something more robust in the futur.

sed -i "s~AUTH_LDAP_PASSWORD~$AUTH_LDAP_PASSWORD~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
sed -i "s~AUTH_LDAP_URL~$AUTH_LDAP_URL~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
sed -i "s~AUTH_LDAP_BINDDN~$AUTH_LDAP_BINDDN~g" "$OPENVPN_CONF_DIR/auth-ldap.conf" 
sed -i "s~AUTH_LDAP_GROUPS_MEMBER~$AUTH_LDAP_GROUPS_MEMBER~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"

openvpn --config "$OPENVPN_CONF_DIR/server.conf"
