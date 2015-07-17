package vm

import "testing"

type TestCase struct {
	Code       []int32
	TopOfStack int32
}

func TestAddLoop(t *testing.T) {
	tests := []TestCase{
		{[]int32{PUSH, 0, PUSH, 1, ADD, HALT}, 1},
		{[]int32{PUSH, 2, PUSH, 3, ADD, JMPNE, 17, 2, HALT}, 17},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[0] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[0])
		}
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

		if vm.stack[0] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[0])
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

		if vm.stack[0] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[0])
		}
	}
}
