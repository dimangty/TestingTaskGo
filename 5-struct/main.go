package main

import (
	"fmt"
	"time"
)

type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

type BinList struct {
	bins []Bin
}

func NewBin(
	id string,
	private bool,
	name string,
) Bin {
	return Bin{
		id:        id,
		private:   private,
		createdAt: time.Now(),
		name:      name,
	}
}

func NewBinList(bins []Bin) BinList {
	return BinList{
		bins: bins,
	}
}

func main() {
	bin1 := NewBin(
		"1",
		true,
		"My private bin",
	)

	bin2 := NewBin(
		"2",
		false,
		"My public bin",
	)

	binList := NewBinList([]Bin{
		bin1,
		bin2,
	})

	fmt.Println(bin1)
	fmt.Println(bin2)
	fmt.Println(binList)
}