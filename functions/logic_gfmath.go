package reedosolomon

import "errors"

// GFMult - Multiplication GF
// ( Умножение GF )
func GFMult(x, y int) int {
	if x == 0 || y == 0 {
		return 0
	}
	return exponentsTable[logsTable[x]+logsTable[y]]
}

// gfPolyAddition - Polynomial addition
// ( добавленеи полинома )
func gfPolyAddition(p, q []int) (result []int) {
	if len(p) > len(q) {
		result = make([]int, len(p))
	} else {
		result = make([]int, len(q))
	}
	for i := 0; i < len(p); i++ {
		result[i+len(result)-len(p)] = p[i]
	}
	for i := 0; i < len(q); i++ {
		result[i+len(result)-len(q)] ^= q[i]
	}
	return
}

// gfDivision - Division GF
// ( Деление GF )
func gfDivision(x, y int) (int, error) {
	if y == 0 {
		return -1, errors.New("Zero division")
	}
	if x == 0 {
		return 0, nil
	}
	return exponentsTable[(logsTable[x]+255-logsTable[y])%255], nil
}

// GFPow - Modular division GF
// ( Деление по модулю GF )
func GFPow(x, power int) int {
	return exponentsTable[negmod(logsTable[x]*power, 255)]
}

// gfPolyScale - multiply polynomial by scalar
// ( умножение полинома на скаляр )
func gfPolyScale(p []int, x int) []int {
	result := make([]int, len(p))
	for i := 0; i < len(p); i++ {
		result[i] = GFMult(p[i], x)
	}
	return result
}

func gfInverse(x int) int {
	return exponentsTable[255-logsTable[x]]
}

// GFDeduction - Addition GF
// ( Добавление GF )
func gfAddition(x, y int) int {
	return x ^ y
}

// gfPolyEvaluate - Evaluates a polynomial in GF(2^p) given the value for x.
// This is based on Horner's scheme for maximum efficiency.
// ( Вычисляет многочлен в GF (2 ^ p) по значению x.
// Это основано на схеме Хорнера для максимальной эффективности.)
func gfPolyEvaluate(p []int, x int) int {
	// example: 01 x4 + 0f x3 + 36 x2 + 78 x + 40 = (((01 x + 0f) x + 36) x + 78) x + 40
	// пример: 01 x4 + 0f x3 + 36 x2 + 78 x + 40 = (((01 x + 0f) x + 36) x + 78) x + 40
	y := p[0]
	for i := 1; i < len(p); i++ {
		y = GFMult(y, x) ^ p[i]
	}
	return y
}

// GFPolyMult - multiply two polynomials inside Galois Field
// (умножение двух многочленов в поле Галуа)
func GFPolyMult(p, q []int) (result []int) {
	result = make([]int, len(p)+len(q)-1)
	// compute the polynomial multiplication like product of two vectors
	for j := 0; j < len(q); j++ {
		for i := 0; i < len(p); i++ {
			result[i+j] ^= GFMult(p[i], q[j])
		}
	}
	return
}

// GFDeduction - Deduction GF
// ( Вичитание GF )
func GFDeduction(x, y int) int {
	return x ^ y
}

func gfPolyDivision(divident, divisor []int) ([]int, []int) {
	// Fast polynomial division by using Extended Synthetic Division and optimized for GF(2^p) computations
	// (Быстрое полиномиальное деление с использованием расширенного синтетического деления и оптимизация для вычислений GF (2 ^ p))
	result := make([]int, len(divident))
	copy(result, divident)

	for i := 0; i < len(divident)-(len(divisor)-1); i++ {
		coef := result[i]
		if coef != 0 {
			for j := 1; j < len(divisor); j++ {
				if divisor[j] != 0 {
					result[i+j] ^= GFMult(divisor[j], coef)
				}
			}
		}
	}
	separator := len(divisor) - 1
	// остаток имеет тот же размер, что и делитель, так как он
	// является тем, что мы не смогли отделить от делителя
	// возвращаем частное, остаток

	// remainder has same size as the divisor, since it's
	// what we couldn't divide from the divident
	// return quotient, remainder
	return result[:len(result)-separator], result[len(result)-separator:]
}
