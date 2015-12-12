default:	
	$(MAKE) all
deps:
	bash -c "go get github.com/tools/godep && godep restore"
test:
	bash -c "./scripts/test.sh $(TEST)"
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
