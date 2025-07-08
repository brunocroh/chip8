test:
	go run cmd/main.go $(ARGS)

test-watch:
	gow run cmd/main.go $(ARGS)

.PHONY: test
