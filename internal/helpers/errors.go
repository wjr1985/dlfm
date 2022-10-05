package helpers

type Error struct {
	text string
}

func (e Error) Error() string { return e.text }
