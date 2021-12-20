package big

// digits は0~9の一桁の数字のスライス
type digits []uint8

// add は |x| + |y| の絶対値による加算を行う
func add(x, y digits) digits {
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
	return norm(abs)
}

// sub は |x| - |y| の絶対値による減算を行う
// 呼び出し側は |x| >= |y| を保証しなければならず、この制約が破られたときpanicする
func sub(x, y digits) digits {
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
	return norm(abs)
}

// mul は |x| * |y| の絶対値による乗算を行う
func mul(x, y digits) digits {
	m, n := len(x), len(y)
	if m == 0 || n == 0 {
		return digits{}
	}
	abs := make(digits, m+n)
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			dx := x[i]
			dy := y[j]
			// 繰り下がりを考慮した各桁の乗算
			abs[i+j+1] += dx * dy
			if abs[i+j+1] >= 10 {
				carry := abs[i+j+1] / 10
				abs[i+j+1] -= carry * 10
				abs[i+j] += carry
			}
		}
	}
	return norm(abs)
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
