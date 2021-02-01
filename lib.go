package dtoa

const (
	kDiySignificandSize = 64
	kDpSignificandSize  = 52
	kDpExponentBias     = 0x3FF + kDpSignificandSize
	kDpMaxExponent      = 0x7FF - kDpExponentBias
	kDpMinExponent      = -kDpExponentBias
	kDpDenormalExponent = -kDpExponentBias + 1
	kDpExponentMask     = 0x7FF0000000000000
	kDpSignificandMask  = 0x000FFFFFFFFFFFFF
	kDpHiddenBit        = 0x0010000000000000

	kSignificandSize  = 52
	kExponentBias     = 0x3FF
	kDenormalExponent = 1 - kExponentBias
	kSignMask         = 0x8000000000000000
	kExponentMask     = 0x7FF0000000000000
	kSignificandMask  = 0x000FFFFFFFFFFFFF
	kHiddenBit        = 0x0010000000000000
)

// 10^-348, 10^-340, ..., 10^340
var kCachedPowers_F = [...]uint64{
	0xfa8fd5a0081c0288, 0xbaaee17fa23ebf76,
	0x8b16fb203055ac76, 0xcf42894a5dce35ea,
	0x9a6bb0aa55653b2d, 0xe61acf033d1a45df,
	0xab70fe17c79ac6ca, 0xff77b1fcbebcdc4f,
	0xbe5691ef416bd60c, 0x8dd01fad907ffc3c,
	0xd3515c2831559a83, 0x9d71ac8fada6c9b5,
	0xea9c227723ee8bcb, 0xaecc49914078536d,
	0x823c12795db6ce57, 0xc21094364dfb5637,
	0x9096ea6f3848984f, 0xd77485cb25823ac7,
	0xa086cfcd97bf97f4, 0xef340a98172aace5,
	0xb23867fb2a35b28e, 0x84c8d4dfd2c63f3b,
	0xc5dd44271ad3cdba, 0x936b9fcebb25c996,
	0xdbac6c247d62a584, 0xa3ab66580d5fdaf6,
	0xf3e2f893dec3f126, 0xb5b5ada8aaff80b8,
	0x87625f056c7c4a8b, 0xc9bcff6034c13053,
	0x964e858c91ba2655, 0xdff9772470297ebd,
	0xa6dfbd9fb8e5b88f, 0xf8a95fcf88747d94,
	0xb94470938fa89bcf, 0x8a08f0f8bf0f156b,
	0xcdb02555653131b6, 0x993fe2c6d07b7fac,
	0xe45c10c42a2b3b06, 0xaa242499697392d3,
	0xfd87b5f28300ca0e, 0xbce5086492111aeb,
	0x8cbccc096f5088cc, 0xd1b71758e219652c,
	0x9c40000000000000, 0xe8d4a51000000000,
	0xad78ebc5ac620000, 0x813f3978f8940984,
	0xc097ce7bc90715b3, 0x8f7e32ce7bea5c70,
	0xd5d238a4abe98068, 0x9f4f2726179a2245,
	0xed63a231d4c4fb27, 0xb0de65388cc8ada8,
	0x83c7088e1aab65db, 0xc45d1df942711d9a,
	0x924d692ca61be758, 0xda01ee641a708dea,
	0xa26da3999aef774a, 0xf209787bb47d6b85,
	0xb454e4a179dd1877, 0x865b86925b9bc5c2,
	0xc83553c5c8965d3d, 0x952ab45cfa97a0b3,
	0xde469fbd99a05fe3, 0xa59bc234db398c25,
	0xf6c69a72a3989f5c, 0xb7dcbf5354e9bece,
	0x88fcf317f22241e2, 0xcc20ce9bd35c78a5,
	0x98165af37b2153df, 0xe2a0b5dc971f303a,
	0xa8d9d1535ce3b396, 0xfb9b7cd9a4a7443c,
	0xbb764c4ca7a44410, 0x8bab8eefb6409c1a,
	0xd01fef10a657842c, 0x9b10a4e5e9913129,
	0xe7109bfba19c0c9d, 0xac2820d9623bf429,
	0x80444b5e7aa7cf85, 0xbf21e44003acdd2d,
	0x8e679c2f5e44ff8f, 0xd433179d9c8cb841,
	0x9e19db92b4e31ba9, 0xeb96bf6ebadf77d9,
	0xaf87023b9bf0ee6b,
}
var kCachedPowers_E = [...]int{
	-1220, -1193, -1166, -1140, -1113, -1087, -1060, -1034, -1007, -980,
	-954, -927, -901, -874, -847, -821, -794, -768, -741, -715,
	-688, -661, -635, -608, -582, -555, -529, -502, -475, -449,
	-422, -396, -369, -343, -316, -289, -263, -236, -210, -183,
	-157, -130, -103, -77, -50, -24, 3, 30, 56, 83,
	109, 136, 162, 189, 216, 242, 269, 295, 322, 348,
	375, 402, 428, 455, 481, 508, 534, 561, 588, 614,
	641, 667, 694, 720, 747, 774, 800, 827, 853, 880,
	907, 933, 960, 986, 1013, 1039, 1066,
}
var kPow10 = [...]uint64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
}

/*
var cDigitsLut = [200]byte{
	'0', '0', '0', '1', '0', '2', '0', '3', '0', '4', '0', '5', '0', '6', '0', '7', '0', '8', '0', '9',
	'1', '0', '1', '1', '1', '2', '1', '3', '1', '4', '1', '5', '1', '6', '1', '7', '1', '8', '1', '9',
	'2', '0', '2', '1', '2', '2', '2', '3', '2', '4', '2', '5', '2', '6', '2', '7', '2', '8', '2', '9',
	'3', '0', '3', '1', '3', '2', '3', '3', '3', '4', '3', '5', '3', '6', '3', '7', '3', '8', '3', '9',
	'4', '0', '4', '1', '4', '2', '4', '3', '4', '4', '4', '5', '4', '6', '4', '7', '4', '8', '4', '9',
	'5', '0', '5', '1', '5', '2', '5', '3', '5', '4', '5', '5', '5', '6', '5', '7', '5', '8', '5', '9',
	'6', '0', '6', '1', '6', '2', '6', '3', '6', '4', '6', '5', '6', '6', '6', '7', '6', '8', '6', '9',
	'7', '0', '7', '1', '7', '2', '7', '3', '7', '4', '7', '5', '7', '6', '7', '7', '7', '8', '7', '9',
	'8', '0', '8', '1', '8', '2', '8', '3', '8', '4', '8', '5', '8', '6', '8', '7', '8', '8', '8', '9',
	'9', '0', '9', '1', '9', '2', '9', '3', '9', '4', '9', '5', '9', '6', '9', '7', '9', '8', '9', '9',
}
*/
const cDigitsLut = "00010203040506070809101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899"

func grow(buf []byte, s int) []byte {
	if s <= 0 {
		return buf
	}
	l := len(buf)
	c := cap(buf)
	if l+s <= c {
		buf = buf[:l+s]
	} else {
		buf = buf[:c]
		if c > l+s-c {
			buf = append(buf, buf[:l+s-c]...)
		} else {
			for i := l + s - c; i > 0; i-- {
				buf = append(buf, 0)
			}
		}
	}
	return buf
}

// http://stackoverflow.com/questions/30614165/is-there-analog-of-memset-in-go
func memsetRepeat(a []byte, v byte) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
