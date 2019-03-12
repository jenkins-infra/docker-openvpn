.PHONY: init

init:
	$(MAKE) -C scripts/easyvpn build
	if ! [ -L easyvpn ]; then ln -s scripts/easyvpn/bin/easyvpn easyvpn; fi
