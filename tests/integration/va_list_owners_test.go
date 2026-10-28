package integration

import "testing"

func TestVaListOwnersInLoopInitializersAreConstructedInPlace(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
fn inspect(marker i32, args @c ...) i32 {
active := true
for owner := libc::va_start(args); active; active = false {
defer libc::va_end(&owner)
if libc::va_arg<libc::c_int>(&owner) != 42 {
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
use "std::libc"
fn inspect(marker i32, args @c ...) i32 {
owner := libc::va_start(args)
defer libc::va_end(&owner)
flags bool[2] := .{false, true}
if false and flags[libc::va_arg<u64>(&owner)] {
return 1
}
if libc::va_arg<u64>(&owner) != 1 {
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
use "std::libc"
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_pointer_drive"]
fn drive(callback Inspect) i32
}
fn inspect(args libc::c_va_list) i32 {
owner := libc::va_copy(args)
defer libc::va_end(&owner)
value := *libc::va_arg<*i32>(&owner)
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
use "std::libc"
fn first(marker i32, args @c ...) libc::c_int {
owner := libc::va_start(args)
value := libc::va_arg<libc::c_int>(&owner)
libc::va_end(&owner)
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
use "std::libc"
struct Pair {
first libc::c_int
second libc::c_int
}
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_order_drive"]
fn drive(callback Inspect) i32
}
fn inspect(args libc::c_va_list) i32 {
owner := libc::va_copy(args)
defer libc::va_end(&owner)
pair := Pair.{second: libc::va_arg<libc::c_int>(&owner), first: libc::va_arg<libc::c_int>(&owner)}
if pair.second != 42 or pair.first != 7 or libc::va_arg<libc::c_int>(&owner) != 99 {
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
use "std::libc"
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_order_drive"]
fn drive(callback Inspect) i32
}
fn pair(first libc::c_int, second libc::c_int) libc::c_int {
return first * 10 + second
}
fn inspect(args libc::c_va_list) i32 {
owner := libc::va_copy(args)
defer libc::va_end(&owner)
if pair(libc::va_arg<libc::c_int>(&owner), libc::va_arg<libc::c_int>(&owner)) != 427 {
return 0
}
if libc::va_arg<libc::c_int>(&owner) != 99 {
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
use "std::libc"
fn count(marker i32, args @c ...) i32 {
owner := libc::va_start(args)
defer libc::va_end(&owner)
if false and libc::va_arg<libc::c_int>(&owner) == 99 {
return 9
}
total i32 := 0
while libc::va_arg<libc::c_int>(&owner) != 0 {
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
use "std::libc"
enum Mode {
Read
Skip
}
fn inspect(mode Mode, marker i32, args @c ...) i32 {
return match mode {
Mode::Read => {
owner := libc::va_start(args)
defer libc::va_end(&owner)
yield libc::va_arg<libc::c_int>(&owner)
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
use "std::libc"
fn first<T>(marker i32, args @c ...) T {
owner := libc::va_start(args)
defer libc::va_end(&owner)
return libc::va_arg<T>(&owner)
}
fn main() i32 {
if first<libc::c_int>(0, 42) != 42 {
return 1
}
if first<libc::c_double>(0, 1.5) != 1.5 {
return 2
}
return 0
}
`)
}

func TestVaStartInitializesNativeStorageInsideCVariadicFunctions(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
type Consume = fn(marker i32, args @c ...) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_owned_drive"]
fn drive(callback Consume) i32
}
@[cc="c"]
fn consume(marker i32, args @c ...) i32 {
owner := libc::va_start(args)
defer libc::va_end(&owner)
number := libc::va_arg<libc::c_int>(&owner)
fraction := libc::va_arg<libc::c_double>(&owner)
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
use "std::libc"
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
}
fn inspect(args libc::c_va_list) i32 {
owner := libc::va_copy(args)
defer libc::va_end(&owner)
number := libc::va_arg<libc::c_int>(&owner)
fraction := libc::va_arg<libc::c_double>(&owner)
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
use "std::libc"
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
@[link_name="fixture_va_check"]
fn check(args libc::c_va_list) i32
}
fn inspect(args libc::c_va_list) i32 {
first := libc::va_copy(args)
defer libc::va_end(&first)
second := libc::va_copy(libc::va_borrow(&first))
defer libc::va_end(&second)
if check(libc::va_borrow(&first)) != 1 {
return 0
}
return check(libc::va_borrow(&second))
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}
