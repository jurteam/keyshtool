#!/usr/bin/make -f

all: keyshtool
keyshtool:
	go build

check: keyshtool
	mkdir -p $@
	bash -x ./run-tests.sh

clean:
	rm -fv keyshtool
	rm -rfv TESTS/ check/

.PHOMY: clean all
