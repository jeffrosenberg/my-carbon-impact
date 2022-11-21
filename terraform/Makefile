.PHONY: init validate plan apply

env := dev

init:
	terraform init

validate:
	terraform validate

plan:
	terraform plan -var-file="$(env).tfvars"

apply:
	terraform apply -var-file="$(env).tfvars"