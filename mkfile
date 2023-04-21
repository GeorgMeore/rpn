PROG=rpn

default:VQ:
	go build -o $PROG

clean:
	rm -f $PROG
