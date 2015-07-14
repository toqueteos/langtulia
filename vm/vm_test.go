package vm

import "testing"

type TestCase struct {
	Code       []int
	TopOfStack int
}

func TestAddLoop(t *testing.T) {
	tests := []TestCase{
		{[]int{PUSH, 0, PUSH, 1, ADD, HALT}, 1},
		{[]int{PUSH, 2, PUSH, 3, ADD, JMPNE, 17, 2, HALT}, 17},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[0] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[0])
		}
	}
}

	}

func TestShiftLogical(t *testing.T) {
	tests := []TestCase{
		{[]int{PUSH, 3, PUSH, 1, SLL, HALT}, 8},
		{[]int{PUSH, 2, PUSH, 8, SRL, HALT}, 2},
	}

	for _, tt := range tests {
		vm := New(tt.Code)
		vm.Run()

		if vm.stack[0] != tt.TopOfStack {
			t.Errorf("expected %d, got %d", tt.TopOfStack, vm.stack[0])
		}
	}
}
