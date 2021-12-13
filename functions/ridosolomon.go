package reedosolomon

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	_ "github.com/cheggaaa/pb/v3"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"time"
)

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func ReadFile(filePath string) ([]byte, string)  {
	fileExt := filepath.Ext(filePath)

	arByte, err := FileToArByte(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return arByte, fileExt
}

func WriteFile(arResultByte []byte, fileExt string, operation string) {
	if ArByteToFile(arResultByte, fmt.Sprintf("%s_File", operation), fileExt, 0644) == nil {
		fmt.Println(fmt.Sprintf("%s_File is written !", operation))
	}
}

func FileToArByte(filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	return b, nil
}

func ArByteToFile(ar []byte, name string, ext string, permissions fs.FileMode) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s%s", name, ext), ar, permissions)
	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}

func CollectArByteFile(arByte []byte, eccsyb int) [][]byte {
	var collectArByte [][]byte

	steps := len(arByte) / (255 - eccsyb)

	for i := 0; i <= steps; i++ {
		if len(arByte) >= (255 - eccsyb) {
			collectArByte = append(collectArByte, arByte[0:(255 - eccsyb)])
			arByte = arByte[(255 - eccsyb):]
		} else {
			collectArByte = append(collectArByte, arByte[0:len(arByte)])
			arByte = arByte[len(arByte):]
		}
	}

	return collectArByte
}

func CollectArByteNotEccFile(arByte []byte) [][]byte {
	var collectArByte [][]byte

	steps := len(arByte) / 255

	for i := 0; i <= steps; i++ {
		if len(arByte) >= 255 {
			collectArByte = append(collectArByte, arByte[0:255])
			arByte = arByte[255:]
		} else {
			collectArByte = append(collectArByte, arByte[0:len(arByte)])
			arByte = arByte[len(arByte):]
		}
	}

	return collectArByte
}

// EncodeByteArray - Кодирование и повреждение битового массива файла.
// В первом аргументе указываем один из двух примитивных многочленов в десятичном представлении (285 или 301).
// Во втором аргументе указываем количество добавочных символов, оно равно в двое больше количества предполагаемых ошибок.
// В третьем аргументе передаем многомерный массив байт.
//
// ( File bitmap encoding and corruption.
// In the first argument, we specify one of the two primitive polynomials in decimal notation (285 or 301).
// In the second argument, we specify the number of additional characters, it is equal to two more than the number of expected errors.
// The third argument is a multidimensional byte array. )
func EncodeByteArray(arByte [] byte, PrimitivePoly int, EccSymbols int) []byte {
	startTime := time.Now()

	rs := RSCodec {
		// Мы используем GF(2^8), потому что каждое кодовое слово занимает 8 бит
		// Можно использовать два приводимых многочлена в десятичном представлении
		// 285 обычно используется для QR-кодов
		// 301 обычно используется для Data Matrix
		PrimitivePoly:  PrimitivePoly,

		// EccSymbols - Кол-во добавочных символов
		// Кол-во ошибок, которое код сможет исправить = EccSymbols / 2
		EccSymbols: EccSymbols,
	}

	rs.InitTables()

	collectArByte := CollectArByteFile(arByte, EccSymbols)

	// create and start new bar
	count := len(collectArByte)
	bar := pb.StartNew(count)

	// Encode //
	var encodeCollectArByte [][]int

	// Итерирование многомерного массива, в котором каждым элементом
	// является массив бит размерностью (255 - EccSymbols)
	// Делается для того, что бы скармливать функции кодирования массивы максимально допустимой длины (255)
	for i := 0; i < len(collectArByte); i++ {
		encoded := rs.RSEncode(collectArByte[i])

		//if i == 0 {
		//	fmt.Println("Encoded: ", encoded)
		//}

		// Закодированный и поврежденный массив битов
		encodeCollectArByte = append(encodeCollectArByte, encoded)
		bar.Increment()
	}

	//fmt.Println("Encoded FINAL: ", encodeCollectArByte[0])

	arResultByte := UnPackArray(encodeCollectArByte)

	duration := time.Since(startTime)
	fmt.Println("Encode Runtime: ", duration)

	return arResultByte
}

func CorruptByteArray(arByte []byte, EccSymbols int) []byte {
	startTime := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	// Corrupt the message //
	// ( повреждение сообщения )

	collectArByte := CollectArByteNotEccFile(arByte)

	var arIntFile [][]int

	for i := 0; i < len(collectArByte); i++ {
		byteMessage := make([]int, len(collectArByte[i]))
		for j, ch := range collectArByte[i] {
			byteMessage[j] = int(ch)
		}

		arIntFile = append(arIntFile, byteMessage)
	}

	// create and start new bar
	count := len(arIntFile)
	bar := pb.StartNew(count)

	var corruptCollectArByte [][]int

	// errors byte
	// ( ошибочные биты )
	for i := 0; i < len(arIntFile); i++ {
		encoded := arIntFile[i]

		// corrupt the message
		// ( повреждение сообщения - ошибочные биты )
		for i := 0; i < (EccSymbols / 2); i++ {
			//rand.Seed(time.Now().UnixNano())
			//
			//randErr := rand.Intn(len(encoded)-0) + 0
			randErr := RandInt(0, len(encoded))

			encoded[randErr] = randErr
		}

		// Поврежденный массив битов
		corruptCollectArByte = append(corruptCollectArByte, encoded)
		bar.Increment()
	}

	//fmt.Println("Corrupted FINAL: ", corruptCollectArByte[0])

	arResultByte := UnPackArray(corruptCollectArByte)

	duration := time.Since(startTime)
	fmt.Println("Corrupt Runtime: ", duration)

	return arResultByte
}

// DecodeAndFixCorruptByteArray - Decoding and recovery of the file bitmap.
// In the first argument, we specify the polynomial used for encoding.
// In the second argument, we indicate the number of additional characters specified during encoding.
// In the third argument, we pass a multidimensional array of bits.
//
// ( Декодирование и восстановление битового массива файла.
// В первом аргументе указываем многочлен используемый при кодировании.
// Во втором аргументе указываем количество добавочных символов, указанное при кодировании.
// In the third argument, we pass the encoded and damaged multidimensional array. )
func DecodeAndFixCorruptByteArray(arByte []byte, PrimitivePoly int, EccSymbols int) []byte  {
	startTime := time.Now()

	// Init RS //
	rs := RSCodec {
		// Мы используем GF(2^8), потому что каждое кодовое слово занимает 8 бит
		// Можно использовать два приводимых многочлена в десятичном представлении
		// 285 обычно используется для QR-кодов
		// 301 обычно используется для Data Matrix
		PrimitivePoly:  PrimitivePoly,

		// EccSymbols - Кол-во добавочных символов
		// Кол-во ошибок, которое код сможет исправить = EccSymbols / 2
		EccSymbols: EccSymbols,
	}

	rs.InitTables()

	collectArByte := CollectArByteNotEccFile(arByte)

	var corruptCollectArByte [][]int

	for i := 0; i < len(collectArByte); i++ {
		byteMessage := make([]int, len(collectArByte[i]))
		for j, ch := range collectArByte[i] {
			byteMessage[j] = int(ch)
		}

		corruptCollectArByte = append(corruptCollectArByte, byteMessage)
	}

	// create and start new bar
	countDec := len(corruptCollectArByte)
	barDec := pb.StartNew(countDec)

	var decodedCollectArByte [][]int

	for i := 0; i < len(corruptCollectArByte); i++ {
		decoded, _ := rs.RSDecode(corruptCollectArByte[i])

		decodedCollectArByte = append(decodedCollectArByte, decoded)
		barDec.Increment()
	}

	//fmt.Println("Decoded FINAL: ", decodedCollectArByte[0])

	arResultByte := UnPackArray(decodedCollectArByte)

	duration := time.Since(startTime)
	fmt.Println("Decode Runtime: ", duration)

	return  arResultByte
}

// UnPackArray - Unpack decodedCollectArByte into one array to create a file
//
// ( Распаковываем decodedCollectArByte в один массив для создания файла )
func UnPackArray(decodedCollectArByte [][]int) []byte {
	var arResultInt []int

	for i := 0; i < len(decodedCollectArByte); i++ {
		arResultInt = append(arResultInt, decodedCollectArByte[i]...)
	}

	arResultByte := make([]byte, len(arResultInt))
	for i, ch := range arResultInt {
		arResultByte[i] = byte(ch)
	}

	return arResultByte
}