package pdf

import (
	"testing"

	"github.com/fjacquet/pdf2md/internal/types"
)

func TestInterpreter_HandleGraphicsState(t *testing.T) {
	in := NewInterpreter(nil, nil)

	// cm
	in.Stack = []Object{Integer(1), Integer(0), Integer(0), Integer(1), Integer(10), Integer(20)}
	if err := in.handleGraphicsState("cm"); err != nil {
		t.Errorf("cm failed: %v", err)
	}
	// CTM should be updated. Initial CTM is Identity.
	// CTM = [1 0 0 1 10 20] * Identity
	if in.State.CTM[4] != 10 || in.State.CTM[5] != 20 {
		t.Errorf("cm did not update CTM correctly: %v", in.State.CTM)
	}

	// q
	if err := in.handleGraphicsState("q"); err != nil {
		t.Errorf("q failed: %v", err)
	}
	if len(in.StateStack) != 1 {
		t.Errorf("q did not push state, stack len: %d", len(in.StateStack))
	}

	// Modify state
	in.State.LineWidth = 5.0

	// Q
	if err := in.handleGraphicsState("Q"); err != nil {
		t.Errorf("Q failed: %v", err)
	}
	if len(in.StateStack) != 0 {
		t.Errorf("Q did not pop state, stack len: %d", len(in.StateStack))
	}
	if in.State.LineWidth != 1.0 { // Default is 1.0
		t.Errorf("Q did not restore state, LineWidth: %f", in.State.LineWidth)
	}

	// w
	in.Stack = []Object{Real(2.5)}
	if err := in.handleGraphicsState("w"); err != nil {
		t.Errorf("w failed: %v", err)
	}
	if in.State.LineWidth != 2.5 {
		t.Errorf("w did not set LineWidth: %f", in.State.LineWidth)
	}
}

func TestInterpreter_HandlePathConstruction(t *testing.T) {
	in := NewInterpreter(nil, nil)

	// m
	in.Stack = []Object{Integer(10), Integer(10)}
	if err := in.handlePathConstruction("m"); err != nil {
		t.Errorf("m failed: %v", err)
	}
	if len(in.State.CurrentPath) != 1 {
		t.Errorf("m did not add path op")
	}
	if in.State.CurrentPath[0].Type != types.PathOpMoveTo {
		t.Errorf("Expected MoveTo, got %v", in.State.CurrentPath[0].Type)
	}

	// l
	in.Stack = []Object{Integer(20), Integer(20)}
	if err := in.handlePathConstruction("l"); err != nil {
		t.Errorf("l failed: %v", err)
	}
	if len(in.State.CurrentPath) != 2 {
		t.Errorf("l did not add path op")
	}

	// c
	in.Stack = []Object{Integer(30), Integer(30), Integer(40), Integer(40), Integer(50), Integer(50)}
	if err := in.handlePathConstruction("c"); err != nil {
		t.Errorf("c failed: %v", err)
	}
	if len(in.State.CurrentPath) != 3 {
		t.Errorf("c did not add path op")
	}

	// re
	in.Stack = []Object{Integer(0), Integer(0), Integer(100), Integer(100)}
	if err := in.handlePathConstruction("re"); err != nil {
		t.Errorf("re failed: %v", err)
	}
	// re adds 4 lines + close = 5 ops?
	// Implementation: MoveTo, LineTo, LineTo, LineTo, Close. Total 5.
	// Previous 3 ops + 5 = 8.
	if len(in.State.CurrentPath) != 8 {
		t.Errorf("re did not add correct ops, count: %d", len(in.State.CurrentPath))
	}

	// h
	if err := in.handlePathConstruction("h"); err != nil {
		t.Errorf("h failed: %v", err)
	}
	if in.State.CurrentPath[len(in.State.CurrentPath)-1].Type != types.PathOpClose {
		t.Errorf("h did not add Close op")
	}
}

func TestInterpreter_HandlePathPainting(t *testing.T) {
	in := NewInterpreter(nil, nil)

	// Create a path big enough to be captured
	in.State.CurrentPath = []types.PathOperation{
		{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 100}}},
	}

	// S
	if err := in.handlePathPainting("S"); err != nil {
		t.Errorf("S failed: %v", err)
	}
	if len(in.Graphics) != 1 {
		t.Errorf("S did not capture graphic")
	}
	if !in.Graphics[0].IsStroked || in.Graphics[0].IsFilled {
		t.Errorf("S captured wrong style: %+v", in.Graphics[0])
	}

	// Reset path
	in.State.CurrentPath = []types.PathOperation{
		{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 100}}},
	}

	// f
	if err := in.handlePathPainting("f"); err != nil {
		t.Errorf("f failed: %v", err)
	}
	if len(in.Graphics) != 2 {
		t.Errorf("f did not capture graphic")
	}
	if in.Graphics[1].IsStroked || !in.Graphics[1].IsFilled {
		t.Errorf("f captured wrong style: %+v", in.Graphics[1])
	}

	// Reset path
	in.State.CurrentPath = []types.PathOperation{
		{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
		{Type: types.PathOpLineTo, Points: []types.Point{{X: 100, Y: 100}}},
	}

	// B
	if err := in.handlePathPainting("B"); err != nil {
		t.Errorf("B failed: %v", err)
	}
	if len(in.Graphics) != 3 {
		t.Errorf("B did not capture graphic")
	}
	if !in.Graphics[2].IsStroked || !in.Graphics[2].IsFilled {
		t.Errorf("B captured wrong style: %+v", in.Graphics[2])
	}

	// n
	in.State.CurrentPath = []types.PathOperation{
		{Type: types.PathOpMoveTo, Points: []types.Point{{X: 0, Y: 0}}},
	}
	if err := in.handlePathPainting("n"); err != nil {
		t.Errorf("n failed: %v", err)
	}
	if len(in.State.CurrentPath) != 0 {
		t.Errorf("n did not clear path")
	}
}
