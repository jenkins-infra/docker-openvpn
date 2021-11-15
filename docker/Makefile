.PHONY: client-connect up

up:
	docker-compose --build --force-recreate up -d
	docker-compose exec ldap /entrypoint/restore

client-connect:
	docker-compose exec --workdir=/root/openvpn-client vpn openvpn --config config.ovpn --auth-user-pass auth
