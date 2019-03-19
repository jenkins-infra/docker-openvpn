.PHONY: init

build.docker:
	$(MAKE) -C docker build

publish.docker:
	$(MAKE) -C docker publish

start:
	$(MAKE) -C docker up

init_windows:
	$(MAKE) -C scripts/easyvpn init_windows

init_linux:
	$(MAKE) -C scripts/easyvpn init_linux
	if [ -f easyvpn ]; then rm easyvpn ;fi
	ln -s $(PWD)/scripts/easyvpn/easyvpn easyvpn

init_osx:
	$(MAKE) -C scripts/easyvpn init_osx
	if [ -f easyvpn ]; then rm easyvpn ;fi
	ln -s $(PWD)/scripts/easyvpn/easyvpn easyvpn
