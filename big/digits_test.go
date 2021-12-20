package big

import (
	"reflect"
	"testing"
)

func Test_norm(t *testing.T) {
	type args struct {
		abs digits
	}
	tests := []struct {
		name string
		args args
		want digits
	}{
		{
			name: "all zero",
			args: args{abs: digits{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: digits{},
		},
		{
			name: "upper digit zero",
			args: args{abs: digits{0, 0, 0, 1, 1, 1}},
			want: digits{1, 1, 1},
		},
		{
			name: "last digit zero",
			args: args{abs: digits{0, 0, 0, 1, 1, 0}},
			want: digits{1, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := norm(tt.args.abs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("norm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmp(t *testing.T) {
	type args struct {
		x digits
		y digits
	}
	tests := []struct {
		name string
		args args
		want int8
	}{
		{
			name: "x == y",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
				y: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: 0,
		},
		{
			name: "len(x) > len(y)",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
				y: digits{1, 2, 3, 4, 5, 6, 7, 8},
			},
			want: 1,
		},
		{
			name: "len(x) < len(y)",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8},
				y: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: -1,
		},
		{
			name: "x[i] > y[i]",
			args: args{
				x: digits{2, 2, 3, 4, 5, 6, 7, 8},
				y: digits{1, 2, 3, 4, 5, 6, 7, 8},
			},
			want: 1,
		},
		{
			name: "x[i] < y[i]",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8},
				y: digits{2, 2, 3, 4, 5, 6, 7, 8},
			},
			want: -1,
		},
		{
			name: "complex",
			args: args{
				x: digits{3, 8, 4, 5, 6, 4, 5, 3, 3, 4},
				y: digits{3, 8, 5, 2, 4, 2, 3, 4, 3, 2},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmp(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}
