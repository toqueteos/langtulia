package vm

import (
	"math/rand"
	"testing"
)

type TestCase struct {
	Code       []int32
	TopOfStack int32
}

func TestCommon(t *testing.T) {
	tests := []TestCase{
		{[]int32{PUSH, 2, NOP, NOP, HALT}, 2},
		{[]int32{PUSH, 2, PUSH, 3, NOP, HALT}, 3},
		{[]int32{PUSH, 2, PUSH, 3, PUSH, 4, PUSH, 5, POP, POP, HALT}, 3},
		{[]int32{PUSH, 7, PUSH, 0x1234, MOD, HALT}, 5},
	}

	for idx, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[vm.sp] != tt.TopOfStack {
			t.Errorf("%d. expected %d, got %d", idx, tt.TopOfStack, vm.stack[vm.sp])
		}
	}
}

func TestArith(t *testing.T) {
	type MulCase struct {
		Code   []int32
		Hi, Lo int32
	}
	tests := []MulCase{
		{[]int32{PUSH, 0x1234, PUSH, 0x4321, MUL, HALT}, 0, 0x04c5f4b4},
		{[]int32{PUSH, 0x66778899, PUSH, 0x66778899, MUL, HALT}, 0x290378aa, 0x3320eb71},
		{[]int32{PUSH, 0x1234, PUSH, 0x4321, DIV, HALT}, 0x00000c85, 0x00000003},
		{[]int32{PUSH, 0x66778899, PUSH, 0x66778899, DIV, HALT}, 0, 1},
	}
	for idx, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.r[HI] != tt.Hi {
			t.Errorf("%d. expected HI %#08x, got %#08x", idx, tt.Hi, vm.r[HI])
		}
		if vm.r[LO] != tt.Lo {
			t.Errorf("%d. expected LO %#08x, got %#08x", idx, tt.Lo, vm.r[LO])
		}
	}
}

func TestAddLoop(t *testing.T) {
	tests := []TestCase{
		{[]int32{PUSH, 0, PUSH, 1, ADD, HALT}, 1},
		{[]int32{PUSH, 2, PUSH, 3, ADD, JNE, 17, 2, HALT}, 17},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[vm.sp] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[vm.sp])
		}
	}
}

func TestStack(t *testing.T) {
	const N = 1 << 16
	var code []int32

	for i := 1; i < N; i++ {
		code = append(code, PUSH)
		code = append(code, rand.Int31n(N))
	}
	for i := 0; i < (N - 2); i++ {
		code = append(code, POP)
	}
	code = append(code, HALT)

	vm := New(code)
	vm.Run()
	t.Logf("stack: % d -- sp: %d", vm.stack, vm.sp)
	if vm.stack[vm.sp] != code[1] {
		t.Errorf("expected %d, got %d", code[1], vm.stack[vm.sp])
	}
}

func TestJumps(t *testing.T) {
	tests := []TestCase{
		{[]int32{PUSH, 2, J, 6, PUSH, 4, PUSH, 3, ADD, HALT}, 5},
		{[]int32{PUSH, 2, J, 5, NOP, PUSH, 3, ADD, HALT}, 5},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[vm.sp] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[vm.sp])
		}
	}
}

func TestShiftLogical(t *testing.T) {
	tests := []TestCase{
		{[]int32{PUSH, 3, PUSH, 1, SLL, HALT}, 8},
		{[]int32{PUSH, 2, PUSH, 8, SRL, HALT}, 2},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[vm.sp] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[vm.sp])
		}
	}
}
