// Don't know why I did this but I did

package main

import (
	"fmt"
	"os"
)

func split(s string) []string {
	tokens := []string{}
	for len(s) > 0 {
		switch s[0] {
		case ' ':
			s = s[1:]
		case '+', '-', '/', '*', ',', '(', ')', '!':
			tokens = append(tokens, s[:1])
			s = s[1:]
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end := 0
			for end < len(s) && s[end] >= '0' && s[end] <= '9' {
				end += 1
			}
			tokens = append(tokens, s[:end])
			s = s[end:]
		default:
			s = s[1:]
		}
	}
	return tokens
}

func join(sep string, ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	joined := ss[0]
	for _, s := range ss[1:] {
		joined = joined + sep + s
	}
	return joined
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	infix := split(os.Args[1])
	rpn := ToRPN(infix,
		Op{sym: ",", prec: 4},
		Op{sym: "!", prec: 3},
		Op{sym: "*", prec: 2},
		Op{sym: "+", prec: 1},
	)
	fmt.Println(join(" ", rpn))
}
