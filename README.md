# CHIP-8 Emulator

A CHIP-8 emulator/interpreter written in Go with both desktop and web support using WebAssembly.

## About CHIP-8

CHIP-8 is an interpreted programming language developed by Joseph Weisbecker in the mid-1970s. It was initially used on the COSMAC VIP and Telmac 1800 8-bit microcomputers. CHIP-8 programs are run on a virtual machine that features:

- 4096 bytes of memory
- 16 8-bit general-purpose registers (V0-VF)
- 64x32 monochrome display
- 16-key hexadecimal keypad
- Two timers (delay and sound)

## Features

- Full CHIP-8 instruction set implementation
- SDL2-based desktop application with graphics (audio need to be implemented yet)
- WebAssembly build for browser execution

## Project Structure

```
chip8/
├── cmd/           # Desktop application entry point
├── cpu/           # CHIP-8 CPU implementation
│   ├── cpu.go     # Main CPU structure and methods
│   ├── decoder.go # Instruction decoding logic
│   ├── instructions.go # Opcode implementations
│   └── timers.go  # Timer management
├── utils/         # Utility functions
│   └── rom.go     # ROM loading utilities
├── wasm/          # WebAssembly entry point
├── public/        # Web assets
│   ├── index.html # Web interface
│   ├── index.js   # JavaScript bridge
│   └── chip8.wasm # Compiled WebAssembly module
└── roms/          # Sample ROM files
```

## Dependencies

- **Go 1.23.6+** - Programming language
- **SDL2** - Graphics and input handling (desktop version)
  - macOS: `brew install sdl2`
  - Ubuntu/Debian: `sudo apt-get install libsdl2-dev`
  - Windows: Download from [libsdl.org](https://www.libsdl.org/)

## Building and Running

### Desktop Application

1. Install SDL2 development libraries
2. Dowload roms from internet, [this repo](https://github.com/dmatlack/chip8) have a bunch of roms to download, save into roms folder
3. Build and run with a ROM file:

using the Makefile:

```bash
make run ARGS="roms/<ROM_NAME>.ch8"
```

For development with auto-reload:

```bash
make run-watch ARGS="roms/<ROM_NAME>.ch8"
```

### WebAssembly Version

1. Build the WebAssembly module:

```bash
make wasm
```

2. install live-server do serve files:

```bash
npm install -g live-server
```

3. Start a local server:

```bash
make server
```

4. Open your browser to the served address and load a ROM file through the file input.

## Controls

The CHIP-8 keypad is mapped to your keyboard as follows:

```
CHIP-8 Keypad    Keyboard
1 2 3 C          1 2 3 4
4 5 6 D    =>    Q W E R
7 8 9 E          A S D F
A 0 B F          Z X C V
```

## ROM Compatibility

Additional ROMs can be found at:

- [CHIP-8 Archive](https://johnearnest.github.io/chip8Archive/)
- [CHIP-8 Test Suite](https://github.com/Timendus/chip8-test-suite)

## Implementation Details

### CPU Architecture

- 16 general-purpose 8-bit registers (V0-VF)
- Program counter and index register
- Stack for subroutine calls
- 64x32 pixel monochrome display buffer
- 16-key input state tracking

### Instruction Set

Implements all 35 standard CHIP-8 instructions including:

- Arithmetic and logic operations
- Memory operations
- Display operations
- Flow control
- Timer operations
- Input handling

### Timing

- CPU runs at 500 Hz (configurable)
- Timers update at 60 Hz
- Display refresh on draw flag

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

The emulator is organized into clear modules:

- `cpu/` contains the core emulation logic
- `cmd/` contains the desktop application
- `wasm/` contains the WebAssembly bridge
- `utils/` contains shared utilities

## References

- [CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [CHIP-8 Wikipedia](https://en.wikipedia.org/wiki/CHIP-8)
- [Tobias V. Langhoff High level guide](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)

## License

This project is open source and available under the MIT License.

