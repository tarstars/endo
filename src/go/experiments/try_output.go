package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stdout, "stdout")
	fmt.Fprintln(os.Stderr, "stderr")
}

/*


	// fmt.Printf("%c", '!')

	var n int

	_, err := fmt.Scanf("%d", &n)
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



*/
