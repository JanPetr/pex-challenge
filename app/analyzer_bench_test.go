package app

import (
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkAnalyze(b *testing.B) {
	f, err := os.Open("testdata/images/1nlgrn49x7ry.jpg")
	require.NoError(b, err)

	img, _, err := image.Decode(f)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := Analyze(img)
		require.NoError(b, err)
	}
}
