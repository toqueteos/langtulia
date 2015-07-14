package vm

import "testing"

func TestAddLoop(t *testing.T) {
	code := []int{
		PUSH, 2,
		PUSH, 3,
		ADD,
		JMPNE, 11, 2,
		PRINT,
		HALT,
	}

	vm := New(code)
	vm.Run()

	if vm.stack[0] != 11 {
		t.Fatalf("expected 11, got %d", vm.stack[0])
	}
}

func TestShiftLogical(t *testing.T) {
	type TestCase struct {
		Code       []int
		TopOfStack int
	}

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
