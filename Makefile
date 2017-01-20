default:
	docker build -t sharebuilder build/.
	docker run --rm -v $(GOPATH):$(GOPATH) -e GOPATH=$(GOPATH) -w `pwd` sharebuilder ./build/build.bash
