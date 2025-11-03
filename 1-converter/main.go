package main

import "fmt"

const (
	USDtoEUR = 0.86
	USDtoRUB = 81.0
	EURtoRUB = USDtoRUB / USDtoEUR
)

func main() {
	fmt.Println("Currency conversion constants:")
	fmt.Printf("USD to EUR: %.2f\n", USDtoEUR)
	fmt.Printf("USD to RUB: %.2f\n", USDtoRUB)
	fmt.Printf("EUR to RUB: %.2f\n", EURtoRUB)
}