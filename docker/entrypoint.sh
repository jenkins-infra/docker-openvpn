#!/bin/bash

set -e

OPENVPN_CONF_DIR='/etc/openvpn/server'

function ensure_required_variables {
  : "${AUTH_LDAP_PASSWORD:? AUTH_LDAP_PASSWORD required}"
  : "${AUTH_LDAP_URL:? AUTH_LDAP_URL required}"
  : "${AUTH_LDAP_BINDDN:? AUTH_LDAP_BINDDN required}"
  : "${AUTH_LDAP_GROUPS_MEMBER:? AUTH_LDAP_GROUPS_MEMBER required}"
  : "${OPENVPN_SERVER_SUBNET:? OPENVPN_SERVER_SUBNET required}"
  : "${OPENVPN_SERVER_NETMASK:? OPENVPN_SERVER_NETMASK required}"
  : "${OPENVPN_NETWORK_NAME:? OPENVPN_NETWORK_NAME required}"
}

function configure_tun {
  [ -d /dev/net ] ||
    mkdir -p /dev/net
  [ -c /dev/net/tun ] ||
    mknod /dev/net/tun c 10 200
}

function configure_certificates {
  # If custom CA are provided in /usr/local/share/ca-certificates/*.crt,
  # then this command ensures that these CAs are added to the default CA bundle
  update-ca-certificates

  if [ ! -f '/etc/openvpn/server/ca.crt' ]; then
    : "${OPENVPN_CA_PEM:? Missing OPENVPN_CA_PEM}"
    echo "$OPENVPN_CA_PEM" > /etc/openvpn/server/ca.crt
  fi

  if [ ! -f '/etc/openvpn/server/server.key' ]; then
    : "${OPENVPN_SERVER_KEY:? Missing OPENVPN_SERVER_KEY }"
    echo "$OPENVPN_SERVER_KEY" > /etc/openvpn/server/server.key
  fi

  if [ ! -f '/etc/openvpn/server/server.crt' ]; then
    : "${OPENVPN_SERVER_PEM:? Missing OPENVPN_SERVER_PEM }"
    echo "$OPENVPN_SERVER_PEM" > /etc/openvpn/server/server.crt
  fi

  if [ ! -f '/etc/openvpn/server/dh.pem' ]; then
    : "${OPENVPN_DH_PEM:? Missing OPENVPN_DH_PEM }"
    echo "$OPENVPN_DH_PEM" > /etc/openvpn/server/dh.pem
  fi
}

function copy_client_configurations_directory {
  mkdir -p /etc/openvpn/server/ccd
  cp /home/openvpn/available-ccds/${OPENVPN_NETWORK_NAME}/* /etc/openvpn/server/ccd
}

# Use ~ in order to avoid wrong interpration with / in sed command.
# Sed should be replaced by something more robust in the futur.

function configure_openvpn_server {
  sed -i "s~OPENVPN_SERVER_SUBNET~$OPENVPN_SERVER_SUBNET~g" "$OPENVPN_CONF_DIR/server.conf"
  sed -i "s~OPENVPN_SERVER_NETMASK~$OPENVPN_SERVER_NETMASK~g" "$OPENVPN_CONF_DIR/server.conf"
}

function configure_openvpn_ldap {
  sed -i "s~AUTH_LDAP_PASSWORD~$AUTH_LDAP_PASSWORD~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
  sed -i "s~AUTH_LDAP_URL~$AUTH_LDAP_URL~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
  sed -i "s~AUTH_LDAP_BINDDN~$AUTH_LDAP_BINDDN~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
  sed -i "s~AUTH_LDAP_GROUPS_MEMBER~$AUTH_LDAP_GROUPS_MEMBER~g" "$OPENVPN_CONF_DIR/auth-ldap.conf"
}

function start_openvpn {
  openvpn --config "$OPENVPN_CONF_DIR/server.conf"
}

ensure_required_variables
copy_client_configurations_directory
configure_tun
configure_certificates
configure_openvpn_server
configure_openvpn_ldap
start_openvpn
