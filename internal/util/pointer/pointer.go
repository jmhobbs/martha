package pointer


func String(v string) *string {
	s := v
	return &s
}

func Float64(v float64) *float64 {
	f := v
	return &f
}

func Bool(v bool) *bool {
	b := v
	return &b
}