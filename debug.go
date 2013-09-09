package luajit

/*
#cgo LDFLAGS: -lluajit
#cgo linux LDFLAGS: -lm -ldl

#include <lua.h>
#include <stddef.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// A Debug is used to carry different pieces of information about an active
// function. Getstack fills only the private part of this structure, for
// later use. To fill the other fields of Debug with useful information,
// call Getinfo.
type Debug struct {
	// A reasonable name for the given function. Because functions in
	// Lua are first-class values, they do not have a fixed name: some
	// functions can be the value of multiple global variables, while
	// others can be stored only in a table field. The State.Getinfo
	// function checks how the function was called to find a suitable
	// name. If it cannot find a name, then name is an empty string.
	Name string
	// Explains the name field. The value of namewhat can be "global",
	// "local", "method", "field", "upvalue", or "" (the empty string),
	// according to how the function was called. (Lua uses the empty
	// string when no other option seems to apply.)
	Namewhat string
	// The string "Lua" if the function is a Lua function, "Go" if it
	// is a Go function, "main" if it is the main part of a chunk, and
	// "tail" if it was a function that did a tail call. In the latter
	// case, Lua has no other information about the function.
	What string
	// If the function was defined in a string, then Source is that
	// string. If the function was defined in a file, then source starts
	// with a '@' followed by the file name.
	Source string
	// "Printable" version of Source, for use in error messages.
	Shortsrc string
	// The current line where the given function is executing. When no
	// line information is available, currentline is set to -1.
	Currentline int
	// The number of upvalues of the function.
	Nups int
	// The line number where the definition of the function starts.
	Linedefined int
	// The line number where the definition of the function ends.
	Lastlinedefined int

	l *C.lua_State
	d C.lua_Debug
}

func Newdebug(s *State) *Debug {
	d := Debug{}
	d.l = s.l
	return &d
}

// Sync a Debug with its C struct.
func (ar *Debug) update() {
	if ar.d.name != nil {
		ar.Name = C.GoString(ar.d.name)
	}
	if ar.d.namewhat != nil {
		ar.Namewhat = C.GoString(ar.d.namewhat)
	}
	if ar.d.what != nil {
		ar.What = C.GoString(ar.d.what)
	}
	if ar.d.source != nil {
		ar.Source = C.GoString(ar.d.source)
	}
	ar.Shortsrc = C.GoString((*C.char)(&ar.d.short_src[0]))
	ar.Currentline = int(ar.d.currentline)
	ar.Nups = int(ar.d.nups)
	ar.Linedefined = int(ar.d.linedefined)
	ar.Lastlinedefined = int(ar.d.lastlinedefined)
}

// Returns information about a specific function or function invocation.
//
// To get information about a function invocation, the parameter ar must be
// a valid activation record that was filled by a previous call to Getstack
// or given as argument to a hook.
//
// To get information about a function you push it onto the stack and start
// the what string with the character '>'. (In that case, Getinfo pops the
// function in the top of the stack.) For instance, to know in which line
// a function f was defined, you can write the following code:
//
// 	d := luajit.Newdebug(s)
// 	s.Getfield(luajit.Globalsindex, "f")  // get global 'f'
// 	d.Getinfo(">S")
// 	fmt.Printf("%d\n", d.Linedefined);
//
// Each character in the string what selects some fields of the structure
// ar to be filled or a value to be pushed on the stack:
//
// 	'n'	fills in the field Name and Namewhat
// 	'S'	fills in the fields Source, Shortsrc, Linedefined,
// 		Lastlinedefined, and What
// 	'l'	fills in the field Currentline
// 	'u'	fills in the field Nups
// 	'f'	pushes onto the stack the function that is running at the
// 		given level
// 	'L'	pushes onto the stack a table whose indices are the numbers of
// 		the lines that are valid on the function. (A valid line is a line
// 		with some associated code, that is, a line where you can put a break
// 		point. Invalid lines include empty lines and comments.)
//
func (d *Debug) Getinfo(what string) error {
	cs := C.CString(what)
	defer C.free(unsafe.Pointer(cs))
	if int(C.lua_getinfo(d.l, cs, &d.d)) == 0 {
		return fmt.Errorf("The significant owl hoots in the night.")
	}
	d.update()
	return nil
}

// Gets information about a local variable of a given activation record. The
// parameter ar must be a valid activation record that was filled by a
// previous call to Getstack or given as argument to a hook. The index
// n selects which local variable to inspect (1 is the first parameter
// or active local variable, and so on, until the last active local
// variable). Getlocal pushes the variable's value onto the stack and
// returns its name.
//
// Variable names starting with '(' (open parentheses) represent internal
// variables (loop control variables, temporaries, and Go function locals).
//
// Returns an empty string (and pushes nothing) when the index is greater
// than the number of active local variables.
func (d *Debug) Getlocal(n int) string {
	cs := C.lua_getlocal(d.l, &d.d, C.int(n))
	if cs == nil {
		return ""
	}
	d.update()
	return C.GoString(cs)
}

// Gets information about the interpreter runtime stack.
//
// This function fills parts of a Debug structure with an identification of
// the activation record of the function executing at a given level. Level
// 0 is the current running function, whereas level n+1 is the function that
// has called level n. When there are no errors, Getstack returns nil; when
// called with a level greater than the stack depth, it returns the error.
func (d *Debug) Getstack(level int) error {
	if int(C.lua_getstack(d.l, C.int(level), &d.d)) == 0 {
		return fmt.Errorf("stack depth exceeded")
	}
	d.update()
	return nil
}

// int lua_sethook(lua_State *L, lua_Hook f, int mask, int count);

// Sets the value of a local variable of a given activation
// record. Parameters d and n are as in Getlocal. Setlocal assigns the
// value at the top of the stack to the variable and returns its name. It
// also pops the value from the stack.
//
// Returns an empty string (and pops nothing) when the index is greater
// than the number of active local variables.
func (d *Debug) Setlocal(n int) string {
	return C.GoString(C.lua_setlocal(d.l, &d.d, C.int(n)))
}
