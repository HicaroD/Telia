package parser

import (
	"testing"
)

func TestRangeLoopBindingLists(t *testing.T) {
	for _, source := range []string{"for key, value in scores {}", "for element in names {}", "for a, b, c in scores {}"} {
		loop, err := ParseForLoopFrom(source, "loop.t")
		if err != nil {
			t.Fatal(err)
		}
		want := 1
		for _, char := range source {
			if char == ',' {
				want++
			}
		}
		if len(loop.Bindings) != want || loop.Range == nil {
			t.Fatalf("lost range bindings: %+v", loop)
		}
	}
}

func TestRangeBindingsDoNotConsumeThreeClauseInitializer(t *testing.T) {
	loop, err := ParseForLoopFrom("for a, b := pair(); a < b; a = a + 1 {}", "loop.t")
	if err != nil {
		t.Fatal(err)
	}
	if loop.Range != nil || loop.Init == nil {
		t.Fatal("three-clause loop parsed as range")
	}
}
