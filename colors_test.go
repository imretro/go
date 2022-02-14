package imretro

import "testing"

// Tests that a single byte would be grayscale.
func TestGrayscaleBytes(t *testing.T) {
	c := colorBytes{0x80}
	tests := []struct {
		channel rune
		want    uint32
	}{
		{'r', 0x8080},
		{'g', 0x8080},
		{'b', 0x8080},
		{'a', 0xFFFF},
	}
	r, g, b, a := c.RGBA()
	channels := []uint32{r, g, b, a}

	for i, tt := range tests {
		actual := channels[i]
		if actual != tt.want {
			t.Errorf(`%c = %04X, want %04X`, tt.channel, actual, tt.want)
		}
	}
}

// Tests that 3 bytes would be RGB.
func TestRGBBytes(t *testing.T) {
	c := colorBytes{0x55, 0x80, 0xAA}
	tests := []struct {
		channel rune
		want    uint32
	}{
		{'r', 0x5555},
		{'g', 0x8080},
		{'b', 0xAAAA},
		{'a', 0xFFFF},
	}
	r, g, b, a := c.RGBA()
	channels := []uint32{r, g, b, a}

	for i, tt := range tests {
		actual := channels[i]
		if actual != tt.want {
			t.Errorf(`%c = %04X, want %04X`, tt.channel, actual, tt.want)
		}
	}
}

// Tests that 4 bytes would be RGBA.
func TestRGBABytes(t *testing.T) {
	c := colorBytes{0x55, 0x80, 0xAA, 0xCC}
	tests := []struct {
		channel rune
		want    uint32
	}{
		{'r', 0x5555},
		{'g', 0x8080},
		{'b', 0xAAAA},
		{'a', 0xCCCC},
	}
	r, g, b, a := c.RGBA()
	channels := []uint32{r, g, b, a}

	for i, tt := range tests {
		actual := channels[i]
		if actual != tt.want {
			t.Errorf(`%c = %04X, want %04X`, tt.channel, actual, tt.want)
		}
	}
}

// Asserts that unreachable code panics.
func TestColorBytesUnreachable(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Invalid length did not cause a panic")
		}
	}()
	c := colorBytes{0, 0, 0, 0, 1, 2, 3}
	c.RGBA()
}
