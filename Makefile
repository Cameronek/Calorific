run: build
	@./bin/app

build: generate
	@go build -o bin/app cmd/main.go 

generate:
	@templ generate