const input = document.querySelector("#load-rom-input");
const go = new Go();
WebAssembly.instantiateStreaming(fetch("chip8.wasm"), go.importObject).then(
  (wasm) => {
    const { instance } = wasm;
    go.run(wasm.instance);

    addEventListener("input", (event) => {
      const fileReader = new FileReader();
      fileReader.readAsArrayBuffer(input.files[0]);
      fileReader.onload = () => {
        const rom = new Uint8Array(fileReader.result);
        window.loadRom(rom);
        // window.start();
      };
    });
  },
);
