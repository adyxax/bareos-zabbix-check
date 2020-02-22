package utils

import "testing"

func TestClen(t *testing.T) {
	normalString := append([]byte("abcd"), 0)
	type args struct {
		n []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty string", args{}, 0},
		{"normal string", args{normalString}, 4},
		{"non null terminated string", args{[]byte("abcd")}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clen(tt.args.n); got != tt.want {
				t.Errorf("Clen() = %v, want %v", got, tt.want)
			}
		})
	}
}
