package c

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSipHash13ReferenceVectors(t *testing.T) {
	// Independent vectors: rust-siphash src/tests.rs test_siphash_1_3,
	// key bytes 00..0f, message bytes 00..(length-1), little-endian output.
	source := `#include <stdint.h>
#include <stddef.h>
#include <assert.h>
#include "runtime_hash.h"
int main(void) {
 const uint64_t expected[] = {0xabac0158050fc4dcULL,0xc9f49bf37d57ca93ULL,
 0x82cb9b024dc7d44dULL,0x8bf80ab8e7ddf7fbULL,0xcf75576088d38328ULL,
 0xdef9d52f49533b67ULL,0xc50d2b50c59f22a7ULL,0xd3927d989bb11140ULL,
 0x369095118d299a8eULL,0x25a48eb36c063de4ULL};
 unsigned char bytes[65]; for(unsigned i=0;i<65;++i) bytes[i]=(unsigned char)(i-1);
 for(unsigned n=0;n<10;++n) {
  _TeliaSip s = _telia_sip_init(0x0706050403020100ULL,0x0f0e0d0c0b0a0908ULL);
  _telia_sip_write(&s,bytes+1,n); assert(_telia_sip_finish(s)==expected[n]);
  s = _telia_sip_init(0x0706050403020100ULL,0x0f0e0d0c0b0a0908ULL);
  for(unsigned i=0;i<n;++i) _telia_sip_write(&s,bytes+1+i,1);
  assert(_telia_sip_finish(s)==expected[n]);
 }
 _TeliaSip s = _telia_sip_init(0x0706050403020100ULL,0x0f0e0d0c0b0a0908ULL);
 _telia_sip_write(&s,bytes+1,32); assert(_telia_sip_finish(s)==0x81157b6c16a7b60dULL);
 const unsigned lengths[] = {15,16,17,31,33,63};
 const uint64_t boundaries[] = {0xd320d86d2a519956ULL,0xcc4fdd1a7d908b66ULL,
 0x9cf2689063dbd80cULL,0x2370dd1f8c21d1bcULL,0x4d54b9e57a8ff9bfULL,0x9d199062b7bbb3a8ULL};
 for(unsigned n=0;n<6;++n) {
  for(unsigned split=0;split<=lengths[n];++split) {
   s = _telia_sip_init(0x0706050403020100ULL,0x0f0e0d0c0b0a0908ULL);
   _telia_sip_write(&s,bytes+1,split);
   _telia_sip_write(&s,bytes+1+split,lengths[n]-split);
   assert(_telia_sip_finish(s)==boundaries[n]);
  }
 }
 return 0;
}`
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) { runHashC(t, source, opt) })
	}
	if runtime.GOOS != "windows" {
		t.Run(
			"ubsan",
			func(t *testing.T) { runHashC(t, source, "-O3", "-fsanitize=undefined", "-fno-sanitize-recover=all") },
		)
	}
}

func runHashC(t *testing.T, source string, flags ...string) {
	t.Helper()
	dir := t.TempDir()
	input, output := filepath.Join(dir, "test.c"), filepath.Join(dir, "test")
	if runtime.GOOS == "windows" {
		output += ".exe"
	}
	if err := os.WriteFile(input, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	include, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	args := append([]string{"-std=gnu11", "-I", include, input, "-o", output}, flags...)
	if os.Getenv("TELIA_TEST_HASH_AMD64") == "1" && runtime.GOOS == "darwin" {
		args = append(args, "-arch", "x86_64")
	}
	if out, err := exec.Command("cc", args...).CombinedOutput(); err != nil {
		t.Fatalf("compile: %v: %s", err, out)
	}
	if out, err := exec.Command(output).CombinedOutput(); err != nil {
		t.Fatalf("run: %v: %s", err, out)
	}
}

func TestHashSeedNative(t *testing.T) {
	source := `#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>
#include <errno.h>
#include <limits.h>
#include <assert.h>
#include <unistd.h>
#ifdef _WIN32
#include <windows.h>
#endif
typedef struct { bool present; } _Error;
#define _TELIA_ERROR(kind,message) ((_Error){true})
#include "runtime_seed.h"
int main(void) {
 uint64_t a,b;
 _Error e=_telia_hash_seed(&a,&b);
 assert(!e.present); // No statistical assertion about secret output bytes.
 return 0;
}`
	// Windows does not provide unistd.h consistently outside MinGW; the target
	// compiler for this test, like the generated-C compiler, is MinGW-w64.
	runHashC(t, source, "-O0")
	runHashC(t, source, "-O3")
}

// Execute each adapter's real control flow with deterministic OS calls. This
// complements, but is not a substitute for, native runs on the CI matrix.
func TestHashSeedAdapters(t *testing.T) {
	const common = `#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>
#include <errno.h>
#include <limits.h>
#include <assert.h>
#undef __linux__
#undef __APPLE__
#undef _WIN32
typedef struct { bool present; int kind; } _Error;
#define _TELIA_ERROR(kind,message) ((_Error){true, 42})
static int mode, calls;
static void fill(unsigned char *p,size_t n) { for(size_t i=0;i<n;++i)p[i]=(unsigned char)i; }
`
	cases := []struct{ name, stub string }{
		{"linux", `#define __linux__ 1
#include <sys/types.h>
static ssize_t getrandom(void *p,size_t n,unsigned flags) {
 (void)flags; ++calls;
 if(mode==1){errno=EIO;return -1;}
 if(mode==2)return 0;
 if(calls==1){errno=EINTR;return -1;}
 size_t count=n>3?3:n; fill(p,count); return (ssize_t)count;
}
`},
		{"darwin", `#define __APPLE__ 1
static int getentropy(void *p,size_t n) {
 ++calls;
 if(mode){errno=EIO;return -1;}
 if(calls==1){errno=EINTR;return -1;}
 fill(p,n); return 0;
}
`},
		{"windows", `#define _WIN32 1
#include <wchar.h>
#define WINAPI
#define BCRYPT_USE_SYSTEM_PREFERRED_RNG 2
#define LOAD_LIBRARY_SEARCH_SYSTEM32 0x800
typedef void *HMODULE;
typedef void *BCRYPT_ALG_HANDLE;
typedef unsigned char *PUCHAR;
typedef unsigned long ULONG;
typedef long NTSTATUS;
static NTSTATUS generate(BCRYPT_ALG_HANDLE h,PUCHAR p,ULONG n,ULONG flags) {
 (void)h; assert(flags==2); ++calls; fill(p,n); return mode==1?-1:0;
}
static HMODULE LoadLibraryExW(const wchar_t *name,void *file,unsigned flags){(void)name;assert(!file && flags==0x800);return mode==2?NULL:(HMODULE)1;}
static void *GetProcAddress(HMODULE m,const char *name){(void)m;(void)name;return mode==3?NULL:(void *)generate;}
static void FreeLibrary(HMODULE m){assert(m);}
`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.Mkdir(filepath.Join(dir, "sys"), 0700); err != nil {
				t.Fatal(err)
			}
			for _, header := range []string{"sys/random.h", "bcrypt.h"} {
				if err := os.WriteFile(filepath.Join(dir, header), nil, 0600); err != nil {
					t.Fatal(err)
				}
			}
			source := common + tc.stub + `#include "runtime_seed.h"
int main(void) {
 uint64_t a=9,b=9;
 _Error e=_telia_hash_seed(&a,&b); assert(!e.present); assert(calls>0);
 for(mode=1;mode<=` + map[string]string{"linux": "2", "darwin": "1", "windows": "3"}[tc.name] + `; ++mode) {
  a=b=9; calls=0; e=_telia_hash_seed(&a,&b);
  assert(e.present && e.kind==42); assert(a==0 && b==0);
 }
 return 0;
}`
			runHashC(t, source, "-O0", "-I", dir)
			runHashC(t, source, "-O3", "-I", dir)
		})
	}
}
