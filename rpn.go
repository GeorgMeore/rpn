package main

type Op struct {
	sym   string // operator string
	prec  int    // precedence (>= 0)
	right bool   // associativity (false - left, true - right)
}

// NOTE: operators with the same precedence must have the same associativity
func (o1 Op) Less(o2 Op) bool {
	return o2.prec > o1.prec || o1.prec == o2.prec && !o2.right
}

// NOTE: parenthesis balance check is not performed
func ToRPN(expr []string, ops... Op) []string {
	opmap := make(map[string]Op)
	for _, op := range ops {
		opmap[op.sym] = op
	}
	rpn, opstack := []string{}, []Op{}
	for _, s := range expr {
		if op, isop := opmap[s]; isop {
			i := len(opstack) - 1
			for ; i >= 0 && op.Less(opstack[i]); i-- {
				rpn = append(rpn, opstack[i].sym)
			}
			opstack = append(opstack[:i+1], op)
		} else if s == "(" {
			opstack = append(opstack, Op{sym: "(", prec: -1})
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
