const input = document.querySelector("#load-rom-input");
const canvas = document.getElementById("canvas");
canvas.width = 1024;
canvas.height = 512;
const ctx = canvas.getContext("2d");
const go = new Go();

const clearCanvas = () => {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
};

const renderCallback = (_video) => {
  const video = convertToUint32Array(_video);
  clearCanvas();

  for (i = 0; i < video.length; i++) {
    if (video[i]) {
      ctx.fillStyle = "white";
    } else {
      ctx.fillStyle = "black";
    }

    const x = (i % 64) * 16;
    const y = (i / 64) * 16;
    const height = 16;
    const width = height;

    ctx.fillRect(x, y, width, height);
  }
};

function convertToUint32Array(uint8Array) {
  const uint32Array = new Uint32Array(2048);
  const dataView = new DataView(uint8Array.buffer);
  for (let i = 0; i < 2048; i++) {
    uint32Array[i] = dataView.getUint32(i * 4, true);
  }
  return uint32Array;
}

const bufferMemory = new ArrayBuffer(8192);
const videoMemory = new Uint8Array(bufferMemory);

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
        window.start(() => renderCallback(videoMemory), videoMemory);
      };
    });
  },
);
