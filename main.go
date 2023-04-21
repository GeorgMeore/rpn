// Don't know why I did this but I did

package main

import (
	"fmt"
	"os"
)

func split(s string) []string {
	tokens := []string{}
	for pos, start := 0, 0; pos <= len(s); pos++ {
		if pos == len(s) || s[pos] == ' ' || s[pos] == '\t' {
			if pos > start {
				tokens = append(tokens, s[start:pos])
			}
			start = pos + 1
		} else if s[pos] == '(' || s[pos] == ')' {
			if pos > start {
				tokens = append(tokens, s[start:pos])
			}
			tokens = append(tokens, s[pos:pos+1])
			start = pos + 1
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
