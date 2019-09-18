package dtoa

func grisuRound(buf []byte, l int, delta, rest, ten_kappa, wp_w uint64) {
	for rest < wp_w && delta-rest >= ten_kappa && (rest+ten_kappa < wp_w || wp_w-rest > rest+ten_kappa-wp_w) {
		buf[l-1]--
		rest += ten_kappa
	}
}

func countDecimalDigit32(n uint32) int {
	// Simple pure C++ implementation was faster than __builtin_clz version in this situation.
	switch {
	case n < 10:
		return 1
	case n < 100:
		return 2
	case n < 1000:
		return 3
	case n < 10000:
		return 4
	case n < 100000:
		return 5
	case n < 1000000:
		return 6
	case n < 10000000:
		return 7
	case n < 100000000:
		return 8
	}
	return 9
}

func digitGen(buf []byte, W, Mp DiyFp, delta uint64, K int) ([]byte, int, int) {
	base := len(buf)
	one := DiyFp{uint64(1) << uint64(-Mp.e), Mp.e}
	wp_w := Mp.Minus(W)
	p1 := uint32(Mp.f >> uint64(-one.e))
	var p2 uint64 = Mp.f & (one.f - 1)
	kappa := countDecimalDigit32(p1)
	l := 0

	for kappa > 0 {
		var d uint32
		switch kappa {
		case 9:
			d = p1 / 100000000
			p1 %= 100000000
		case 8:
			d = p1 / 10000000
			p1 %= 10000000
		case 7:
			d = p1 / 1000000
			p1 %= 1000000
		case 6:
			d = p1 / 100000
			p1 %= 100000
		case 5:
			d = p1 / 10000
			p1 %= 10000
		case 4:
			d = p1 / 1000
			p1 %= 1000
		case 3:
			d = p1 / 100
			p1 %= 100
		case 2:
			d = p1 / 10
			p1 %= 10
		case 1:
			d = p1
			p1 = 0
		default:
			d = 0
		}
		if d != 0 || l != 0 {
			buf = append(buf, '0'+byte(d))
			l++
		}
		kappa--
		var tmp uint64 = (uint64(p1) << uint64(-one.e)) + p2
		if tmp <= delta {
			K += kappa
			grisuRound(buf[base:], l, delta, tmp, kPow10[kappa]<<uint64(-one.e), wp_w.f)
			return buf, l, K
		}
	}

	// kappa = 0
	for {
		p2 *= 10
		delta *= 10
		d := byte(p2 >> uint64(-one.e))
		if d != 0 || l != 0 {
			buf = append(buf, '0'+d)
			l++
		}
		p2 &= one.f - 1
		kappa--
		if p2 < delta {
			K += kappa
			grisuRound(buf[base:], l, delta, p2, one.f, wp_w.f*kPow10[-kappa])
			return buf, l, K
		}
	}
	return buf, l, K
}

func Grisu2(buf []byte, value float64) ([]byte, int, int) {
	v := DiyFpDouble(value)
	w_m, w_p := v.NormalizedBoundaries()

	c_mk, K := GetCachedPower(w_p.e)
	W := v.Normalize().Multiplication(c_mk)
	Wp := w_p.Multiplication(c_mk)
	Wm := w_m.Multiplication(c_mk)
	Wm.f++
	Wp.f--
	return digitGen(buf, W, Wp, Wp.f-Wm.f, K)
}

func writeExponent(buf []byte, K int) []byte {
	if K < 0 {
		buf = append(buf, '-')
		K = -K
	}

	if K >= 100 {
		buf = append(buf, '0'+byte(K/100))
		K %= 100
		d := K * 2
		buf = append(buf, cDigitsLut[d], cDigitsLut[d+1])
	} else if K >= 10 {
		d := K * 2
		buf = append(buf, cDigitsLut[d], cDigitsLut[d+1])
	} else {
		buf = append(buf, '0'+byte(K))
	}
	return buf
}

func Prettify(buf []byte, l, k, maxDecimalPlaces int) []byte {
	base := len(buf) - l
	kk := l + k // 10^(kk-1) <= v < 10^kk

	if 0 <= k && kk <= 21 {
		// 1234e7 -> 12340000000
		buf = grow(buf, k+2)
		memsetRepeat(buf[base+l:base+kk], '0')
		buf[base+kk] = '.'
		buf[base+kk+1] = '0'
		return buf[:base+kk+2]
	} else if 0 < kk && kk <= 21 {
		// 1234e-2 -> 12.34
		buf = grow(buf, 1)
		// std::memmove(&buf[kk + 1], &buf[kk], static_cast<size_t>(l - kk));
		copy(buf[base+kk+1:], buf[base+kk:base+l])
		buf[base+kk] = '.'
		if 0 > k+maxDecimalPlaces {
			// When maxDecimalPlaces = 2, 1.2345 -> 1.23, 1.102 -> 1.1
			// Remove extra trailing zeros (at least one) after truncation.
			for i := kk + maxDecimalPlaces; i > kk+1; i-- {
				if buf[base+i] != '0' {
					return buf[:base+i+1]
				}
			}
			return buf[:base+kk+2] // Reserve one zero
		} else {
			return buf[:base+l+1]
		}
	} else if -6 < kk && kk <= 0 {
		// 1234e-6 -> 0.001234
		offset := 2 - kk
		buf = grow(buf, offset)
		// std::memmove(&buf[offset], &buf[0], static_cast<size_t>(l));
		copy(buf[base+offset:], buf[base:base+l])
		buf[base] = '0'
		buf[base+1] = '.'
		memsetRepeat(buf[base+2:base+offset], '0')
		if l-kk > maxDecimalPlaces {
			// When maxDecimalPlaces = 2, 0.123 -> 0.12, 0.102 -> 0.1
			// Remove extra trailing zeros (at least one) after truncation.
			for i := maxDecimalPlaces + 1; i > 2; i-- {
				if buf[base+i] != '0' {
					return buf[:base+i+1]
				}
			}
			return buf[:base+3] // Reserve one zero
		}
		return buf[:base+l+offset]
	} else if kk < -maxDecimalPlaces {
		// Truncate to zero
		buf = grow(buf, 3-(len(buf)-base))
		buf[base] = '0'
		buf[base+1] = '.'
		buf[base+2] = '0'
		return buf[:base+3]
	} else if l == 1 {
		// 1e30
		buf = append(buf, 'e')
		return writeExponent(buf, kk-1)
	} else {
		// 1234e30 -> 1.234e33
		buf = grow(buf, 2)
		// std::memmove(&buf[2], &buf[1], static_cast<size_t>(l - 1));
		copy(buf[base+2:], buf[base+1:base+1+l-1])
		buf[base+1] = '.'
		buf[base+l+1] = 'e'
		return writeExponent(buf, kk-1)
	}
}

func PrettifySimple(buf []byte, l, k, maxDecimalPlaces int) []byte {
	buf = Prettify(buf, l, k, maxDecimalPlaces)
	s := len(buf)
	if s > 2 && buf[s-1] == '0' && buf[s-2] == '.' {
		return buf[:s-2]
	}
	return buf
}

func Dtoa(buf []byte, value float64, maxDecimalPlaces int) []byte {
	d := NewDouble(value)
	if d.IsZero() {
		if d.Sign() {
			buf = append(buf, '-')
		}
		return append(buf, '0', '.', '0')
	} else {
		if value < 0 {
			buf = append(buf, '-')
			value = -value
		}
		if maxDecimalPlaces < 0 {
			maxDecimalPlaces = 324
		}
		var l, K int
		buf, l, K = Grisu2(buf, value)
		return Prettify(buf, l, K, maxDecimalPlaces)
	}
}

func DtoaSimple(buf []byte, value float64, maxDecimalPlaces int) []byte {
	d := NewDouble(value)
	if d.IsZero() {
		if d.Sign() {
			buf = append(buf, '-')
		}
		return append(buf, '0')
	} else {
		if value < 0 {
			buf = append(buf, '-')
			value = -value
		}
		if maxDecimalPlaces < 0 {
			maxDecimalPlaces = 324
		}
		var l, K int
		buf, l, K = Grisu2(buf, value)
		return PrettifySimple(buf, l, K, maxDecimalPlaces)
	}
}
