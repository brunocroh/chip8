run:
	go run cmd/main.go $(ARGS)

run-watch:
	gow run cmd/main.go $(ARGS)

test:
	go test -v --cover ./...

build:
	go build -o chip8 cmd/main.go 

wasm:
	GOOS=js GOARCH=wasm go build -o chip8.wasm wasm/wasm.go
	mv chip8.wasm public/

server:
	live-server public/

.PHONY: run test run-watch wasm server
