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
	MUL                // Multiplication
	DIV                // Division
	MOD                // Modulo
	PRINT              // Print
	JEQ                // Jump If Equal
	JNE                // Jump If Not Equal
	J                  // Jump Inconditionally
	SLT                // Set If Less Than
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
	HALT:  op{"halt", 0},
	NOP:   op{"nop", 0},
	PUSH:  op{"push", 1},
	POP:   op{"pop", 0},
	ADD:   op{"add", 0},
	MUL:   op{"mul", 0},
	DIV:   op{"div", 0},
	MOD:   op{"mod", 0},
	PRINT: op{"print", 0},
	JEQ:   op{"jeq", 2},
	JNE:   op{"jne", 2},
	J:     op{"j", 1},
	SLT:   op{"slt", 1},
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
	stats struct {
		push   int
		pop    int
		grow   int
		shrink int
	}
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

func (v *VM) resizeStack() {
	// Need a bigger stack?
	if (v.stats.push - v.stats.pop) >= (cap(v.stack) >> 1) {
		stack := make([]int32, cap(v.stack)*2)
		copy(stack, v.stack)
		v.stack = stack
		v.stats.grow++
	}
	// Want a smaller stack?
	if (v.stats.push - v.stats.pop) < (cap(v.stack) >> 2) {
		stack := make([]int32, cap(v.stack)>>1)
		copy(stack, v.stack)
		v.stack = stack
		v.stats.shrink++
	}
}

func (v *VM) Run() {
	for {
		v.maybeTrace()

		// Fetch
		op := v.code[v.pc]
		v.pc++

		// Resize stack
		v.resizeStack()

		// Decode
		switch op {
		case NOP:
			// Derp
		case PUSH:
			val := v.code[v.pc]
			v.pc++

			v.sp++
			v.stack[v.sp] = val
			v.stats.push++
		case POP:
			v.sp--
			v.stats.pop++
		case ADD:
			a := v.stack[v.sp]
			v.sp--
			b := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a + b
		case MUL:
			a := int64(v.stack[v.sp])
			v.sp--
			b := int64(v.stack[v.sp])
			v.sp--

			c := a * b
			v.r[LO] = int32((c << 32) >> 32)
			v.r[HI] = int32(c >> 32)
		case DIV:
			a := v.stack[v.sp]
			v.sp--
			b := v.stack[v.sp]
			v.sp--

			v.r[LO] = a / b
			v.r[HI] = a % b
		case MOD:
			a := v.stack[v.sp]
			v.sp--
			b := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a % b
		case PRINT:
			val := v.stack[v.sp]
			v.sp--
			fmt.Println(val)
		case JEQ:
			eq := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] == eq {
				v.pc = addr
			}
		case JNE:
			ne := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] != ne {
				v.pc = addr
			}
		case J:
			addr := v.code[v.pc]
			v.pc = addr
		case SLT:
			a := v.stack[v.sp]
			v.sp--
			b := v.code[v.pc]
			v.pc++

			v.sp++
			v.stack[v.sp] = bint32[a < b]
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

var bint32 = map[bool]int32{
	false: 0,
	true:  1,
}
