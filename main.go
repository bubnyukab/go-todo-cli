package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Todo CLI")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		scanner.Scan()

		line := scanner.Text()
		args := strings.Split(line, " ")
		fmt.Println(line, args)
	}
}
