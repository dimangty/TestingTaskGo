package main

import "fmt"

const (
	USDtoEUR = 0.86
	USDtoRUB = 81.0
	EURtoRUB = USDtoRUB / USDtoEUR
)

func readUserInput() (float64, string, string) {
	var amount float64
	var fromCurrency, toCurrency string
	
	fmt.Print("Введите сумму: ")
	fmt.Scanf("%f", &amount)
	
	fmt.Print("Исходная валюта (USD/EUR/RUB): ")
	fmt.Scanf("%s", &fromCurrency)
	
	fmt.Print("Целевая валюта (USD/EUR/RUB): ")
	fmt.Scanf("%s", &toCurrency)
	
	return amount, fromCurrency, toCurrency
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) float64 {
	return 0.0
}

func main() {
	fmt.Println("Currency conversion constants:")
	fmt.Printf("USD to EUR: %.2f\n", USDtoEUR)
	fmt.Printf("USD to RUB: %.2f\n", USDtoRUB)
	fmt.Printf("EUR to RUB: %.2f\n", EURtoRUB)
}