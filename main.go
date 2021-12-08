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

	// Read File //
	arByte, fileExt := reedosolomon.ReadFile(flag.Args()[1])

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

		encodedArByte := reedosolomon.EncodeByteArray(arByte, primitive, eccsymbols)

		// Write File //
		reedosolomon.WriteFile(encodedArByte, fileExt, "Encoded")
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

		corruptedArByte := reedosolomon.CorruptByteArray(arByte, eccsymbols)

		// Write File //
		reedosolomon.WriteFile(corruptedArByte, fileExt, "Corrupted")
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

		decodedArByte := reedosolomon.DecodeAndFixCorruptByteArray(arByte, primitive, eccsymbols)

		// Write File //
		reedosolomon.WriteFile(decodedArByte, fileExt, "Decoded")
	default:
		fmt.Println("Не правильно переданы параметры, см. документацию")
	}
}