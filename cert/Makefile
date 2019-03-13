.PHONY: clean encrypt decrypt isNameDefined request revoke show-req show-cert sign

PKI_DIR=pki

clean: 
	rm $(PKI_DIR)/private/ca.key

encrypt:
	sops -e $(PKI_DIR)/private/ca.key > $(PKI_DIR)/private/ca.key.enc

decrypt:
	sops -d $(PKI_DIR)/private/ca.key.enc > $(PKI_DIR)/private/ca.key

isNameDefined:
	if [ -z "$(name)" ]; then echo "Missing argument name=<username>"; exit 1 ; fi

request: isNameDefined
	./easyrsa gen-req $(name)
	git add $(PKI_DIR)/reqs/$(name).req
	git commit $(PKI_DIR)/reqs/$(name).req -m 'Submit certificate request for $(name)'
	git push origin

revoke: isNameDefined decrypt
	./easyrsa revoke $(name)
	./easyrsa gen-crl
	git add \
		$(PKI_DIR)/crl.pem \
		$(PKI_DIR)/index.txt \
		$(PKI_DIR)/index.txt.attr \
		$(PKI_DIR)/reqs \
		$(PKI_DIR)/certs_by_serial \
		$(PKI_DIR)/issued $(PKI_DIR)/revoked
	rm $(PKI_DIR)/index.txt.old
	git commit -m "Revoked cert for $(name)"
	git push origin
	make clean

show-req: isNameDefined
	./easyrsa show-req $(name)

show-cert: isNameDefined
	./easyrsa show-cert $(name)

sign: isNameDefined decrypt
	./easyrsa sign-req client $(name)
	rm $(PKI_DIR)/index.txt.old $(PKI_DIR)/index.txt.attr $(PKI_DIR)/serial.old $(PKI_DIR)/extensions.temp
	git add $(PKI_DIR)/issued/$(name).crt $(PKI_DIR)/index.txt $(PKI_DIR)/certs_by_serial/ $(PKI_DIR)/serial
	git commit $(PKI_DIR)/issued/$(name).crt $(PKI_DIR)/index.txt $(PKI_DIR)/serial $(PKI_DIR)/certs_by_serial/ -m 'Sign certificate request for $(name)'
	git push origin
	make clean

# Missing Revoking procedure
