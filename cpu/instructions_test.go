package cpu

import "testing"

func TestCls(t *testing.T) {
	chip8 := NewChip8()
	chip8.Init()

	ins := NewInstructions()

	for i := range chip8.Video {
		chip8.Video[i] = 1
	}

	ins.cls(chip8)

	for _, v := range chip8.Video {

		if v != 0 {
			t.Errorf("Video memory not cleaned")
			break
		}
	}
}
