.PHONY: build test zip init plan apply

build:
	$(MAKE) -C go build

test:
	$(MAKE) -C go test

zip:
	./build.sh

init:
	$(MAKE) -C terraform plan

plan:
	$(MAKE) -C terraform plan

apply:
	$(MAKE) -C terraform apply