package main

import "fmt"

func main() {
	// fmt.Printf("%c", '!')

	var n int

	scanf, err := fmt.Scanf("%d", &n)
	if err != nil {
		return
	}

	fmt.Println(n * 2)

	switch n {
	case 7:
		fmt.Print("It is seven!")
	case 13, 17:
		fmt.Print("It is either thirteen or seventeen.")
	default:
		fmt.Print("Nothing matched")
	}
}
