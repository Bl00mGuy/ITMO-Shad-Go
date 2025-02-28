//go:build !solution

package genericsum

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"math/cmplx"
)

func Min[T constraints.Ordered](first, second T) T {
	if first < second {
		return first
	}
	return second
}

func SortSlice[T constraints.Ordered](elements []T) {
	slices.Sort(elements)
}

func MapsEqual[M1, M2 ~map[K]V, K, V comparable](map1 M1, map2 M2) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value1 := range map1 {
		if value2, ok := map2[key]; !ok || value1 != value2 {
			return false
		}
	}
	return true
}

func SliceContains[E comparable](slice []E, element E) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func MergeChans[T any](channels ...<-chan T) <-chan T {
	mergedChannel := make(chan T)
	go func() {
		openChannels := len(channels)
		for {
			for _, channel := range channels {
				select {
				case value, ok := <-channel:
					if ok {
						mergedChannel <- value
						continue
					}
					openChannels--
					if openChannels < 1 {
						close(mergedChannel)
						return
					}
				default:
					continue
				}
			}
		}
	}()
	return mergedChannel
}

type Numeric interface {
	constraints.Integer | constraints.Complex | constraints.Float
}

func CompareComplex[T Numeric](firstComplex, secondComplex T) bool {
	switch value := any(firstComplex).(type) {
	case complex64:
		if secondValue, ok := any(secondComplex).(complex64); ok {
			return real(value) == real(secondValue) && imag(value) == -imag(secondValue)
		}
	case complex128:
		if secondValue, ok := any(secondComplex).(complex128); ok {
			return cmplx.Conj(value) == secondValue
		}
	}
	return false
}

func CompareNumbers[T Numeric](firstNumber, secondNumber T) bool {
	return firstNumber == secondNumber
}

func isComplex[T Numeric](element T) bool {
	_, isComplex := any(element).(complex64)
	if isComplex {
		return true
	}
	_, isComplex = any(element).(complex128)
	return isComplex
}

func IsHermitianMatrix[T Numeric](matrix [][]T) bool {
	rows, columns := len(matrix), len(matrix[0])
	if rows != columns {
		return false
	}

	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			element := matrix[row][column]
			transposeElement := matrix[column][row]

			if isComplex(element) {
				if !CompareComplex(element, transposeElement) {
					return false
				}
			} else {
				if !CompareNumbers(element, transposeElement) {
					return false
				}
			}
		}
	}
	return true
}
