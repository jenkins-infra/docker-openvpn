.PHONY: init

build.docker:
	$(MAKE) -C docker build

start:
	$(MAKE) -C docker up

init_windows:
	$(MAKE) -C scripts/easyvpn init_windows

init_linux:
	$(MAKE) -C scripts/easyvpn init_linux
	ln -s $(PWD)/scripts/easyvpn/easyvpn easyvpn

init_osx:
	$(MAKE) -C scripts/easyvpn init_osx
	ln -s $(PWD)/scripts/easyvpn/easyvpn easyvpn
