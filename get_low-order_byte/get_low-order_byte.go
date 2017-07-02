package main

import (
	"fmt"
)

func getFirstLowOrderByte() {
	numOfBits := uint8(8)          // number of bits to be shifted
	mask := uint(1)<<numOfBits - 1 // The shift operators shift the left operand by the shift count specified by the right operand.
	hash := uint(1406845795978887831)
	firstLowOrderByte := hash & mask

	fmt.Printf("M: %064b\n", mask)
	fmt.Printf("H: %064b\n", hash)
	fmt.Printf("L: %064b\n", firstLowOrderByte)
}

func getFirstHighOrderByte() {
	numOfBits := uint8(8)*8 - 8 // number of bits to be shifted.
	hash := uint(1406845795978887831)
	firstHighOrderByte := uint8(hash >> numOfBits) // The shift operators shift the left operand by the shift count specified by the right operand.

	fmt.Printf("H: %064b\n", hash)
	fmt.Printf("H: %064b\n", firstHighOrderByte)
}

func main() {
	getFirstLowOrderByte()
	getFirstHighOrderByte()
}
