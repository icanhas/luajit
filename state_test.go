package luajit

import "testing"
import (
	"fmt"
	"math/rand"
)

func TestPushpop(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Error("Newstate failed")
	}
	r := rand.New(rand.NewSource(1))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(6 * 1000) {
		t.Error("not enough slots in stack")
	}
	for i = 0; i < len(vals); i++ {
		vals[i] = r.Float64() * 1000.0
		s.Pushnumber(vals[i])
		s.Pushstring(fmt.Sprintf("!!%f", vals[i]))
		s.Pushinteger(int(vals[i]))
		s.Pushnil()
		s.Pushboolean(true)
		s.Pushboolean(false)
	}
	for i--; i >= 0; i-- {
		n := vals[i]
		if s.Toboolean(-1) {
			t.Errorf("expected false, got true")
		}
		if !s.Toboolean(-2) {
			t.Errorf("expected true, got true")
		}
		s.Pop(3)	// pop the nil as well
		if nn := s.Tointeger(-1); nn != int(n) {
			t.Errorf("expected int %d, got %d", int(n), nn)
		}
		ns := fmt.Sprintf("!!%f", n)
		if str := s.Tostring(-2); str != ns {
			t.Errorf("expected string %s, got %s", ns, str)
		}
		if f := s.Tonumber(-3); f != n {
			t.Errorf("expected float64 %f, got %f", n, f)
		}
		s.Pop(3)
	}
	s.Close()
}

func TestStacktypes(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Error("Newstate failed")
	}
	r := rand.New(rand.NewSource(2))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(5 * 1000) {
		t.Error("not enough slots in stack")
	}
	for i = 0; i < len(vals); i++ {
		vals[i] = r.Float64() * 1000.0
		s.Pushnumber(vals[i])
		s.Pushstring(fmt.Sprintf("!!%f", vals[i]))
		s.Pushinteger(int(vals[i]))
		s.Pushnil()
		s.Pushboolean(true)
	}
	for i--; i >= 0; i-- {
		if !s.Isboolean(-1) {
			t.Errorf("expected boolean")
		}
		if !s.Isnil(-2) {
			t.Errorf("expected nil")
		}
		if !s.Isnumber(-3) {
			t.Errorf("expected number (from int)")
		}
		if !s.Isstring(-4) {
			t.Errorf("expected string")
		}
		if !s.Isnumber(-5) {
			t.Errorf("expected number")
		}
		s.Pop(5)
	}
	s.Close()
}
