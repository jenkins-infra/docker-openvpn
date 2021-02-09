.PHONY: init_linux init_windows init_osx start publish.docker build.docker

build.docker:
	$(MAKE) -C docker build

publish.docker:
	$(MAKE) -C docker publish

start:
	$(MAKE) -C docker up
	# Sleeping is good for your health buddy!
	sleep 5
	$(MAKE) -C docker client-connect

init_windows:
	$(MAKE) -C utils/easyvpn init_windows

init_linux:
	$(MAKE) -C utils/easyvpn init_linux
	ln -f -s $(PWD)/utils/easyvpn/easyvpn easyvpn

init_osx:
	$(MAKE) -C utils/easyvpn init_osx
	ln -f -s $(PWD)/utils/easyvpn/easyvpn easyvpn
