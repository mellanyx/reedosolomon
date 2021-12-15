package reedosolomon

import "log"

// calcErrorLocatorPoly - compute the errors locator polynomial from the errors positions as input
// ( вычисляем полином локатора ошибок по позициям ошибок в качестве входных данных )
func calcErrorLocatorPoly(errorPositions []int) []int {
	erasureLocations := []int{1}
	// erasures location = product(1 - x*alpha**i) for i in error positions (alpha is the alpha choosen to eval polynomials)

	// ( erasures location = product (1 - x * alpha ** i) для i в ошибочных позициях (alpha - это такая alpha, выбранная для оценки многочленов) )
	for _, p := range errorPositions {
		erasureLocations = GFPolyMult(erasureLocations, GFPolyAddition([]int{1}, []int{GFPow(2, p), 0}))
	}
	return erasureLocations
}

func calcErrorPoly(syndrom, erasureLocations []int, nsym int) []int {
	// compute the error evaluator polynomial Omega from the
	// syndrome locator Sigma
	// Omega(x) = [ Synd(x) * Error_loc(x) ] mod x^(n-k+1)

	// ( вычисляем полином Omega оценщика ошибок из
	// локатор синдромов Sigma )
	placeholder := make([]int, nsym+1)
	placeholder = append([]int{1}, placeholder...)

	_, remainder := gfPolyDivision(GFPolyMult(syndrom, erasureLocations), placeholder)

	return remainder
}

// unknownErrorLocator - Find error locator and evaluator polynomials with Berlekamp-Massey algorithm
// ( Поиск многочленов локатора и вычислителя ошибок с помощью алгоритма Берлекампа-Месси )
func unknownErrorLocator(synd []int, nsym int) []int {
	// The idea is that BM will iteratively estimate the error locator polynomial
	// To do this, it will compute a Discrepancy term called Delta, which will tell
	// if the error locator polynomial needs update or not

	// ( Идея состоит в том, что BM будет итеративно оценивать полином локатора ошибок
	// Для этого он вычислит член несоответствия под названием Дельта, который скажет
	// если полином локатора ошибок нуждается в обновлении или нет )
	errLocator := []int{1} // Sigma
	oldLocator := []int{1} // BM - итерационный алгоритм, значения предыдущих итераций Sigma

	syndShift := 0
	if len(synd) > nsym {
		syndShift = len(synd) - nsym
	}

	for i := 0; i < nsym; i++ {
		K := i + syndShift
		// compute the discrepance Delta
		// ( вычисление дельты несоответствия )
		delta := synd[K]
		for j := 1; j < len(errLocator); j++ {
			delta ^= GFMult(errLocator[len(errLocator)-(j+1)], synd[K-j])
		}

		// Shift polynomials to compute next degree
		// Полиномы сдвига для вычисления следующей степени
		oldLocator = append(oldLocator, 0)

		// iteratively estimate the errata locator and evaluator polynomials
		// ( оценка полиномов локатора ошибок и вычислителя )
		if delta != 0 {
			if len(oldLocator) > len(errLocator) {
				// computing Sigma
				// ( вычисление Sigma )
				newLocator := GFPolyScale(oldLocator, delta)
				oldLocator = GFPolyScale(errLocator, gfInverse(delta))
				errLocator = newLocator
			}

			// update with the discrepancy
			// ( обновление с не соответствием )
			errLocator = GFPolyAddition(errLocator, GFPolyScale(oldLocator, delta))
		}
	}

	// drop leading zeroes
	// ( отбрасываем ведущие нули )
	for len(errLocator) > 0 && errLocator[0] == 0 {
		errLocator = errLocator[1:]
	}

	errs := len(errLocator) - 1
	if (errs * 2) > nsym {
		log.Printf("Too many errors to correct: %d\n", errs)
	}

	return errLocator
}

func findErrors(errLocator []int, messageLen int) []int {
	// find the roots of polynomial by brute-force iter
	// ( поиск корней многочлена с помощью перебора )
	errs := len(errLocator) - 1
	errPos := []int{}

	for i := 0; i < messageLen; i++ {
		if gfPolyEvaluate(errLocator, GFPow(2, i)) == 0 {
			errPos = append(errPos, messageLen-1-i)
		}
	}

	if len(errPos) != errs {
		log.Println("Too many (or few) errors found by Chien Search")
	}
	return errPos
}

func correctErrors(msg, syndrom, errPos []int) []int {
	coefPos := make([]int, len(errPos))

	for i, p := range errPos {
		coefPos[i] = len(msg) - 1 - p
	}

	// compute the error locator polynomial
	// ( вычисление полинома локатора ошибок )
	errorLocatorPolynomial := calcErrorLocatorPoly(coefPos)

	// reverse errLocator
	reverse(syndrom)
	errorPolynomial := calcErrorPoly(syndrom, errorLocatorPolynomial, len(errorLocatorPolynomial)-1)

	// get the error location polynomial from the error positions in errPos
	// ( получение полинома местоположения ошибки из позиций ошибки в errPos )
	locationPolynomial := []int{}
	for i := 0; i < len(coefPos); i++ {
		l := 255 - coefPos[i]
		locationPolynomial = append(locationPolynomial, GFPow(2, -l))
	}

	// Forney algorithm: compute the magnitudes
	// ( Алгоритм Форни: вычисление величин )
	E := forneyAlgo(msg, errorPolynomial, locationPolynomial, errPos)

	// Simply add correction vector to message
	// ( Добавление вектора коррекции к сообщению )
	msg = GFPolyAddition(msg, E)
	return msg
}
