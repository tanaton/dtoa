package dtoa

func U32toa(buf []byte, value uint32) []byte {
	d := cDigitsLut
	if value < 10000 {
		var tmp [4]byte
		d1 := (value / 100) << 1
		d2 := (value % 100) << 1

		start := 3
		switch {
		case value >= 1000:
			tmp[0] = d[d1]
			start--
			fallthrough
		case value >= 100:
			tmp[1] = d[d1+1]
			start--
			fallthrough
		case value >= 10:
			tmp[2] = d[d2]
			start--
		}
		tmp[3] = d[d2+1]
		buf = append(buf, tmp[start:]...)
	} else if value < 100000000 {
		// value = bbbbcccc
		var tmp [8]byte
		b := value / 10000
		c := value % 10000

		d1 := (b / 100) << 1
		d2 := (b % 100) << 1

		d3 := (c / 100) << 1
		d4 := (c % 100) << 1

		start := 3
		switch {
		case value >= 10000000:
			tmp[0] = d[d1]
			start--
			fallthrough
		case value >= 1000000:
			tmp[1] = d[d1+1]
			start--
			fallthrough
		case value >= 100000:
			tmp[2] = d[d2]
			start--
		}
		tmp[3] = d[d2+1]
		tmp[4] = d[d3]
		tmp[5] = d[d3+1]
		tmp[6] = d[d4]
		tmp[7] = d[d4+1]
		buf = append(buf, tmp[start:]...)
	} else {
		// value = aabbbbcccc in decimal
		var tmp [10]byte
		a := value / 100000000 // 1 to 42
		value %= 100000000

		b := value / 10000 // 0 to 9999
		c := value % 10000 // 0 to 9999

		d1 := (b / 100) << 1
		d2 := (b % 100) << 1

		d3 := (c / 100) << 1
		d4 := (c % 100) << 1

		var start int
		if a >= 10 {
			start = 0
			i := a << 1
			tmp[0] = d[i]
			tmp[1] = d[i+1]
		} else {
			start = 1
			tmp[1] = '0' + byte(a)
		}
		tmp[2] = d[d1]
		tmp[3] = d[d1+1]
		tmp[4] = d[d2]
		tmp[5] = d[d2+1]
		tmp[6] = d[d3]
		tmp[7] = d[d3+1]
		tmp[8] = d[d4]
		tmp[9] = d[d4+1]
		buf = append(buf, tmp[start:]...)
	}
	return buf
}

func I32toa(buf []byte, value int32) []byte {
	u := uint32(value)
	if value < 0 {
		buf = append(buf, '-')
		u = -u
	}
	return U32toa(buf, u)
}

func U64toa(buf []byte, value uint64) []byte {
	const kTen8 uint64 = 100000000
	const kTen9 uint64 = kTen8 * 10
	const kTen10 uint64 = kTen8 * 100
	const kTen11 uint64 = kTen8 * 1000
	const kTen12 uint64 = kTen8 * 10000
	const kTen13 uint64 = kTen8 * 100000
	const kTen14 uint64 = kTen8 * 1000000
	const kTen15 uint64 = kTen8 * 10000000
	const kTen16 uint64 = kTen8 * kTen8
	d := cDigitsLut

	if value < kTen8 {
		v := uint32(value)
		if v < 10000 {
			var tmp [4]byte
			d1 := (v / 100) << 1
			d2 := (v % 100) << 1

			start := 3
			switch {
			case v >= 1000:
				tmp[0] = d[d1]
				start--
				fallthrough
			case v >= 100:
				tmp[1] = d[d1+1]
				start--
				fallthrough
			case v >= 10:
				tmp[2] = d[d2]
				start--
			}
			tmp[3] = d[d2+1]
			buf = append(buf, tmp[start:]...)
		} else {
			// value = bbbbcccc
			var tmp [8]byte
			b := v / 10000
			c := v % 10000

			d1 := (b / 100) << 1
			d2 := (b % 100) << 1

			d3 := (c / 100) << 1
			d4 := (c % 100) << 1

			start := 3
			switch {
			case value >= 10000000:
				tmp[0] = d[d1]
				start--
				fallthrough
			case value >= 1000000:
				tmp[1] = d[d1+1]
				start--
				fallthrough
			case value >= 100000:
				tmp[2] = d[d2]
				start--
			}
			tmp[3] = d[d2+1]
			tmp[4] = d[d3]
			tmp[5] = d[d3+1]
			tmp[6] = d[d4]
			tmp[7] = d[d4+1]
			buf = append(buf, tmp[start:]...)
		}
	} else if value < kTen16 {
		var tmp [16]byte
		v0 := uint32(value / kTen8)
		v1 := uint32(value % kTen8)

		b0 := v0 / 10000
		c0 := v0 % 10000

		d1 := (b0 / 100) << 1
		d2 := (b0 % 100) << 1

		d3 := (c0 / 100) << 1
		d4 := (c0 % 100) << 1

		b1 := v1 / 10000
		c1 := v1 % 10000

		d5 := (b1 / 100) << 1
		d6 := (b1 % 100) << 1

		d7 := (c1 / 100) << 1
		d8 := (c1 % 100) << 1

		start := 8
		switch {
		case value >= kTen15:
			tmp[0] = d[d1]
			start--
			fallthrough
		case value >= kTen14:
			tmp[1] = d[d1+1]
			start--
			fallthrough
		case value >= kTen13:
			tmp[2] = d[d2]
			start--
			fallthrough
		case value >= kTen12:
			tmp[3] = d[d2+1]
			start--
			fallthrough
		case value >= kTen11:
			tmp[4] = d[d3]
			start--
			fallthrough
		case value >= kTen10:
			tmp[5] = d[d3+1]
			start--
			fallthrough
		case value >= kTen9:
			tmp[6] = d[d4]
			start--
			fallthrough
		case value >= kTen8:
			tmp[7] = d[d4+1]
			start--
		}
		tmp[8] = d[d5]
		tmp[9] = d[d5+1]
		tmp[10] = d[d6]
		tmp[11] = d[d6+1]
		tmp[12] = d[d7]
		tmp[13] = d[d7+1]
		tmp[14] = d[d8]
		tmp[15] = d[d8+1]
		buf = append(buf, tmp[start:]...)
	} else {
		var tmp [20]byte
		a := uint32(value / kTen16) // 1 to 1844
		value %= kTen16

		v0 := uint32(value / kTen8)
		v1 := uint32(value % kTen8)

		b0 := v0 / 10000
		c0 := v0 % 10000

		d1 := (b0 / 100) << 1
		d2 := (b0 % 100) << 1

		d3 := (c0 / 100) << 1
		d4 := (c0 % 100) << 1

		b1 := v1 / 10000
		c1 := v1 % 10000

		d5 := (b1 / 100) << 1
		d6 := (b1 % 100) << 1

		d7 := (c1 / 100) << 1
		d8 := (c1 % 100) << 1

		var start int
		if a < 10 {
			start = 3
			tmp[3] = '0' + byte(a)
		} else if a < 100 {
			start = 2
			i := a << 1
			tmp[2] = d[i]
			tmp[3] = d[i+1]
		} else if a < 1000 {
			start = 1
			tmp[1] = '0' + byte(a/100)
			i := (a % 100) << 1
			tmp[2] = d[i]
			tmp[3] = d[i+1]
		} else {
			start = 0
			i := (a / 100) << 1
			j := (a % 100) << 1
			tmp[0] = d[i]
			tmp[1] = d[i+1]
			tmp[2] = d[j]
			tmp[3] = d[j+1]
		}
		tmp[4] = d[d1]
		tmp[5] = d[d1+1]
		tmp[6] = d[d2]
		tmp[7] = d[d2+1]
		tmp[8] = d[d3]
		tmp[9] = d[d3+1]
		tmp[10] = d[d4]
		tmp[11] = d[d4+1]
		tmp[12] = d[d5]
		tmp[13] = d[d5+1]
		tmp[14] = d[d6]
		tmp[15] = d[d6+1]
		tmp[16] = d[d7]
		tmp[17] = d[d7+1]
		tmp[18] = d[d8]
		tmp[19] = d[d8+1]
		buf = append(buf, tmp[start:]...)
	}

	return buf
}

func I64toa(buf []byte, value int64) []byte {
	u := uint64(value)
	if value < 0 {
		buf = append(buf, '-')
		u = -u
	}
	return U64toa(buf, u)
}
