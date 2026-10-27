package diagnostics

// ColorEnabled resolves an already-validated CLI mode. Explicit flags override
// NO_COLOR; auto requires a terminal and the absence of NO_COLOR.
func ColorEnabled(mode string, terminal, noColor bool) bool {
	switch mode {
	case "always":
		return true
	case "never":
		return false
	default:
		return terminal && !noColor
	}
}
