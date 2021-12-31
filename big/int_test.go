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
				t.Errorf("NewInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt_SetString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "pos",
			args: args{s: "123456789"},
			want: NewInt(123456789),
		},
		{
			name: "pos with + sign",
			args: args{s: "+123456789"},
			want: NewInt(123456789),
		},
		{
			name: "neg",
			args: args{s: "-123456789"},
			want: NewInt(-123456789),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Int).SetString(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetString() = %v, want %v", got, tt.want)
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
		{
			name: "1234567890... * 1234567890... with karatsuba",
			args: args{
				x: new(Int).SetString("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"),
				y: new(Int).SetString("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"),
			},
			want: new(Int).SetString("1524157875323883675049535156256668194500838287337600975522511812231126352691000152415888766956267751562263087639079520012193273126047859425087639153757049236500533455762536198787501905199875019052100"),
		},
		{
			name: "9876543210... * 1234567890... with karatsuba",
			args: args{
				x: new(Int).SetString("9876543210987654321098765432109876543210987654321098765432109876543210987654321098765432109876543210"),
				y: new(Int).SetString("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"),
			},
			want: new(Int).SetString("12193263113702179522618503273386678859451150739156363359236761164455788599298790108215200135650052123609205801112635258986434993786160646167367779295611949397448712086533622923332237463801111263526900"),
		},
		{
			name: "9876543210... * 9876543210... with karatsuba",
			args: args{
				x: new(Int).SetString("9876543210987654321098765432109876543210987654321098765432109876543210987654321098765432109876543210"),
				y: new(Int).SetString("9876543210987654321098765432109876543210987654321098765432109876543210987654321098765432109876543210"),
			},
			want: new(Int).SetString("97546105798506325258725803993760097546164761469295351318397422648986531016613331976832801085200426877762536208901082153002591068511507392172275567749340039628145252248135650053345677488187778997104100"),
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
