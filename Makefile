.PHONY: all

CURDIR := $(shell pwd)

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

serve_helloworld: ## Start hello world application
	npx parcel ./helloworld

serve_sessioncookie: ## Start session cookie sample application
	npx parcel build ./sessioncookie/index.html
	cd ./sessioncookie && go run main.go

serve_usermanagement: ## Start user management sample application
	npx parcel build ./usermanagement/index.html -d ./usermanagement/dist --public-url /static
	go run ./usermanagement/main.go

serve_sessionmanagement: ## Start session management sample application
	npx parcel build ./sessionmanagement/index.html -d ./sessionmanagement/dist --public-url /static
	go run ./sessionmanagement/session.go ./sessionmanagement/main.go
