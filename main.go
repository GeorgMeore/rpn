// Don't know why I did this but I did

// TODO: write README

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func errorf(format string, vs ...any) error {
	return errors.New(fmt.Sprintf(format, vs...))
}

type operator struct {
	sym   string // operator string
	prec  int    // precedence (>= 0)
	right bool   // associativity (false - left, true - right)
	unary bool   // arity (false - binary, true - unary)
}

// TODO: some research (I am not sure this is a correct extension to shunting yard, but it seems to work)
// convert infix expression to postfix
func toRPN(expr []string, opmap map[string]operator) []string {
	rpn, ops := []string{}, []operator{}
	for _, s := range expr {
		if s == "(" {
			// push pseudo-operator '(' with lowest precedence
			ops = append(ops, operator{sym: "(", prec: -1})
		} else if s == ")" {
			i := len(ops) - 1
			for ; i > 0 && ops[i].sym != "("; i-- {
				rpn = append(rpn, ops[i].sym)
			}
			ops = ops[:i]
		} else if op, isop := opmap[s]; !isop {
			rpn = append(rpn, s)
		} else if op.unary {
			if op.right {
				ops = append(ops, op)
			} else {
				i := len(ops) - 1
				for ; i >= 0 && (ops[i].unary && ops[i].prec<op.prec || !ops[i].unary && ops[i].prec>op.prec); i-- {
					rpn = append(rpn, ops[i].sym)
				}
				ops = ops[:i+1]
				rpn = append(rpn, op.sym)
			}
		} else {
			i := len(ops) - 1
			if op.right {
				for ; i >= 0 && ops[i].prec > op.prec; i-- {
					rpn = append(rpn, ops[i].sym)
				}
			} else {
				for ; i >= 0 && ops[i].prec >= op.prec; i-- {
					rpn = append(rpn, ops[i].sym)
				}
			}
			ops = append(ops[:i+1], op)
		}
	}
	for i := len(ops) - 1; i >= 0; i-- {
		rpn = append(rpn, ops[i].sym)
	}
	return rpn
}

// split line by spaces and parentheses
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

// join strings with space
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

// "expected" flags
const (
	fnothing = 1 << iota
	frparen
	flparen
	farg
	finfix
	fprefix
	fpostfix
)

// check if i-th element of infix expression is something we expect
func expect(flags int, opmap map[string]operator, infix []string, i int) bool {
	if i < 0 || i == len(infix) {
		return flags & fnothing != 0
	}
	if infix[i] == "(" {
		return flags & flparen != 0
	}
	if infix[i] == ")" {
		return flags & frparen != 0
	}
	op, isop := opmap[infix[i]]
	if !isop {
		return flags & farg != 0
	}
	if !op.unary {
		return flags & finfix != 0
	}
	if !op.right {
		return flags & fpostfix != 0
	}
	return flags & fprefix != 0
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
		} else if op, isop := opmap[s]; !isop {
			if !expect(fnothing | flparen | finfix | fprefix, opmap, infix, i - 1) {
				return errorf("expected nothing or a '(' or an binary or prefix operator before '%s'", s)
			}
			if !expect(fnothing | frparen | finfix | fpostfix, opmap, infix, i + 1) {
				return errorf("expected nothing or a ')' or an binary or postfix operator after '%s'", s)
			}
		} else if op.unary {
			if op.right {
				if !expect(fnothing | flparen | finfix | fprefix, opmap, infix, i - 1) {
					return errorf("expected nothing or a '(' or an binary or prefix operator before '%s'", s)
				}
				if !expect(flparen | farg | fprefix, opmap, infix, i + 1) {
					return errorf("expected a '(' or an argument or a prefix operator after '%s'", s)
				}
			} else {
				if !expect(frparen | farg | fpostfix, opmap, infix, i - 1) {
					return errorf("expected a ')' or an argument or a postfix operator before '%s'", s)
				}
				if !expect(fnothing | frparen | finfix | fpostfix, opmap, infix, i + 1) {
					return errorf("expected nothing or a ')' or an binary or postfix operator after '%s'", s)
				}
			}
		} else {
			if !expect(frparen | farg | fpostfix, opmap, infix, i - 1) {
				return errorf("expected a ')' or an argument or a postfix operator before '%s'", s)
			}
			if !expect(flparen | farg | fprefix, opmap, infix, i + 1) {
				return errorf("expected a '(' or an argument or a prefix operator after '%s'", s)
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
		right, unary := false, false
		for _, flag := range desc[0] {
			if flag == 'r' {
				right = true
			} else if flag == 'l' {
				right = false
			} else if flag == 'u' {
				unary = true
			} else if flag == 'b' {
				unary = false
			} else {
				return nil, errorf("'%s': unknown flag: %c", arg, flag)
			}
		}
		for _, sym := range desc[1:] {
			if len(sym) == 0 {
				return nil, errorf("'%s': empty operator", arg)
			}
			ops = append(ops, operator{sym: sym, right: right, prec: prec, unary: unary})
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
