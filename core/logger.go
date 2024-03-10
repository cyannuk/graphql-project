package core

import (
	"io"

	"github.com/rs/zerolog/log"
)

type logger struct{}

func (l logger) Write(p []byte) (int, error) {
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[:n-1]
	}
	log.Info().Msg(string(p))
	return n, nil
}

func LoggerWriter() io.Writer {
	return logger{}
}
