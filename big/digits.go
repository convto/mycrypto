package big

// digits は0~9の一桁の数字のスライス
type digits []uint8

// add は |x| + |y| の絶対値による加算を行う
func add(x, y digits) digits {
	return norm(basicAdd(x, y))
}

// basicAdd は一桁ずつ加算し結果を返す
// 結果は大きい方の桁+1の長さで返し、上位の桁に0を含む可能性がある
func basicAdd(x, y digits) digits {
	m, n := len(x), len(y)
	var l int
	if m > n {
		l = m
	} else {
		l = n
	}
	abs := make(digits, l+1) // capは繰り上がり考慮で+1
	for i := l - 1; i >= 0; i-- {
		var dx, dy uint8
		if (m-1)-(l-1-i) >= 0 {
			dx = x[(m-1)-(l-1-i)]
		}
		if (n-1)-(l-1-i) >= 0 {
			dy = y[(n-1)-(l-1-i)]
		}
		// 繰り上がりを考慮した各桁の加算
		abs[i+1] = abs[i+1] + dx + dy
		if abs[i+1] >= 10 {
			// 結果が10以上だったら次の桁に繰り上がり
			abs[i+1] -= 10
			abs[i] = 1
		}
	}
	return abs
}

// sub は |x| - |y| の絶対値による減算を行う
// 呼び出し側は |x| >= |y| を保証しなければならず、この制約が破られたときpanicする
func sub(x, y digits) digits {
	return norm(basicSub(x, y))
}

// basicSub は一桁ずつ減算し結果を返す
// 結果は大きい方の桁数とおなじ桁数で返し、上位の桁に0を含む可能性がある
func basicSub(x, y digits) digits {
	m, n := len(x), len(y)
	if m < n || cmp(x, y) == -1 {
		panic("underflow")
	}
	var l int
	if m > n {
		l = m
	} else {
		l = n
	}
	abs := make(digits, l)
	for i := l - 1; i >= 0; i-- {
		var dx, dy uint8
		if (m-1)-(l-1-i) >= 0 {
			dx = x[(m-1)-(l-1-i)]
		}
		if (n-1)-(l-1-i) >= 0 {
			dy = y[(n-1)-(l-1-i)]
		}

		// 繰り下がりを考慮した各桁の減算
		if dx < abs[i] || dx-abs[i] < dy {
			abs[i] = (dx - abs[i] + 10) - dy
			// 次の桁に繰り下がりのflagを立てておく
			// |x| >= |y| が保証されているので最上位桁の繰り下がりは考慮しなくてよい
			abs[i-1] = 1
		} else {
			abs[i] = (dx - abs[i]) - dy
		}
	}
	return abs
}

// mul は |x| * |y| の絶対値による乗算を行う
func mul(x, y digits) digits {
	m, n := len(x), len(y)
	switch {
	case m < n:
		return norm(mul(y, x))
	case m == 0 || n == 0:
		return digits{}
	}

	if m < karatsubaThreshold && n < karatsubaThreshold {
		return norm(basicMul(x, y))
	}

	// karatsubaThreshold までが2のべき乗となるようにpaddingをとる
	k := karatsubaLen(m, karatsubaThreshold)
	px := leftPad(x, k-len(x))
	py := leftPad(y, k-len(y))
	return norm(karatsuba(px, py))
}

// basicMul は long multiplication で乗算を行う
func basicMul(x, y digits) digits {
	m, n := len(x), len(y)
	abs := make(digits, m+n)
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			dx := x[i]
			dy := y[j]
			// 積が0以外のときのみ計算
			if dx != 0 && dy != 0 {
				abs[i+j+1] += dx * dy
			}
			// 繰り上がり処理
			if abs[i+j+1] >= 10 {
				carry := abs[i+j+1] / 10
				abs[i+j+1] -= carry * 10
				abs[i+j] += carry
			}
		}
	}
	return abs
}

// karatsubaLen はnが閾値まで2べきならnをそのまま返し、そうでなければ閾値まで2べきとなるようなn以上のできるだけ小さい数を返す
func karatsubaLen(n, threshold int) int {
	var i uint = 0
	for n > threshold {
		if n&1 == 1 {
			n++
		}
		n >>= 1
		i++
	}
	n <<= i
	return n
}

// karatsuba法は定数倍が大きいので、40桁以上の乗算について適用させるようにする
const karatsubaThreshold = 40

// karatsuba は karatsuba's algorithm で乗算を行う
// 呼び出し側は len(x) == len(y) かつ karatsubaThreshold まで2のべき乗のサイズになっていることを保証すること
// TODO: 再帰するたびに内部でメモリallocするのでうまく結果を共有してalloc回数を減らしたい
func karatsuba(x, y digits) digits {
	m := len(x)
	// len(x) について、奇数/閾値以下/0のいずれかなら通常の乗算にて計算する
	if m&1 != 0 || m <= karatsubaThreshold || m < 2 {
		return basicMul(x, y)
	}
	m2 := m >> 1
	x1, x0 := x[m2:], x[0:m2]
	y1, y0 := y[m2:], y[0:m2]

	x0y0 := karatsuba(x0, y0)
	x1y1 := karatsuba(x1, y1)

	// x0y1 + x1y0 = (x0-x1)(y1-y0) + (x0y0 + x1y1) となるのでその計算
	// (x0+x1)(y1+y0) - (x0y0 + x1y1) の形にも整理できるが、
	// 加算はcarryが発生する可能性があり後続の再帰処理にて2のべき乗のサイズとならない可能性があるため減算の形で扱っている
	s := 1
	xd := make(digits, m2)
	if cmp(x0, x1) >= 0 {
		xd = basicSub(x0, x1)
	} else {
		s = -s
		xd = basicSub(x1, x0)
	}
	yd := make(digits, m2)
	if cmp(y1, y0) >= 0 {
		yd = basicSub(y1, y0)
	} else {
		s = -s
		yd = basicSub(y0, y1)
	}
	var p digits
	if s < 0 {
		p = basicSub(basicAdd(x0y0, x1y1), karatsuba(xd, yd))
	} else {
		p = basicAdd(karatsuba(xd, yd), basicAdd(x0y0, x1y1))
	}

	// x1y1*(10^m) + p*(10^m2) + x0y0
	x0y0 = rightPad(x0y0, m)
	p = rightPad(p, m2)
	return basicAdd(basicAdd(x0y0, x1y1), p)
}

// div は |x| / |y| の絶対値による除算を行い、商を quo あまりを rem で返す
// 呼び出し側は y != 0 を保証しなければならず、この条件が守られないときpanicする
func div(x, y digits) (quo digits, rem digits) {
	m, n := len(x), len(y)
	if n == 0 {
		panic("division by zero")
	}
	if cmp(x, y) < 0 {
		return digits{}, x
	}
	// 結果の桁数lを求める
	l := m - n // 少なくとも x > y*10^(m-n-1) で、そのとき商の桁数は m-n
	if cmp(x[:n-1], y) >= 0 {
		l += 1 // y*10^(m-n) したときにxのほうが大きければさらにその桁についても商が存在する
	}

	// 内部の乗算のコストをかるくするため、remは除算中の桁に関連するあまりだけ保持する
	// y*10^l のときの10の指数ぶんの桁数はそのステップでは値が変動しないので初期値から除いておく
	remPos := (m) - (l - 1)
	rem = x[:remPos]
	quo = make(digits, l)
	for i := 0; i < l; i++ {
		for j := 1; j < 10; j++ {
			// 現在のあまりを超えない範囲で最大の y * j (jは1-9) を求めて今の桁の商とする
			if cmp(mul(y, digits{uint8(j)}), rem) <= 0 {
				quo[i] = uint8(j)
			} else {
				break
			}
		}
		rem = sub(rem, mul(y, digits{quo[i]}))
		if i < l-1 {
			remPos += 1
			rem = append(rem, x[remPos-1])
		}
	}
	return quo, rem
}

// norm は上の桁から連続して0になる部分をtrimする
func norm(abs digits) digits {
	i := 0
	for i < len(abs) && abs[i] == 0 {
		i++
	}
	return abs[i:]
}

// cmp はx, yを比較して以下の結果を返す
// x > y  -> 1
// x == y -> 0
// x < y  -> -1
func cmp(x, y digits) int8 {
	// サイズが異なる場合サイズの大きいほうが絶対値が大きい
	m, n := len(x), len(y)
	if m != n {
		if m > n {
			return 1
		} else {
			return -1
		}
	}

	// サイズが同一の場合はより上位の桁の値が大きいほうが絶対値が大きい
	for i := 0; i < m; i++ {
		dx, dy := x[i], y[i]
		switch {
		case dx > dy:
			return 1
		case dx < dy:
			return -1
		}
	}
	return 0
}

// leftPad は指定された数だけ上位の桁に0を追加します
func leftPad(x digits, n int) digits {
	m := len(x)
	l := m + n
	abs := make(digits, l)
	for i := 0; i < l; i++ {
		if i < n {
			abs[i] = 0
		} else {
			abs[i] = x[i-n]
		}
	}
	return abs
}

// rightPad は指定された数だけ下位の桁に0を追加します
func rightPad(x digits, n int) digits {
	m := len(x)
	l := m + n
	abs := make(digits, l)
	for i := 0; i < l; i++ {
		if i < m {
			abs[i] = x[i]
		} else {
			abs[i] = 0
		}
	}
	return abs
}
