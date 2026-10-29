package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

// CVaListType is a borrowed native C va_list parameter. Its ABI is supplied by
// <stdarg.h>, including targets where a va_list parameter undergoes adjustment.
type CVaListType struct{}

func (*CVaListType) String() string { return "c::c_va_list" }

// VaListOwnerType owns local native storage initialized in place, never copied.
type VaListOwnerType struct{}

func (*VaListOwnerType) String() string { return "c::VaList" }

type VaIntrinsic struct{ Name *token.Token }
