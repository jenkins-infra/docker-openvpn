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
	ln -f -s $(PWD)/scripts/easyvpn/easyvpn easyvpn

init_osx:
	$(MAKE) -C scripts/easyvpn init_osx
	ln -f -s $(PWD)/scripts/easyvpn/easyvpn easyvpn
