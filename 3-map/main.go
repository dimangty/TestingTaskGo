package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	USDtoEUR = 0.86
	USDtoRUB = 81.0
	EURtoRUB = USDtoRUB / USDtoEUR
)

var conversionRates = map[string]map[string]float64{
	"USD": {
		"USD": 1,
		"EUR": USDtoEUR,
		"RUB": USDtoRUB,
	},
	"EUR": {
		"USD": 1 / USDtoEUR,
		"EUR": 1,
		"RUB": EURtoRUB,
	},
	"RUB": {
		"USD": 1 / USDtoRUB,
		"EUR": 1 / EURtoRUB,
		"RUB": 1,
	},
}

func validateCurrency(currency string) bool {
	currency = strings.ToUpper(currency)

	_, exists := conversionRates[currency]
	return exists
}

func getCurrencyInput(prompt string) string {
	var currency string

	for {
		fmt.Print(prompt + " (USD/EUR/RUB): ")
		fmt.Scanf("%s", &currency)

		currency = strings.ToUpper(currency)

		if validateCurrency(currency) {
			return currency
		}

		fmt.Println("Ошибка: введите корректную валюту (USD, EUR или RUB)")
	}
}

func getAmountInput() float64 {
	var input string

	for {
		fmt.Print("Введите сумму: ")
		fmt.Scanf("%s", &input)

		amount, err := strconv.ParseFloat(input, 64)
		if err != nil || amount <= 0 {
			fmt.Println("Ошибка: введите корректную положительную сумму")
			continue
		}

		return amount
	}
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) float64 {
	rate := conversionRates[fromCurrency][toCurrency]
	return amount * rate
}

func main() {
	fmt.Println("=== КАЛЬКУЛЯТОР ВАЛЮТ ===")
	fmt.Println()

	fmt.Println("Шаг 1: Выберите исходную валюту")
	fromCurrency := getCurrencyInput("Исходная валюта")

	fmt.Println()
	fmt.Println("Шаг 2: Введите сумму для конвертации")
	amount := getAmountInput()

	fmt.Println()
	fmt.Println("Шаг 3: Выберите целевую валюту")
	toCurrency := getCurrencyInput("Целевая валюта")

	fmt.Println()
	result := convertCurrency(amount, fromCurrency, toCurrency)

	fmt.Printf("Результат: %.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}
