.PHONY: build run shell

IMAGE = 'olblak/openvpn'
TAG = $(shell git rev-parse HEAD | cut -c1-6)

OPENVPN_CONF_DIR="certs"

build:
	docker build --no-cache -t $(IMAGE):$(TAG) .

publish:
	docker push $(IMAGE):$(TAG)

gen_cert: gen_cert_vpn gen_cert_ldap

gen_cert_vpn:
	mkdir $(OPENVPN_CONF_DIR) || true
	openssl dhparam -out $(OPENVPN_CONF_DIR)/dh.pem 2048
	openssl req \
	   -newkey rsa:2048 \
	   -nodes \
	   -keyout $(OPENVPN_CONF_DIR)/server.key \
	   -out $(OPENVPN_CONF_DIR)/server.csr\
	   -subj "/C=BE/ST=**/L=**/O=**/CN=vpn/emailAddress=**"
	openssl req \
	   -newkey rsa:2048 \
	   -nodes \
	   -keyout $(OPENVPN_CONF_DIR)/ca.key \
	   -x509 \
	   -days 365 \
	   -out $(OPENVPN_CONF_DIR)/ca.pem\
	   -subj "/C=BE/ST=**/L=**/O=**/CN=CA/emailAddress=**"
	openssl x509 -req \
	    -in $(OPENVPN_CONF_DIR)/server.csr \
	    -CA $(OPENVPN_CONF_DIR)/ca.pem \
	    -CAkey $(OPENVPN_CONF_DIR)/ca.key \
	    -CAcreateserial \
	    -out $(OPENVPN_CONF_DIR)/server.pem
	openssl verify -CAfile $(OPENVPN_CONF_DIR)/ca.pem $(OPENVPN_CONF_DIR)//server.pem
	chown -R 101:101 $(OPENVPN_CONF_DIR)

gen_cert_ldap:
	mkdir ssl || true
	openssl req \
       -newkey rsa:2048 \
       -nodes \
       -keyout $(OPENVPN_CONF_DIR)/ldap.key \
       -out $(OPENVPN_CONF_DIR)/ldap.csr\
       -subj "/C=BE/O=JENKINSPROJECT/CN=ldap"
	openssl x509 -req \
        -in $(OPENVPN_CONF_DIR)/ldap.csr \
        -CA $(OPENVPN_CONF_DIR)/ca.pem \
        -CAkey $(OPENVPN_CONF_DIR)/ca.key \
        -CAcreateserial \
        -out $(OPENVPN_CONF_DIR)/ldap.pem\
	chown -R 101:101 $(OPENVPN_CONF_DIR)

up: 
	docker-compose build
	docker-compose up -d
