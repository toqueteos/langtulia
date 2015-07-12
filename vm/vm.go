package vm

import (
	"fmt"
)

const (
	HALT int = iota
	PUSH
	POP
	ADD
	PRINT
	JMPEQ
	JMPNE
	J
)

type op struct {
	name  string
	nargs int
}

var ops = map[int]op{
	PUSH:  op{"push", 1},
	ADD:   op{"add", 0},
	PRINT: op{"print", 0},
	HALT:  op{"halt", 0},
	JMPEQ: op{"jmpeq", 2},
	JMPNE: op{"jmpne", 2},
	J:     op{"j", 1},
}

type VM struct {
	code []int
	pc   int

	stack []int
	sp    int

	trace bool
}

func New(code []int) *VM {
	return &VM{
		stack: make([]int, 128),
		sp:    -1,
		code:  code,
		pc:    0,
		trace: false,
	}
}

func (v *VM) maybeTrace() {
	if !v.trace {
		return
	}
	addr := v.pc
	op := ops[v.code[v.pc]]
	args := v.code[v.pc+1 : v.pc+op.nargs+1]
	stack := v.stack[0 : v.sp+1]

	fmt.Printf("%04d: %s %v \t%v\n", addr, op.name, args, stack)
}

func (v *VM) Run() {
	for {
		v.maybeTrace()

		// Fetch
		op := v.code[v.pc]
		v.pc++

		// Decode
		switch op {
		case PUSH:
			val := v.code[v.pc]
			v.pc++

			v.sp++
			v.stack[v.sp] = val
		case ADD:
			a := v.stack[v.sp]
			v.sp--
			b := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a + b
		case PRINT:
			val := v.stack[v.sp]
			v.sp--
			fmt.Println(val)
		case JMPEQ:
			eq := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] == eq {
				v.pc = addr
			}
		case JMPNE:
			ne := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] != ne {
				v.pc = addr
			}
		case J:
			addr := v.code[v.pc]
			// v.pc++
			v.pc = addr
		case HALT:
			return
		}
	}
}
