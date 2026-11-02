package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

// EqualityEligibility classifies resolved concrete types. Key mode uses the
// same traversal, but excludes otherwise comparable function values and owners.
// Pointers terminate traversal: only their identity participates.
func EqualityEligibility(ty *ExprType, key bool) error {
	return equalityEligibility(ty, key, "value")
}

func equalityEligibility(ty *ExprType, key bool, path string) error {
	reject := func(reason string) error { return fmt.Errorf("%s (%s): %s", path, ty.T, reason) }
	switch ty.Kind {
	case EXPR_TYPE_ALIAS:
		return equalityEligibility(ty.T.(*TypeAlias).Type, key, path)
	case EXPR_TYPE_POINTER:
		pointee := ty.T.(*PointerType).Type
		if pointee != nil && pointee.Kind == EXPR_TYPE_BASIC &&
			pointee.T.(*BasicType).Kind == token.UNTYPED_NULLPTR {
			return reject("bare nil is not a concrete key or component type")
		}
		return nil
	case EXPR_TYPE_ERROR_KIND:
		return nil
	case EXPR_TYPE_FUNCTION:
		if key {
			return reject("function values are not eligible keys")
		}
		return nil
	case EXPR_TYPE_BASIC:
		b := ty.T.(*BasicType)
		if ty.IsNumeric() || ty.IsBoolean() || b.Kind == token.RAWPTR_TYPE ||
			b.Kind == token.STRING_TYPE ||
			b.Kind == token.CSTRING_TYPE {
			return nil
		}
		if ty.IsError() {
			return reject(
				"general errors have no whole-value equality; compare presence against nil or compare kind",
			)
		}
		return reject("whole-value equality is not supported")
	case EXPR_TYPE_SLICE:
		return reject("slices have no whole-value equality and are not eligible keys")
	case EXPR_TYPE_OPAQUE:
		return reject("opaque values have no whole-value equality; use pointer identity")
	case EXPR_TYPE_C_VA_LIST, EXPR_TYPE_VA_LIST_OWNER:
		return reject("native argument-list lifecycle types have no whole-value equality")
	case EXPR_TYPE_ARRAY:
		if ty.T.(*SequenceType).Foreign {
			return reject(
				"foreign inline arrays have no whole-value equality; compare indexed elements",
			)
		}
		return equalityEligibility(ty.T.(*SequenceType).Element, key, path+"[element]")
	case EXPR_TYPE_TUPLE:
		for i, t := range ty.T.(*TupleType).Types {
			if err := equalityEligibility(t, key, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
		return nil
	case EXPR_TYPE_STRUCT:
		d := ty.T.(*StructType).Decl
		if key && standardOwningHandle(d) {
			return reject("standard-library owning handle is not an eligible key")
		}
		if d.Foreign != nil {
			return reject("foreign records and unions have no whole-value equality")
		}
		for _, f := range d.Fields {
			if err := equalityEligibility(f.Type, key, path+"."+f.Name.Name()); err != nil {
				return err
			}
		}
		return nil
	case EXPR_TYPE_ENUM:
		for _, v := range ty.T.(*EnumType).Decl.Variants {
			for _, f := range v.Fields {
				if err := equalityEligibility(f.Type, key, path+"::"+v.Name.Name()+"."+f.Name.Name()); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return reject("whole-value equality is not supported")
	}
}

// Match the declaration's actual std package, not its source spelling at a use
// site or an unrelated user package/name. Instances retain template identity.
func standardOwningHandle(d *StructDecl) bool {
	if d.Template != nil {
		d = d.Template
	}
	if d.Owner == nil {
		return false
	}
	switch d.Owner.StandardPath + "::" + d.Name.Name() {
	case "strings::Owned", "strings::COwned", "memory::Arena", "lists::List",
		"fs::File", "process::Args", "process::Child", "terminal::Handle", "terminal::SavedMode",
		"maps::Map", "sets::Set":
		return true
	}
	return false
}
