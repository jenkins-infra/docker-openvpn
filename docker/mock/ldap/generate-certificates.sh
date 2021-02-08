#!/bin/bash

set -eux -o pipefail

## Generate certificates in current dir

## Generate LDAP Certificate + Key
openssl req \
  -newkey rsa:2048 \
  -nodes \
  -keyout ./ldap.key \
  -out ./ldap.csr\
  -subj "/C=BE/O=JENKINSPROJECT/CN=ldap"

# Generate Certificate Authority Certificate + Key
openssl req \
  -newkey rsa:2048 \
  -nodes \
  -keyout ./ca.key \
  -x509 \
  -days 365 \
  -out ./ca.crt\
  -subj "/C=BE/O=JENKINSPROJECT/CN=ldap"

## Sign LDAP Certificate with Certificate Authority
openssl x509 -req \
  -in ./ldap.csr \
  -CA ./ca.crt \
  -CAkey ./ca.key \
  -CAcreateserial \
  -CAserial ./ca.srl \
  -out ./ldap.crt
