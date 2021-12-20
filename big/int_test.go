package big

import (
	"reflect"
	"testing"
)

func TestNewInt(t *testing.T) {
	type args struct {
		x int64
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "pos",
			args: args{x: 1234567890},
			want: &Int{
				neg: false,
				abs: digits{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			},
		},
		{
			name: "neg",
			args: args{x: -1234567890},
			want: &Int{
				neg: true,
				abs: digits{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntFromInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt_String(t *testing.T) {
	type fields struct {
		neg bool
		abs digits
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "pos",
			fields: fields{
				neg: false,
				abs: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: "123456789",
		},
		{
			name: "neg",
			fields: fields{
				neg: true,
				abs: digits{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: "-123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Int{
				neg: tt.fields.neg,
				abs: tt.fields.abs,
			}
			if got := b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		x *Int
		y *Int
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "x + y",
			args: args{
				x: NewInt(11111),
				y: NewInt(1111),
			},
			want: NewInt(12222),
		},
		{
			name: "x + y (with carry)",
			args: args{
				x: NewInt(123456789),
				y: NewInt(123456789),
			},
			want: NewInt(246913578),
		},
		{
			name: "x + y (with top digit carry)",
			args: args{
				x: NewInt(99999),
				y: NewInt(99999),
			},
			want: NewInt(199998),
		},
		{
			name: "x + y (with repeatable carry)",
			args: args{
				x: NewInt(99999),
				y: NewInt(1),
			},
			want: NewInt(100000),
		},
		{
			name: "(-x) + (-y) (with carry)",
			args: args{
				x: NewInt(-123456789),
				y: NewInt(-123456789),
			},
			want: NewInt(-246913578),
		},
		{
			name: "-x + y = 0",
			args: args{
				x: NewInt(-123456789),
				y: NewInt(123456789),
			},
			want: Zero,
		},
		{
			name: "x + (-y) (|x| >= |y|)",
			args: args{
				x: NewInt(9384),
				y: NewInt(-3894),
			},
			want: NewInt(5490),
		},
		{
			name: "(-x) + y (|x| >= |y|)",
			args: args{
				x: NewInt(-9384),
				y: NewInt(3894),
			},
			want: NewInt(-5490),
		},
		{
			name: "x + (-y) (|x| < |y|)",
			args: args{
				x: NewInt(3894),
				y: NewInt(-9384),
			},
			want: NewInt(-5490),
		},
		{
			name: "(-x) + y (|x| < |y|)",
			args: args{
				x: NewInt(-3894),
				y: NewInt(9384),
			},
			want: NewInt(5490),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	type args struct {
		x *Int
		y *Int
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "x - y (|x| >= |y|)",
			args: args{
				x: NewInt(11111),
				y: NewInt(1111),
			},
			want: NewInt(10000),
		},
		{
			name: "x - y (|x| < |y|)",
			args: args{
				x: NewInt(123456789),
				y: NewInt(234567890),
			},
			want: NewInt(-111111101),
		},
		{
			name: "x - y (with borrow)",
			args: args{
				x: NewInt(123456789),
				y: NewInt(99999999),
			},
			want: NewInt(23456790),
		},
		{
			name: "(-x) - (-y) (with borrow)",
			args: args{
				x: NewInt(-123456789),
				y: NewInt(-99999999),
			},
			want: NewInt(-23456790),
		},
		{
			name: "(-x) - y",
			args: args{
				x: NewInt(-123456789),
				y: NewInt(234567890),
			},
			want: NewInt(-358024679),
		},
		{
			name: "x - (-y)",
			args: args{
				x: NewInt(123456789),
				y: NewInt(-234567890),
			},
			want: NewInt(358024679),
		},
		{
			name: "x - y = 0",
			args: args{
				x: NewInt(123456789),
				y: NewInt(123456789),
			},
			want: Zero,
		},
		{
			name: "x - y (with repeated borrow)",
			args: args{
				x: NewInt(100000),
				y: NewInt(1),
			},
			want: NewInt(99999),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sub(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMul(t *testing.T) {
	type args struct {
		x *Int
		y *Int
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "x * y",
			args: args{
				x: NewInt(987654321),
				y: NewInt(987654321),
			},
			want: NewInt(975461057789971041),
		},
		{
			name: "x * y (different digit length)",
			args: args{
				x: NewInt(987654321),
				y: NewInt(87654321),
			},
			want: NewInt(86572168889971041),
		},
		{
			name: "-x * -y",
			args: args{
				x: NewInt(-987654321),
				y: NewInt(-987654321),
			},
			want: NewInt(975461057789971041),
		},
		{
			name: "-x * y",
			args: args{
				x: NewInt(-987654321),
				y: NewInt(987654321),
			},
			want: NewInt(-975461057789971041),
		},
		{
			name: "x * -y",
			args: args{
				x: NewInt(987654321),
				y: NewInt(-987654321),
			},
			want: NewInt(-975461057789971041),
		},
		{
			name: "0 * y",
			args: args{
				x: Zero,
				y: NewInt(-987654321),
			},
			want: Zero,
		},
		{
			name: "x * 0",
			args: args{
				x: NewInt(987654321),
				y: Zero,
			},
			want: Zero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mul(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		x *Int
		y *Int
	}
	tests := []struct {
		name    string
		args    args
		wantQuo *Int
		wantRem *Int
	}{
		{
			name: "x / y (no rem)",
			args: args{
				x: NewInt(1735745558982),
				y: NewInt(4984423),
			},
			wantQuo: NewInt(348234),
			wantRem: Zero,
		},
		{
			name: "x / y (with rem)",
			args: args{
				x: NewInt(1735745558983),
				y: NewInt(4984423),
			},
			wantQuo: NewInt(348234),
			wantRem: NewInt(1),
		},
		{
			name: "-x / -y (with rem)",
			args: args{
				x: NewInt(-1735745558983),
				y: NewInt(-4984423),
			},
			wantQuo: NewInt(348234),
			wantRem: NewInt(-1),
		},
		{
			name: "-x / y (with rem)",
			args: args{
				x: NewInt(-1735745558983),
				y: NewInt(4984423),
			},
			wantQuo: NewInt(-348234),
			wantRem: NewInt(-1),
		},
		{
			name: "x / -y (with rem)",
			args: args{
				x: NewInt(1735745558983),
				y: NewInt(-4984423),
			},
			wantQuo: NewInt(-348234),
			wantRem: NewInt(1),
		},
		{
			name: "x < y",
			args: args{
				x: NewInt(17357),
				y: NewInt(4984423),
			},
			wantQuo: Zero,
			wantRem: NewInt(17357),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuo, gotRem := Div(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(gotQuo, tt.wantQuo) {
				t.Errorf("Div() gotQuo = %v, want %v", gotQuo, tt.wantQuo)
			}
			if !reflect.DeepEqual(gotRem, tt.wantRem) {
				t.Errorf("Div() gotRem = %v, want %v", gotRem, tt.wantRem)
			}
		})
	}
}
