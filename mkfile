PROG=rpn

$PROG: main.go
	go build -o $PROG

clean:
	rm -f $PROG
