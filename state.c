#include <lua.h>
#include <stddef.h>
#include <stdlib.h>
#include "_cgo_export.h"

typedef struct Readbuf	Readbuf;
struct Readbuf {
	void*	reader;
	char*	buf;
	size_t	bufsz;
};

// a lua_Alloc
static void*
lalloc(void *ud, void *p, size_t osize, size_t nsize)
{
	if(nsize != 0)
		return realloc(p, nsize);
	free(p);
	return NULL;
}

// a lua_Reader
static const char*
readchunk(lua_State *l, void *data, size_t *size)
{
	Readbuf *rb;
	size_t sz;
	
	rb = data;
	if((sz = goreadchunk(rb->reader, rb->buf, rb->bufsz)) < 1){
		free(rb->buf);
		free(rb);
		return NULL;
	}
	*size = sz;
	return rb->buf;
}

// a lua_Writer
static int
writechunk(lua_State *l, const void *p, size_t sz, void *ud)
{
	if(gowritechunk(ud, (void*)p, sz) != sz)
		return 1;
	return 0;
}

lua_State*
newstate(void)
{
	return lua_newstate(lalloc, NULL);
}

int
load(lua_State *l, void *reader, size_t bufsz, const char *chunkname)
{
	char *buf;
	Readbuf *rb;
	
	buf = malloc(bufsz);		// both allocs are freed by readchunk
	if(buf == NULL)
		return LUA_ERRMEM;
	rb = malloc(sizeof *rb);
	if(rb == NULL){
		free(buf);
		return LUA_ERRMEM;
	}
	rb->reader = reader;
	rb->buf = buf;
	rb->bufsz = bufsz;
	return lua_load(l, readchunk, rb, chunkname);
}

int
dump(lua_State *l, void *ud)
{
	return lua_dump(l, writechunk, ud);
}
