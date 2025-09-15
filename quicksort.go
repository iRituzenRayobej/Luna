package main

import "fmt"

// Função QuickSort
func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	// Escolhe o pivô (último elemento)
	pivot := arr[len(arr)-1]
	left := []int{}
	right := []int{}

	// Divide os elementos em menores e maiores que o pivô
	for _, v := range arr[:len(arr)-1] {
		if v <= pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}

	// Recursão e concatenação
	return append(append(quickSort(left), pivot), quickSort(right)...)
}

func main() {
	arr := []int{33, 10, 55, 71, 29, 3, 18}
	fmt.Println("Antes:", arr)
	sorted := quickSort(arr)
	fmt.Println("Depois:", sorted)
}
