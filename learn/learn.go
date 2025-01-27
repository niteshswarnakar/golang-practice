package learn

type Letter interface {
	Count() int
}

type Number struct {
	Number string
}

func (n *Number) Count() int {
	return len(n.Number)
}

func NewNumber(num string) Letter {
	return &Number{Number: num}
}
