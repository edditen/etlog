package common

const (
	// CurrentLine means log filename and current line number.
	CurrentLine = 3

	flagFunc = 1 << 16

	// CurrentLineF means log filename, current line number and function name
	CurrentLineF = CurrentLine | flagFunc

	// NoneSource logs no source info
	NoneSource = 0

	// MaxSkip is the max depth of runtime caller, this is intend to prevent misusing `WithSource`
	MaxSkip = 16
)
