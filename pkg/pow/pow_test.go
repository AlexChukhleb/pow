package pow

import (
	"crypto"
	"encoding/hex"
	"testing"
)

func TestSearchSHA1(t *testing.T) {
	arrShort, _ := hex.DecodeString("5348412d310010")
	_, err := SearchKey(arrShort, 0x3FFFFFFF)
	if err == nil {
		t.Fatal("skip err")
		return
	}

	//arr := Gen(crypto.SHA1, 16)
	arr, _ := hex.DecodeString("5348412d3100100029e5466300557d45b38fcd70d2cabe506c07a16c58a6cb7be7a921552302ce79a45f0de20d43e83cd219e7ab7cf16d2d6841f5a5472b7bfcb139ae9565006e01ebe6b6d53d09dab87af7ea15fe96c1c781d27662c05d5922ac156ac1d64b2e71f3c9e9327f9b8b9034f95a214b775209dcb15242a6")

	//log.Info(hex.EncodeToString(arr))

	key, err := SearchKey(arr, 0x3FFFFFFF)
	if err != nil {
		t.Fatal(err)
		return
	}

	ok, err := Check(arr, key, 0x3FFFFFFF)
	if err != nil {
		t.Fatal(err)
		return
	}

	if !ok {
		t.Fatal("!ok")
	}

	ok, err = Check(arr, key, 1)
	if err == nil {
		t.Fatal("skip err")
		return
	}

	if ok {
		t.Fatal("ok")
	}
}

func BenchmarkCheckSHA256_1(b *testing.B)  { benchmarkCheck(crypto.SHA256, 1, b) }
func BenchmarkCheckSHA256_20(b *testing.B) { benchmarkCheck(crypto.SHA256, 20, b) }
func BenchmarkCheckSHA1_1(b *testing.B)    { benchmarkCheck(crypto.SHA1, 1, b) }
func BenchmarkCheckSHA1_20(b *testing.B)   { benchmarkCheck(crypto.SHA1, 20, b) }
func BenchmarkCheckMD5_1(b *testing.B)     { benchmarkCheck(crypto.MD5, 1, b) }
func BenchmarkCheckMD5_20(b *testing.B)    { benchmarkCheck(crypto.MD5, 20, b) }

func benchmarkCheck(alg crypto.Hash, zeroLen uint8, b *testing.B) {
	arr := Gen(alg, zeroLen)
	key := []byte{0}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = Check(arr, key, 10)
	}
}

func BenchmarkSearchSHA256_1(b *testing.B)  { benchmarkSearch(crypto.SHA256, 1, b) }
func BenchmarkSearchSHA256_10(b *testing.B) { benchmarkSearch(crypto.SHA256, 10, b) }
func BenchmarkSearchSHA1_1(b *testing.B)    { benchmarkSearch(crypto.SHA1, 1, b) }
func BenchmarkSearchSHA1_10(b *testing.B)   { benchmarkSearch(crypto.SHA1, 10, b) }
func BenchmarkSearchMD5_1(b *testing.B)     { benchmarkSearch(crypto.MD5, 1, b) }
func BenchmarkSearchMD5_10(b *testing.B)    { benchmarkSearch(crypto.MD5, 10, b) }

func benchmarkSearch(alg crypto.Hash, zeroLen uint8, b *testing.B) {
	arr := Gen(alg, zeroLen)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = SearchKey(arr, 10)
	}
}

//goos: linux
//goarch: amd64
//pkg: pow/pkg/pow
//cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
//BenchmarkCheckSHA256_1
//BenchmarkCheckSHA256_1-4     	 1134361	       967.8 ns/op
//BenchmarkCheckSHA256_20
//BenchmarkCheckSHA256_20-4    	 1237808	       959.9 ns/op
//BenchmarkCheckSHA1_1
//BenchmarkCheckSHA1_1-4       	 2614327	       450.5 ns/op
//BenchmarkCheckSHA1_20
//BenchmarkCheckSHA1_20-4      	 2573126	       479.6 ns/op
//BenchmarkCheckMD5_1
//BenchmarkCheckMD5_1-4        	 3113239	       456.8 ns/op
//BenchmarkCheckMD5_20
//BenchmarkCheckMD5_20-4       	 3106795	       415.2 ns/op
//BenchmarkSearchSHA256_1
//BenchmarkSearchSHA256_1-4    	    4503	    261607 ns/op
//BenchmarkSearchSHA256_10
//BenchmarkSearchSHA256_10-4   	    2276	    513525 ns/op
//BenchmarkSearchSHA1_1
//BenchmarkSearchSHA1_1-4      	    8616	    128587 ns/op
//BenchmarkSearchSHA1_10
//BenchmarkSearchSHA1_10-4     	    5320	    212444 ns/op
//BenchmarkSearchMD5_1
//BenchmarkSearchMD5_1-4       	   10000	    106684 ns/op
//BenchmarkSearchMD5_10
//BenchmarkSearchMD5_10-4      	    5397	    219848 ns/op
