default: build

# note: not using $GOPATH/src to use local imports in the poc apps
build_container:
	docker run --rm -v "$(CURDIR)":/go/src/github.com/kcq/poc-ipblock-pool -w /go/src/github.com/kcq/poc-ipblock-pool golang:1.11.4 make build
	docker build -t poc/ipblock-pool $(CURDIR)

build:
	'$(CURDIR)/scripts/build.sh'

build_with_local_gopath:
	'$(CURDIR)/scripts/local_gopath.sh'
	. '$(CURDIR)/scripts/env.sh'; '$(CURDIR)/scripts/build.sh'

clean:
	'$(CURDIR)/scripts/clean.sh'

fmt:
	'$(CURDIR)/scripts/fmt.sh'

consul_run:
	docker run -it --rm --name="consul_only" -p 8500:8500 consul:1.4.2 agent -dev -ui -client 0.0.0.0

.PHONY: default build_container build_with_local_gopath build clean fmt consul_run