package integration

import "testing"

func TestErrorReturningCallsCompareWithNilAndEvaluateOnce(t *testing.T) {
	runRecordFixture(t, `package main
fn failed(calls *i32) error {
  *calls = *calls + 1
  return error("expected failure")
}
fn succeeded() error {
  return nil
}
fn main() i32 {
  calls i32 := 0
  if failed(&calls) == nil {
    return 1
  }
  if nil == failed(&calls) {
    return 2
  }
  if succeeded() != nil or nil != succeeded() {
    return 3
  }
  if calls != 2 {
    return 4
  }
  return 0
}
`)
}
