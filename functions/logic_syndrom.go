package reedosolomon

func calcSyndromes(msg []int, EccSymbols int) []int {
	syndrom := make([]int, EccSymbols)
	for i := 0; i < EccSymbols; i++ {
		syndrom[i] = gfPolyEvaluate(msg, GFPow(2, i))
	}
	syndrom = append([]int{0}, syndrom...)
	return syndrom
}

func checkSyndromes(syndrom []int) bool {
	for _, v := range syndrom {
		if v > 0 {
			return false
		}
	}
	return true
}
