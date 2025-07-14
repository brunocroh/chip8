test:
	go run cmd/main.go $(ARGS)

test-watch:
	gow run cmd/main.go $(ARGS)

wasm:
	GOOS=js GOARCH=wasm go build -o chip8.wasm wasm/main.go
	mv chip8.wasm public/

server:
	live-server public/

.PHONY: test test-watch wasm server
