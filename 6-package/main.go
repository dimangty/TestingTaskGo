package main

import (
	"fmt"
	"log"

	"module-name/bins"
	"module-name/storage"
)

const storagePath = "bins.json"

func main() {
	bin1 := bins.NewBin(
		"1",
		true,
		"My private bin",
	)

	bin2 := bins.NewBin(
		"2",
		false,
		"My public bin",
	)

	if err := storage.SaveBin(storagePath, bin1); err != nil {
		log.Fatal(err)
	}

	if err := storage.SaveBin(storagePath, bin2); err != nil {
		log.Fatal(err)
	}

	storedBins, err := storage.ReadBins(storagePath)
	if err != nil {
		log.Fatal(err)
	}

	binList := bins.NewBinList(storedBins)

	fmt.Println(bin1)
	fmt.Println(bin2)
	fmt.Println(binList)
}
