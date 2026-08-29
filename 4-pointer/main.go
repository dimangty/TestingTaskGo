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

func validateCurrency(currency *string, rates *map[string]map[string]float64) bool {
	*currency = strings.ToUpper(*currency)

	_, exists := (*rates)[*currency]
	return exists
}

func getCurrencyInput(prompt string, currency *string, rates *map[string]map[string]float64) {
	for {
		fmt.Print(prompt + " (USD/EUR/RUB): ")
		fmt.Scanf("%s", currency)

		if validateCurrency(currency, rates) {
			return
		}

		fmt.Println("Ошибка: введите корректную валюту (USD, EUR или RUB)")
	}
}

func getAmountInput(amount *float64) {
	var input string

	for {
		fmt.Print("Введите сумму: ")
		fmt.Scanf("%s", &input)

		value, err := strconv.ParseFloat(input, 64)
		if err != nil || value <= 0 {
			fmt.Println("Ошибка: введите корректную положительную сумму")
			continue
		}

		*amount = value
		return
	}
}

func convertCurrency(
	amount *float64,
	fromCurrency *string,
	toCurrency *string,
	rates *map[string]map[string]float64,
	result *float64,
) {
	rate := (*rates)[*fromCurrency][*toCurrency]
	*result = *amount * rate
}

func main() {
	fmt.Println("=== КАЛЬКУЛЯТОР ВАЛЮТ ===")
	fmt.Println()

	var fromCurrency string
	var amount float64
	var toCurrency string
	var result float64

	fmt.Println("Шаг 1: Выберите исходную валюту")
	getCurrencyInput("Исходная валюта", &fromCurrency, &conversionRates)

	fmt.Println()
	fmt.Println("Шаг 2: Введите сумму для конвертации")
	getAmountInput(&amount)

	fmt.Println()
	fmt.Println("Шаг 3: Выберите целевую валюту")
	getCurrencyInput("Целевая валюта", &toCurrency, &conversionRates)

	fmt.Println()
	convertCurrency(&amount, &fromCurrency, &toCurrency, &conversionRates, &result)

	fmt.Printf("Результат: %.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}
