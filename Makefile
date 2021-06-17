
.PHONY: clean
clean: ## clean logs
	@rm -f log/*.log

.PHONY: test
test: ## Test all files with unit mode
	-@go test ./...

.PHONY: example
example: ## Run example
	@go run example/example.go