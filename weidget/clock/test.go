package clock

import "fmt"

func Testascii() {
	block := joinHorizontal(asciiDigits[1], asciiDigits[2], asciiDigits[3])
	fmt.Println(block)
}
