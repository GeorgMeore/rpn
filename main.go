// Don't know why I did this but I did

package main

import (
	"fmt"
	"os"
)

type operator struct {
	sym   string // operator string
	prec  int    // precedence (>= 0)
	right bool   // associativity (false - left, true - right)
}

// NOTE: operators with the same precedence must have the same associativity
func (o1 operator) less(o2 operator) bool {
	return o2.prec > o1.prec || o1.prec == o2.prec && !o2.right
}

// NOTE: parenthesis balance check is not performed
func toRPN(expr []string, ops... operator) []string {
	opmap := make(map[string]operator)
	for _, op := range ops {
		opmap[op.sym] = op
	}
	rpn, opstack := []string{}, []operator{}
	for _, s := range expr {
		if op, isop := opmap[s]; isop {
			i := len(opstack) - 1
			for ; i >= 0 && op.less(opstack[i]); i-- {
				rpn = append(rpn, opstack[i].sym)
			}
			opstack = append(opstack[:i+1], op)
		} else if s == "(" {
			opstack = append(opstack, operator{sym: "(", prec: -1})
		} else if s == ")" {
			i := len(opstack) - 1
			for ; i > 0 && opstack[i].sym != "("; i-- {
				rpn = append(rpn, opstack[i].sym)
			}
			opstack = opstack[:i]
		} else {
			rpn = append(rpn, s)
		}
	}
	for i := len(opstack)-1; i >= 0; i-- {
		rpn = append(rpn, opstack[i].sym)
	}
	return rpn
}

// split line by spaces and brackets
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

// join words with space
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

// TODO: add support for validation
// TODO: add support for unary opearators
// TODO: read operators from file (or better yet - arguments)
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s EXPR\n", os.Args[0])
		os.Exit(1)
	}
	infix := split(os.Args[1])
	rpn := toRPN(infix,
		operator{sym: ",", prec: 4},
		operator{sym: "!", prec: 3},
		operator{sym: "/", prec: 2},
		operator{sym: "*", prec: 2},
		operator{sym: "-", prec: 1},
		operator{sym: "+", prec: 1},
	)
	fmt.Println(join(" ", rpn))
}
