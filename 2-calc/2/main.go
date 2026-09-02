package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var calculations = map[string]func([]float64) float64{
	"AVG": calculateAverage,
	"SUM": calculateSum,
	"MED": calculateMedian,
}

var operationNames = map[string]string{
	"AVG": "Среднее арифметическое",
	"SUM": "Сумма",
	"MED": "Медиана",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Калькулятор статистических операций ===")
		fmt.Println("Доступные операции: AVG, SUM, MED")
		fmt.Print("Введите операцию (или 'exit' для выхода): ")

		operationInput, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) && operationInput == "" {
			fmt.Println("\nПрограмма завершена.")
			return
		}
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("Ошибка ввода:", err)
			return
		}

		operation := strings.TrimSpace(strings.ToUpper(operationInput))
		if operation == "EXIT" {
			fmt.Println("Программа завершена.")
			return
		}

		calculate, ok := calculations[operation]
		if !ok {
			fmt.Println("Ошибка: Неверная операция. Доступны: AVG, SUM, MED")
			if errors.Is(err, io.EOF) {
				return
			}
			continue
		}

		fmt.Print("Введите числа через запятую (например: 2, 10, 9): ")
		numbersInput, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) && numbersInput == "" {
			fmt.Println("\nПрограмма завершена.")
			return
		}
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("Ошибка ввода:", err)
			return
		}

		numbers, invalidInputs := parseNumbers(numbersInput)
		if len(invalidInputs) > 0 {
			fmt.Printf(
				"Предупреждение: Некорректные значения '%s' пропущены\n",
				strings.Join(invalidInputs, "', '"),
			)
		}

		if len(numbers) == 0 {
			fmt.Println("Ошибка: Не введено ни одного корректного числа")
			if errors.Is(err, io.EOF) {
				return
			}
			continue
		}

		result := calculate(numbers)
		fmt.Printf("\n%s: %.2f\n", operationNames[operation], result)
		fmt.Printf("Обработано чисел: %d\n", len(numbers))
		fmt.Printf("Исходные данные: %v\n", numbers)

		if errors.Is(err, io.EOF) {
			return
		}
	}
}

func parseNumbers(input string) ([]float64, []string) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	numbers := make([]float64, 0, len(parts))
	invalidInputs := make([]string, 0)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		number, err := strconv.ParseFloat(part, 64)
		if err != nil {
			invalidInputs = append(invalidInputs, part)
			continue
		}
		numbers = append(numbers, number)
	}

	return numbers, invalidInputs
}

func calculateAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return calculateSum(numbers) / float64(len(numbers))
}

func calculateSum(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func calculateMedian(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	sortedNumbers := append([]float64(nil), numbers...)
	sort.Float64s(sortedNumbers)

	middle := len(sortedNumbers) / 2
	if len(sortedNumbers)%2 == 1 {
		return sortedNumbers[middle]
	}
	return (sortedNumbers[middle-1] + sortedNumbers[middle]) / 2
}
