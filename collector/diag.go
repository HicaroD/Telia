package collector

type Diag struct {
	Message string
}

func (diag Diag) String() string {
	return diag.Message
}
