MKSHELL=rc
PROG=rpn

$PROG: main.go
	go build -o $PROG main.go

test:VQ: $PROG
	for (t in test/*) {
		echo -n $t' '
		if (cmp -s <{xargs -a $t/args ./$prereq <$t/stdin} $t/stdout) echo OK
		if not echo FAIL
	}

clean:
	rm -f $PROG
