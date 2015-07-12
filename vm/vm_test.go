package vm

import "testing"

func TestAllOps(t *testing.T) {
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
