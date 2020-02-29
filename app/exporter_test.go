package app

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportCSV(t *testing.T) {
	outputFile := "testdata/out/output.csv"

	testAnalyzedImages := []AnalyzedImage{
		{"foo", []string{"1", "2", "3"}},
		{"bar", []string{"4", "5", "6"}},
		{"baz", []string{"7", "8", "9"}},
	}

	expected := []string{"foo,1,2,3", "bar,4,5,6", "baz,7,8,9", ""}

	analyzedImages := make(chan AnalyzedImage, len(testAnalyzedImages))
	for _, tai := range testAnalyzedImages {
		analyzedImages <- tai
	}
	close(analyzedImages)

	done, erre := ExportCSV(outputFile, analyzedImages)
	select {
	case <-done:
		of, err := os.Open(outputFile)
		require.NoError(t, err)

		b, err := ioutil.ReadAll(of)
		require.NoError(t, err)

		lines := strings.Split(string(b), "\n")
		require.EqualValues(t, expected, lines)
		return
	case err := <-erre:
		require.NoError(t, err)
		return
	}
}
