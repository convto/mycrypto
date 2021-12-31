package big

import (
	"reflect"
	"testing"
)

func Test_karatsubaLen(t *testing.T) {
	type args struct {
		n         int
		threshold int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "n = 40",
			args: args{
				n:         40,
				threshold: karatsubaThreshold,
			},
			want: 40,
		},
		{
			name: "n = 200",
			args: args{
				n:         200,
				threshold: karatsubaThreshold,
			},
			want: 200,
		},
		{
			name: "n = 500",
			args: args{
				n:         500,
				threshold: karatsubaThreshold,
			},
			want: 512,
		},
		{
			name: "n = 900",
			args: args{
				n:         900,
				threshold: karatsubaThreshold,
			},
			want: 928,
		},
		{
			name: "n = 999",
			args: args{
				n:         999,
				threshold: karatsubaThreshold,
			},
			want: 1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := karatsubaLen(tt.args.n, tt.args.threshold); got != tt.want {
				t.Errorf("karatsubaLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_leftPad(t *testing.T) {
	type args struct {
		x digits
		n int
	}
	tests := []struct {
		name string
		args args
		want digits
	}{
		{
			name: "left 5 pad",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n: 5,
			},
			want: digits{0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := leftPad(tt.args.x, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("leftPad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rightPad(t *testing.T) {
	type args struct {
		x digits
		n int
	}
	tests := []struct {
		name string
		args args
		want digits
	}{
		{
			name: "right 5 pad",
			args: args{
				x: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
				n: 5,
			},
			want: digits{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rightPad(tt.args.x, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rightPad() = %v, want %v", got, tt.want)
			}
		})
	}
}
