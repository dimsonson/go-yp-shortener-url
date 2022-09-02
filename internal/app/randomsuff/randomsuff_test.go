package randomsuff

import (
	"testing"
)

func TestRandSeq(t *testing.T) {

	//rnd:= RandSeq(args.n)
	tests := []struct {
		name    string
		args    int
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "zero",
			args:    0,
			wantErr: true,
		},
		{
			name:    "negative",
			args:    -1,
			wantErr: true,
		},
		{
			name:    "one",
			args:    1,
			wantErr: true,
		},
		{
			name:    "more than one number",
			args:    5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RandSeq(tt.args); err != nil != tt.wantErr {
				t.Errorf("RandSeq() = %v, want %v", got, tt.want)
			}

		})
	}
}
