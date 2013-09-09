// Package luajit provides an interface to LuaJIT, a just-in-time compiler
// and interpreter for the Lua programming language.
package luajit

/*
#include <lua.h>
*/
import "C"
import "fmt"

const (
	Version    = C.LUA_VERSION
	Versionnum = C.LUA_VERSION_NUM
	Copyright  = C.LUA_COPYRIGHT
	Authors    = C.LUA_AUTHORS
)

const (
	Signature = C.LUA_SIGNATURE // mark for precompiled code (`<esc>Lua')
	Multret   = C.LUA_MULTRET   // option for multiple returns in 'call' functions
	Minstack  = C.LUA_MINSTACK  // minimum Lua stack available to a Go function
)

// Thread status; 0 is OK
const (
	Yield     = C.LUA_YIELD
	Errrun    = C.LUA_ERRRUN
	Errsyntax = C.LUA_ERRSYNTAX
	Errmem    = C.LUA_ERRMEM
	Errerr    = C.LUA_ERRERR
)

var errmsgs map[int]string = map[int]string{
	Errrun:    "run time error",
	Errsyntax: "syntax error",
	Errmem:    "out of memory",
	Errerr:    "error in error handling",
}

func err2msg(errnum int) error {
	if errnum == 0 {
		return nil
	}
	return fmt.Errorf("%s", errmsgs[errnum])
}

// Pseudo-indices. Unless otherwise noted, any function that accepts valid
// indices can also be called with these pseudo-indices, which represent
// some Lua values that are accessible to Go code but which are not in
// the stack. Pseudo-indices are used to access the thread environment,
// the function environment, the registry, and the upvalues of a Go function.
//
// The thread environment (where global variables live) is always at
// pseudo-index Globalsindex. The environment of the running Go function
// is always at pseudo-index Environindex.
//
// To access and change the value of global variables, you can use regular
// table operations over an environment table. For instance, to access the
// value of a global variable, do:
//	s.Getfield(luajit.Globalsindex, varname);
const (
	Registryindex = C.LUA_REGISTRYINDEX
	Environindex  = C.LUA_ENVIRONINDEX // env of running Go function
	Globalsindex  = C.LUA_GLOBALSINDEX // thread env, where globals live
)

// Equivalent of lua_upvalueindex.
func Upvalueindex(i int) int {
	return Globalsindex - i
}

// Basic types
const (
	Tnone          = C.LUA_TNONE
	Tnil           = C.LUA_TNIL
	Tboolean       = C.LUA_TBOOLEAN
	Tlightuserdata = C.LUA_TLIGHTUSERDATA
	Tnumber        = C.LUA_TNUMBER
	Tstring        = C.LUA_TSTRING
	Ttable         = C.LUA_TTABLE
	Tfunction      = C.LUA_TFUNCTION
	Tuserdata      = C.LUA_TUSERDATA
	Tthread        = C.LUA_TTHREAD
)

// Garbage-collection function and options
const (
	// Stops the garbage collector.
	GCstop = C.LUA_GCSTOP
	// Restarts the garbage collector.
	GCrestart = C.LUA_GCRESTART
	// Performs a full garbage-collection cycle.
	GCcollect = C.LUA_GCCOLLECT
	// Returns the current amount of memory (in Kbytes) in use by Lua.
	GCcount = C.LUA_GCCOUNT
	// Returns the remainder of dividing the current amount of bytes of memory
	// in use by Lua by 1024.
	GCcountb = C.LUA_GCCOUNTB
	// Performs an incremental step of garbage collection. The step "size" is
	// controlled by data (larger values mean more steps) in a non-specified
	// way. If you want to control the step size you must experimentally
	// tune the value of data. The function returns 1 if the step finished a
	// garbage-collection cycle.
	GCstep = C.LUA_GCSTEP
	// Sets data as the new value for the pause of the collector. The function
	// returns the previous value of the pause.
	GCsetpause = C.LUA_GCSETPAUSE
	// Sets data as the new value for the step multiplier of the collector. The
	// function returns the previous value of the step multiplier.
	GCsetstepmul = C.LUA_GCSETSTEPMUL
)
