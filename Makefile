.PHONY: init

init:
	$(MAKE) -C scripts/easyvpn build
	if ! [ -f easyvpn ]; then ln scripts/easyvpn/bin/easyvpn easyvpn; fi
