package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Калькулятор статистических операций ===")
		fmt.Println("Доступные операции: AVG, SUM, MED")
		fmt.Print("Введите операцию (или 'exit' для выхода): ")
		
		operation, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}
		
		operation = strings.TrimSpace(strings.ToUpper(operation))
		
		if operation == "EXIT" {
			fmt.Println("Программа завершена.")
			break
		}
		
		if operation != "AVG" && operation != "SUM" && operation != "MED" {
			fmt.Println("Ошибка: Неверная операция. Доступны: AVG, SUM, MED")
			continue
		}
		
		fmt.Print("Введите числа через запятую (например: 2, 10, 9): ")
		numbersInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}
		
		numbersInput = strings.TrimSpace(numbersInput)
		
		// Разбиваем строку по запятым
		parts := strings.Split(numbersInput, ",")
		
		// Преобразуем строки в числа
		var numbers []float64
		var invalidInputs []string
		
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			
			num, err := strconv.ParseFloat(part, 64)
			if err != nil {
				invalidInputs = append(invalidInputs, part)
				continue
			}
			numbers = append(numbers, num)
		}
		
		if len(invalidInputs) > 0 {
			fmt.Printf("Предупреждение: Некорректные значения '%s' пропущены\n", 
				strings.Join(invalidInputs, "', '"))
		}
		
		if len(numbers) == 0 {
			fmt.Println("Ошибка: Не введено ни одного корректного числа")
			continue
		}
		
		// Выполняем расчет в зависимости от операции
		var result float64
		var operationName string
		
		switch operation {
		case "AVG":
			result = calculateAverage(numbers)
			operationName = "Среднее арифметическое"
		case "SUM":
			result = calculateSum(numbers)
			operationName = "Сумма"
		case "MED":
			result = calculateMedian(numbers)
			operationName = "Медиана"
		}
		
		// Выводим результат
		fmt.Printf("\n%s: %.2f\n", operationName, result)
		fmt.Printf("Обработано чисел: %d\n", len(numbers))
		if len(numbers) > 0 {
			fmt.Printf("Исходные данные: %v\n", numbers)
		}
	}
}

// calculateAverage вычисляет среднее арифметическое
func calculateAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	sum := calculateSum(numbers)
	return sum / float64(len(numbers))
}

// calculateSum вычисляет сумму чисел
func calculateSum(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// calculateMedian вычисляет медиану
func calculateMedian(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	
	// Создаем копию для сортировки, чтобы не изменять исходный слайс
	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	sort.Float64s(sorted)
	
	n := len(sorted)
	if n%2 == 1 {
		// Нечетное количество - средний элемент
		return sorted[n/2]
	}
	// Четное количество - среднее двух средних элементов
	return (sorted[n/2-1] + sorted[n/2]) / 2.0
}