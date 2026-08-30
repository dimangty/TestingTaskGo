package main

import (
	"fmt"
	"log"

	"module-name/bins"
	"module-name/file"
	"module-name/storage"
)

const storagePath = "bins.json"

func main() {
	localFiles := file.NewLocal()
	binStorage := storage.NewJSONStorage(storagePath, localFiles)
	binList, err := bins.NewBinList(binStorage)
	if err != nil {
		log.Fatal(err)
	}

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

	if err := binList.Add(bin1); err != nil {
		log.Fatal(err)
	}

	if err := binList.Add(bin2); err != nil {
		log.Fatal(err)
	}

	fmt.Println(bin1)
	fmt.Println(bin2)
	fmt.Println(binList.Bins())
}
