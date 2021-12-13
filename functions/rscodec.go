package reedosolomon

import (
	"log"
)

// RSCodec - Кодер-декодер Рида-Соломона
// ( Reed-Solomon coder/decoder )
type RSCodec struct {
	// PrimitivePoly - Десятичное представление примитивного полинома для создания таблицы поиска
	// ( Decimal representation of primitive polynomial to create lookup table )
	PrimitivePoly int
	// EccSymbols - Количество дополнительных символов
	// ( Number of additional characters )
	EccSymbols int
}

var exponentsTable = make([]int, 512)
var logsTable = make([]int, 256)

// InitTables - заполняет экспоненциальные и логарифмические таблицы
func (rs *RSCodec) InitTables() {
	// RUS //
	// Предварительно вычисляем логарифм и анти логарифмические таблицы для более быстрого вычисления, используя предоставленный примитивный полином.
	// b ** (log_b (x), log_b (y)) == x * y, где b - основание или генератор логарифма =>
	// мы можем использовать любое значение b для предварительного вычисления логарифмических и анти логарифмических таблиц, используемых для умножения двух чисел x и y.

	// EN //
	// Precompute the logarithm and anti-log tables for faster computation, using the provided primitive polynomial.
	// b**(log_b(x), log_b(y)) == x * y, where b is the base or generator of the logarithm =>
	// we can use any b to precompute logarithm and anti-log tables to use for multiplying two numbers x and y.
	initEL := 1
	for i := 0; i < 255; i++ {
		exponentsTable[i] = initEL
		logsTable[initEL] = i
		initEL = russianPeasantMult(initEL, 2, rs.PrimitivePoly, 256, true)
	}

	for i := 255; i < 512; i++ {
		exponentsTable[i] = exponentsTable[i-255]
	}
}

// RSEncode - кодируем данное сообщение кодом Рида-Соломона
// ( encode this message with the Reed-Solomon code )
func (rs *RSCodec) RSEncode(arByte []byte) (encoded []int) {
	transformedByteAr := make([]int, len(arByte))

	irredPoly := PolyGen(rs.EccSymbols)
	placeholder := make([]int, len(irredPoly)-1)

	for i, j := range arByte { transformedByteAr[i] = int(j) }

	// Дополнение сообщения и разделиние его на неприводимый порождающий многочлен
	// Adding the message and splitting it into an irreducible generator polynomial
	_, residue := gfPolyDivision(append(transformedByteAr, placeholder...), irredPoly)

	encoded = append(transformedByteAr, residue...)
	return
}

// RSDecode - декодирование и коррекция ошибок в сообщении
// ( decoding and error correction in a message )
func (rs *RSCodec) RSDecode(arStart []int) ([]int, []int) {
	arDecoded := arStart

	if len(arStart) > 255 {
		log.Fatalf("Сообщение слишком длинное, максимально допустимый размер (The message is too long, the maximum size allowed) %d\n", 255)
	}

	syndrom := calcSyndromes(arStart, rs.EccSymbols)

	if checkSyndromes(syndrom) {
		m := len(arDecoded) - rs.EccSymbols
		return arDecoded[:m], arDecoded[m:]
	}

	// вычисление полинома локатора ошибок с помощью Берлекампа-Масси
	// calculating the error locator polynomial using Berlekamp-Massey
	errLoc := unknownErrorLocator(syndrom, rs.EccSymbols)

	reverse(errLoc)
	errPos := findErrors(errLoc, len(arDecoded))

	arDecoded = correctErrors(arDecoded, syndrom, errPos)

	slice := len(arDecoded) - rs.EccSymbols

	syndrom = calcSyndromes(arDecoded, rs.EccSymbols)

	if !checkSyndromes(syndrom) {
		log.Fatalf("Could not correct message\n")
	}

	return arDecoded[:slice], arDecoded[slice:]
}