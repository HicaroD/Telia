package integration

import "testing"

func TestVaListOwnersInLoopInitializersAreConstructedInPlace(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn inspect(marker i32, args @c ...) i32 {
active := true
for owner := c::va_start(args); active; active = false {
defer c::va_end(&owner)
if c::va_arg<c::c_int>(&owner) != 42 {
return 1
}
}
return 0
}
fn main() i32 {
return inspect(0, 42)
}
`)
}

func TestVaArgArrayIndicesStayInsideShortCircuitBranches(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn inspect(marker i32, args @c ...) i32 {
owner := c::va_start(args)
defer c::va_end(&owner)
flags bool[2] := .{false, true}
if false and flags[c::va_arg<u64>(&owner)] {
return 1
}
if c::va_arg<u64>(&owner) != 1 {
return 2
}
return 0
}
fn main() i32 {
index u64 := 1
sentinel u64 := 0
return inspect(0, index, sentinel)
}
`)
}

func TestDereferencingVaArgPointersAdvancesTheCursorOnlyOnce(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_pointer_drive"]
fn drive(callback Inspect) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
owner := c::va_copy(args)
defer c::va_end(&owner)
value := *c::va_arg<*i32>(&owner)
if value != 42 {
return 0
}
return 1
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}

func TestVaListOwnerCanBeEndedExplicitlyBeforeReturning(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn first(marker i32, args @c ...) c::c_int {
owner := c::va_start(args)
value := c::va_arg<c::c_int>(&owner)
c::va_end(&owner)
return value
}
fn main() i32 {
if first(0, 42) != 42 {
return 1
}
return 0
}
`)
}

func TestVaArgRecordInitializersUseSourceRatherThanCFieldOrder(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
struct Pair {
first c::c_int
second c::c_int
}
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_order_drive"]
fn drive(callback Inspect) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
owner := c::va_copy(args)
defer c::va_end(&owner)
pair := Pair{second: c::va_arg<c::c_int>(&owner), first: c::va_arg<c::c_int>(&owner)}
if pair.second != 42 or pair.first != 7 or c::va_arg<c::c_int>(&owner) != 99 {
return 0
}
return 1
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}

func TestVaArgCallArgumentsAdvanceTheCursorInSourceOrder(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_order_drive"]
fn drive(callback Inspect) i32
}
fn pair(first c::c_int, second c::c_int) c::c_int {
return first * 10 + second
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
owner := c::va_copy(args)
defer c::va_end(&owner)
if pair(c::va_arg<c::c_int>(&owner), c::va_arg<c::c_int>(&owner)) != 427 {
return 0
}
if c::va_arg<c::c_int>(&owner) != 99 {
return 0
}
return 1
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}

func TestVaArgStaysAtShortCircuitAndLoopEvaluationSites(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn count(marker i32, args @c ...) i32 {
owner := c::va_start(args)
defer c::va_end(&owner)
if false and c::va_arg<c::c_int>(&owner) == 99 {
return 9
}
total i32 := 0
while c::va_arg<c::c_int>(&owner) != 0 {
total = total + 1
}
return total
}
fn main() i32 {
if count(0, 10, 20, 0) != 2 {
return 1
}
return 0
}
`)
}

func TestVaListOwnerCleanupRunsOnSelectedMatchYield(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
enum Mode {
Read
Skip
}
fn inspect(mode Mode, marker i32, args @c ...) i32 {
return match mode {
Mode::Read => {
owner := c::va_start(args)
defer c::va_end(&owner)
yield c::va_arg<c::c_int>(&owner)
}
Mode::Skip => 0
}
}
fn main() i32 {
if inspect(Mode::Read, 0, 42) != 42 or inspect(Mode::Skip, 0) != 0 {
return 1
}
return 0
}
`)
}

func TestGenericVaListOwnersKeepIndependentTypesAndNativeParameters(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn first<T>(marker i32, args @c ...) T {
owner := c::va_start(args)
defer c::va_end(&owner)
return c::va_arg<T>(&owner)
}
fn main() i32 {
if first<c::c_int>(0, 42) != 42 {
return 1
}
if first<c::c_double>(0, 1.5) != 1.5 {
return 2
}
return 0
}
`)
}

func TestVaStartInitializesNativeStorageInsideCVariadicFunctions(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Consume = fn(marker i32, args @c ...) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_owned_drive"]
fn drive(callback Consume) i32
}
@[cc="c"]
fn consume(marker i32, args @c ...) i32 {
owner := c::va_start(args)
defer c::va_end(&owner)
number := c::va_arg<c::c_int>(&owner)
fraction := c::va_arg<c::c_double>(&owner)
if number != 42 or fraction != 1.5 {
return 0
}
return 1
}
fn main() i32 {
if drive(consume) != 1 or consume(0, 42, 1.5) != 1 {
return 1
}
return 0
}
`)
}

func TestVaArgReadsPromotedValuesFromOwnedNativeStorage(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
owner := c::va_copy(args)
defer c::va_end(&owner)
number := c::va_arg<c::c_int>(&owner)
fraction := c::va_arg<c::c_double>(&owner)
if number != 42 or fraction != 1.5 {
return 0
}
return 1
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}

func TestVaCopyCreatesIndependentNativeCursors(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
@[link_name="fixture_va_check"]
fn check(args c::c_va_list) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
first := c::va_copy(args)
defer c::va_end(&first)
second := c::va_copy(c::va_borrow(&first))
defer c::va_end(&second)
if check(c::va_borrow(&first)) != 1 {
return 0
}
return check(c::va_borrow(&second))
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}
