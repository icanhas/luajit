package luajit

import "testing"
import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func TestPushpop(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate failed")
	}
	r := rand.New(rand.NewSource(1))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(6 * 1000) {
		t.Fatal("not enough slots in stack")
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
		s.Pop(3) // pop the nil as well
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
		t.Fatal("Newstate failed")
	}
	r := rand.New(rand.NewSource(2))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(5 * 1000) {
		t.Fatal("not enough slots in stack")
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

func TestLoad(t *testing.T) {
	txt := `
		function f(x)
			return math.sqrt(x)
		end
		testx = f(400)
		testy = f(1)
		testz = f(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	r := bufio.NewReader(strings.NewReader(txt))
	if r == nil {
		t.Fatal("NewReader returned nil")
	}
	s.Openlibs()
	if err := s.Load(r, "TestLoad"); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
	s.Close()
}

func TestLoadstring(t *testing.T) {
	txt := `
		function f(x)
			return math.sqrt(x)
		end
		testx = f(400)
		testy = f(1)
		testz = f(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	s.Openlibs()
	if err := s.Loadstring(txt); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
	s.Close()
}

func TestRegister(t *testing.T) {
	txt := `
		testx = mysqrt(400)
		testy = mysqrt(1)
		testz = mysqrt(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	s.Openlibs()
	s.Register(func(s *State) int {
		n := s.Tonumber(-1)
		s.Pushnumber(math.Sqrt(n))
		return 1
	}, "mysqrt")
	if err := s.Loadstring(txt); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		errdetail := s.Tostring(-1)
		s.Close()
		t.Fatalf("%s -- %s", err.Error(), errdetail)
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
	s.Close()
}
