// Don't know why I did this but I did

// TODO: add support for unary opearators (or not)
// TODO: write README

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// this should be a builtin...
func has[T comparable, V any](m map[T]V, key T) bool {
	_, ok := m[key]
	return ok
}

func errorf(format string, vs ...any) error {
	return errors.New(fmt.Sprintf(format, vs...))
}

type operator struct {
	sym   string // operator string
	prec  int    // precedence (>= 0)
	right bool   // associativity (false - left, true - right)
}

// NOTE: operators with the same precedence must have the same associativity
func (o1 operator) less(o2 operator) bool {
	return o2.prec > o1.prec || o1.prec == o2.prec && !o2.right
}

// convert infix expression to postfix
func toRPN(expr []string, opmap map[string]operator) []string {
	rpn, opstack := []string{}, []operator{}
	for _, s := range expr {
		if op, isop := opmap[s]; isop {
			i := len(opstack) - 1
			for ; i >= 0 && op.less(opstack[i]); i-- {
				rpn = append(rpn, opstack[i].sym)
			}
			opstack = append(opstack[:i+1], op)
		} else if s == "(" {
			// push pseudo-operator '(' with lowest precedence
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
	for i := len(opstack) - 1; i >= 0; i-- {
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

// TODO: error location reporting
// check if expression is a valid infix expression
func check(infix []string, opmap map[string]operator) error {
	if len(infix) == 0 {
		return nil
	}
	parens := 0
	for i, s := range infix {
		if s == "(" {
			parens += 1
			if i < len(infix)-1 && infix[i+1] == ")" {
				return errorf("empty brackets")
			}
		} else if s == ")" {
			if parens == 0 {
				return errorf("unmatched ')'")
			}
			parens -= 1
		} else if has(opmap, s) {
			if !(i > 0 && (!has(opmap, infix[i-1]) || infix[i-1] == ")")) {
				return errorf("expected a ')' or an argument before '%s'", s)
			}
			if !(i < len(infix)-1 && (!has(opmap, infix[i+1]) || infix[i+1] == "(")) {
				return errorf("expected a '(' or an argument after '%s'", s)
			}
		} else {
			if !(i == 0 || (has(opmap, infix[i-1]) || infix[i-1] == "(")) {
				return errorf("expected nothing or a '(' or an operator before '%s'", s)
			}
			if !(i == len(infix)-1 || (has(opmap, infix[i+1]) || infix[i+1] == ")")) {
				return errorf("expected nothing or a ')' or an operator after '%s'", s)
			}
		}
	}
	if parens > 0 {
		return errorf("unmatched '('")
	}
	return nil
}

// parse operator descriptions
func getops(args []string) ([]operator, error) {
	ops := []operator{}
	for prec, arg := range args {
		desc := strings.Split(arg, ":")
		if len(desc) < 2 {
			return nil, errorf("'%s': invalid format", arg)
		}
		right := false
		for _, flag := range desc[0] {
			if flag == 'r' {
				right = true
			} else if flag == 'l' {
				right = false
			} else {
				return nil, errorf("'%s': unknown flag: %c", arg, flag)
			}
		}
		for _, sym := range desc[1:] {
			if len(sym) == 0 {
				return nil, errorf("'%s': empty operator", arg)
			}
			ops = append(ops, operator{sym: sym, right: right, prec: prec})
		}
	}
	return ops, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s OPERATORS\n", os.Args[0])
		os.Exit(1)
	}
	ops, err := getops(os.Args[1:])
	if err != nil {
		fmt.Printf("error: bad argument: %s\n", err.Error())
		os.Exit(1)
	}
	opmap := make(map[string]operator)
	for _, op := range ops {
		opmap[op.sym] = op
	}
	istty := false
	if info, _ := os.Stdout.Stat(); info.Mode()&os.ModeCharDevice != 0 {
		istty = true
	}
	scanner := bufio.NewScanner(os.Stdin)
	if istty {
		fmt.Print("> ")
	}
	for scanner.Scan() {
		infix := split(scanner.Text())
		if err := check(infix, opmap); err != nil {
			fmt.Printf("error: bad expression: %s\n", err.Error())
			continue
		}
		rpn := toRPN(infix, opmap)
		fmt.Println(join(" ", rpn))
		if istty {
			fmt.Print("> ")
		}
	}
}
