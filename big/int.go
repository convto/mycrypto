package big

import (
	"fmt"
)

// Int は巨大な整数演算をする型です
type Int struct {
	neg bool
	abs digits
}

// neg = false, abs = [] を0とする
var Zero = &Int{neg: false, abs: digits{}}

func NewInt(x int64) *Int {
	neg := false
	if x < 0 {
		neg = true
		x = -x
	}
	abs := make(digits, 0)
	for i := x; i >= 1; {
		d := uint8(i % 10)
		if len(abs) == 0 {
			abs = append(abs, d)
		} else {
			abs = append([]uint8{d}, abs...)
		}
		i /= 10
	}
	return &Int{
		neg: neg,
		abs: abs,
	}
}

func (b *Int) String() string {
	var s string
	if b.neg {
		s += "-"
	}
	for _, v := range b.abs {
		s += fmt.Sprintf("%d", v)
	}
	return s
}

// Add は整数の和を求める
func Add(x, y *Int) *Int {
	neg := x.neg
	var abs digits
	if x.neg == y.neg {
		abs = add(x.abs, y.abs)
	} else {
		// xとyの正負が異なっていれば絶対値についての減算と言い換えることができる
		if cmp(x.abs, y.abs) >= 0 {
			abs = sub(x.abs, y.abs)
		} else {
			// (-x) + y == y - x となるので結果のsignを反転させれば成り立つ
			abs = sub(y.abs, x.abs)
			neg = !neg
		}
	}
	return &Int{
		neg: len(abs) > 0 && neg,
		abs: abs,
	}
}

// Sub は整数の差を求める
func Sub(x, y *Int) *Int {
	neg := x.neg
	var abs digits
	if x.neg != y.neg {
		// xとyの正負が異なれば絶対値についての加算と言い換えることができる
		abs = add(x.abs, y.abs)
	} else {
		if cmp(x.abs, y.abs) >= 0 {
			abs = sub(x.abs, y.abs)
		} else {
			// (-x) - (-y) == y - x となるので結果のsignを反転させれば成り立つ
			abs = sub(y.abs, x.abs)
			neg = !neg
		}
	}
	return &Int{
		neg: len(abs) > 0 && neg,
		abs: abs,
	}
}

// Mul は整数の積を求める
func Mul(x, y *Int) *Int {
	abs := mul(x.abs, y.abs)
	return &Int{
		neg: len(abs) > 0 && x.neg != y.neg,
		abs: abs,
	}
}

// Div は整数の商とあまりを求める
func Div(x, y *Int) (quo *Int, rem *Int) {
	q, r := div(x.abs, y.abs)
	quo = &Int{
		neg: len(q) > 0 && x.neg != y.neg,
		abs: q,
	}
	rem = &Int{
		neg: len(r) > 0 && x.neg,
		abs: r,
	}
	return quo, rem
}

// Cmp はx, yが等しいかどうかを判定します
func Cmp(x, y *Int) int8 {
	if x.neg != y.neg {
		if x.neg {
			return -1
		}
		return 1
	}
	return cmp(x.abs, y.abs)
}
