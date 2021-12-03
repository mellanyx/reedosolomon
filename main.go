package main

import (
	"flag"
	"fmt"
	"github.com/mellanyx/reedosolomon/functions"
	"log"
	"strconv"
)

func main() {
	flag.Parse()

	switch flag.Args()[0] {
	case "encode":
		// mode path primitive eccsybmols
		if len(flag.Args()) != 4 {
			log.Fatal("Не правильно переданы параметры, см. документацию")
		}

		// primitive //
		primitive, err := strconv.Atoi(flag.Args()[2])
		if err != nil {
			log.Fatal(err)
		}

		// primitive //
		eccsymbols, err := strconv.Atoi(flag.Args()[3])
		if err != nil {
			log.Fatal(err)
		}

		reedosolomon.EncodeFile(flag.Args()[1], primitive, eccsymbols)
	case "corrupt":
		// mode path eccsybmols

		if len(flag.Args()) != 3 {
			log.Fatal("Не правильно переданы параметры, см. документацию")
		}

		// primitive //
		eccsymbols, err := strconv.Atoi(flag.Args()[2])
		if err != nil {
			log.Fatal(err)
		}

		//errorsArr := map[int]int{
		//	0: 99,
		//	1: 0,
		//	2: 3,
		//	3: 68,
		//	4: 21,
		//}

		reedosolomon.CorruptFile(flag.Args()[1], eccsymbols)
	case "decode":
		// mode path primitive eccsybmols

		if len(flag.Args()) != 4 {
			log.Fatal("Не правильно переданы параметры, см. документацию")
		}

		// primitive //
		primitive, err := strconv.Atoi(flag.Args()[2])
		if err != nil {
			log.Fatal(err)
		}

		// primitive //
		eccsymbols, err := strconv.Atoi(flag.Args()[3])
		if err != nil {
			log.Fatal(err)
		}

		reedosolomon.DecodeAndFixCorruptFile(flag.Args()[1], primitive, eccsymbols)
	default:
		fmt.Println("Не правильно переданы параметры, см. документацию")
	}
}