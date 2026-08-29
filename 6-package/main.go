package main

import (
	"fmt"

	"module-name/bins"
)

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

	binList := bins.NewBinList([]bins.Bin{
		bin1,
		bin2,
	})

	fmt.Println(bin1)
	fmt.Println(bin2)
	fmt.Println(binList)
}
