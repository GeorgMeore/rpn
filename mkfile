PROG=rpn

default: main.go
	go build -o $PROG main.go

clean:
	rm -f $PROG
