package imretro

import (
	"errors"
	"image/color"
	"testing"
)

// CompareColors helps compare colors to each other.
func CompareColors(t *testing.T, actual, want color.Color) {
	t.Helper()
	r, g, b, a := actual.RGBA()
	wr, wg, wb, wa := want.RGBA()
	comparisons := [4]channelComparison{
		{"red", r, wr},
		{"green", g, wg},
		{"blue", b, wb},
		{"alpha", a, wa},
	}

	for _, comparison := range comparisons {
		if comparison.actual != comparison.want {
			t.Errorf(
				`%s channel = %v, want %v`,
				comparison.name, comparison.actual,
				comparison.want,
			)
		}
	}
}

// ChannelComparison is used to compare color channels.
type channelComparison struct {
	name         string
	actual, want uint32
}

// CappedWriter is a writer with a fixed capacity.
type cappedWriter struct {
	len int
	cap int
}

func (w *cappedWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.len += n
	if w.len > w.cap {
		err = errors.New("Max capacity")
	}
	return
}
