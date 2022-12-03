package helpers

import (
	"reflect"
	"testing"

	"github.com/ProninIgorr/fingerprint/internal/matrix"
)

func TestLoadImage(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name string
		args args
		want *matrix.M
	}{

		{
			name: "empty file name",
			args: args{
				fname: "",
			},
			want: nil,
		},
		{
			name: "non existing file name",
			args: args{
				fname: "non-existing-file",
			},
			want: nil,
		},

		{

			name: "existing file name and wrong image format",
			args: args{
				fname: "examples/example-input-1.txt",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadImage(tt.args.fname); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
