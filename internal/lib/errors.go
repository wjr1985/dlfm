package lib

type Error struct {
	text string

	isEmptyField bool
}

// ==================== Functions ====================

func NewError(text string) Error {
	return Error{
		text: text,
	}
}

func ErrEmptyField(field string) Error {
	return Error{
		text:         field + " is empty",
		isEmptyField: true,
	}
}

func (e Error) IsEmptyField() bool {
	return e.isEmptyField
}

func (e Error) Error() string {
	return e.text
}

// ==================== Functions ====================

var (
	ErrNilLoggerPtr        = NewError("nil log.Logger pointer")
	ErrNilDGoSession       = NewError("nil discordgo session")
	ErrDiscordDisconnected = NewError("discord disconnected")
	ErrIncorrectNowPlaying = NewError("incorrect ctrack.NowPlaying")
)
