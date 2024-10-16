package scorereader

import (
	"io"
	"os"
)

type IScoreReader interface {
	HandleScores(scores io.Reader)
}

type ConsoleScoreReader struct{}

func NewConsoleScoreReader() *ConsoleScoreReader {
	return &ConsoleScoreReader{}
}

func (r ConsoleScoreReader) HandleScores(scores io.Reader) {
	io.Copy(os.Stdout, scores)
}

func ReadScores(reader IScoreReader, scores io.Reader) {
	reader.HandleScores(scores)
}
