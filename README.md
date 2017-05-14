# dtoa
本ライブラリは浮動小数点数値をバイト列に変換する高速関数を提供します。

実装は[Milo Yip](https://github.com/miloyip)氏のC++実装[dtoa.h](https://github.com/miloyip/rapidjson/blob/master/include/rapidjson/internal/dtoa.h)をGoに移植したものとなっています。

おまけでitoaの実装も含んでいます。

## ベンチマーク
|数値|小数点以下桁数|strconv.AppendFloat|dtoa.Dtoa|速度差|
|:-:|:-:|:-:|:-:|:-:|
|1.1|全部|434 ns/op|173 ns/op|2.5倍|
|3.1415926535|全部|316 ns/op|272 ns/op|1.16倍|
|2.225073858507201e-308|1桁|15694 ns/op|325 ns/op|48.3倍|

## 使い方
```go:main.go
// strconv.FormatFloat()風
fmt.Println(string(dtoa.Dtoa(nil, 1.12345678, 5)))
// 1.12345

// strconv.AppendFloat()風
buf := []byte{}
buf = dtoa.Dtoa(buf, 0.1, 5)
buf = append(buf, ',')
buf = dtoa.Dtoa(buf, 0.2, 5)
fmt.Println(string(buf))
// 0.1,0.2
```

## License
MIT
