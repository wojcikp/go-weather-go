package scorereader

import (
	"io"
	"os"
)

type IScoreReader interface {
	ReadScores(scores io.Reader)
}

type ConsoleScoreReader struct{}

func NewConsoleScoreReader() *ConsoleScoreReader {
	return &ConsoleScoreReader{}
}

func (r ConsoleScoreReader) ReadScores(scores io.Reader) {
	io.Copy(os.Stdout, scores)
}
