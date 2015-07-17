package vm

import (
	"fmt"
)

const (
	HALT  int32 = iota // Halt
	NOP                // Nop
	PUSH               // Push to Top of Stack
	POP                // Remove from Top of Stack
	ADD                // Add
	PRINT              // Print
	JMPEQ              // Jump If Equal
	JMPNE              // Jump If Not Equal
	J                  // Jump Inconditionally
	SLL                // Shift Left Logical
	SRL                // Shift Right Logical
)

const (
	ZERO = iota
	LO
	HI
)

type op struct {
	name  string
	nargs int32
}

var ops = map[int32]op{
	NOP:   op{"nop", 0},
	PUSH:  op{"push", 1},
	POP:   op{"pop", 0},
	ADD:   op{"add", 0},
	PRINT: op{"print", 0},
	HALT:  op{"halt", 0},
	JMPEQ: op{"jmpeq", 2},
	JMPNE: op{"jmpne", 2},
	J:     op{"j", 1},
	SLL:   op{"sll", 2},
	SRL:   op{"srl", 2},
}

type VM struct {
	r [32]int32 // registers

	code []int32
	pc   int32

	stack []int32
	sp    int32

	trace bool
}

func New(code []int32) *VM {
	return &VM{
		stack: make([]int32, 128),
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
		case NOP:
			// Derp
		case PUSH:
			val := v.code[v.pc]
			v.pc++

			v.sp++
			v.stack[v.sp] = val
		case POP:
			v.sp--
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
		case SLL:
			a := uint(v.stack[v.sp])
			v.sp--
			b := uint(v.stack[v.sp])
			v.sp--

			v.sp++
			v.stack[v.sp] = int32(a << b)
		case SRL:
			a := uint(v.stack[v.sp])
			v.sp--
			b := uint(v.stack[v.sp])
			v.sp--

			v.sp++
			v.stack[v.sp] = int32(a >> b)
		case HALT:
			return
		}
	}
}
