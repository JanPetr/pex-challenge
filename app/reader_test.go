package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadURLs(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		processDuplicateURLs bool
		wantValues           []string
		wantError            bool
	}{
		{
			name:                 "non existent file",
			input:                "testdata/foobar.txt",
			processDuplicateURLs: true,
			wantValues:           nil,
			wantError:            true,
		},
		{
			name:                 "empty",
			input:                "testdata/empty.txt",
			processDuplicateURLs: true,
			wantValues:           nil,
		},
		{
			name:                 "regular",
			input:                "testdata/input.txt",
			processDuplicateURLs: true,
			wantValues: []string{
				"http://i.imgur.com/FApqk3D.jpg",
				"http://i.imgur.com/TKLs9lo.jpg",
				"https://i.redd.it/d8021b5i2moy.jpg",
				"https://i.redd.it/4m5yk8gjrtzy.jpg",
				"https://i.redd.it/xae65ypfqycy.jpg",
				"http://i.imgur.com/lcEUZHv.jpg",
				"https://i.redd.it/1nlgrn49x7ry.jpg",
				"http://i.imgur.com/M3NOzLC.jpg",
				"https://i.redd.it/w5q6gldnvcuy.jpg",
				"https://i.redd.it/s5viyluv421z.jpg",
			},
		},
		{
			name:                 "remove duplicates",
			input:                "testdata/duplicates.txt",
			processDuplicateURLs: false,
			wantValues: []string{
				"http://i.imgur.com/FApqk3D.jpg",
				"http://i.imgur.com/TKLs9lo.jpg",
				"https://i.redd.it/d8021b5i2moy.jpg",
				"https://i.redd.it/4m5yk8gjrtzy.jpg",
				"https://i.redd.it/xae65ypfqycy.jpg",
				"http://i.imgur.com/lcEUZHv.jpg",
				"https://i.redd.it/1nlgrn49x7ry.jpg",
				"http://i.imgur.com/M3NOzLC.jpg",
				"https://i.redd.it/w5q6gldnvcuy.jpg",
				"https://i.redd.it/s5viyluv421z.jpg",
			},
		},
		{
			name:                 "keep duplicates",
			input:                "testdata/duplicates.txt",
			processDuplicateURLs: true,
			wantValues: []string{
				"http://i.imgur.com/FApqk3D.jpg",
				"http://i.imgur.com/TKLs9lo.jpg",
				"https://i.redd.it/d8021b5i2moy.jpg",
				"https://i.redd.it/4m5yk8gjrtzy.jpg",
				"https://i.redd.it/xae65ypfqycy.jpg",
				"http://i.imgur.com/lcEUZHv.jpg",
				"https://i.redd.it/1nlgrn49x7ry.jpg",
				"http://i.imgur.com/M3NOzLC.jpg",
				"https://i.redd.it/w5q6gldnvcuy.jpg",
				"https://i.redd.it/s5viyluv421z.jpg",
				"http://i.imgur.com/FApqk3D.jpg",
				"http://i.imgur.com/TKLs9lo.jpg",
				"https://i.redd.it/d8021b5i2moy.jpg",
				"https://i.redd.it/4m5yk8gjrtzy.jpg",
				"https://i.redd.it/xae65ypfqycy.jpg",
				"http://i.imgur.com/lcEUZHv.jpg",
				"https://i.redd.it/1nlgrn49x7ry.jpg",
				"http://i.imgur.com/M3NOzLC.jpg",
				"https://i.redd.it/w5q6gldnvcuy.jpg",
				"https://i.redd.it/s5viyluv421z.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us, err := ReadURLs(tt.input, tt.processDuplicateURLs)

			var values []string
			for v := range us {
				values = append(values, v)
			}

			require.EqualValues(t, tt.wantValues, values)

			if tt.wantError {
				require.NotEmpty(t, err)
			}
		})
	}
}
