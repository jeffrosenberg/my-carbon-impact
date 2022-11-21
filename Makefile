.PHONY: generate build test zip init plan apply deploy

generate:
	$(MAKE) -C go generate

build:
	$(MAKE) -C go build

test:
	$(MAKE) -C go test

zip:
	./build.sh

init:
	$(MAKE) -C terraform plan

validate:
	$(MAKE) -C terraform validate

plan:
	$(MAKE) -C terraform plan

apply:
	$(MAKE) -C terraform apply

deploy: test zip apply