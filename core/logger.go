package core

import (
	"io"

	"github.com/rs/zerolog/log"
	"github.com/savsgio/gotils/strconv"
)

type logger struct{}

func (l logger) Write(p []byte) (int, error) {
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[:n-1]
	}
	log.Info().Msg(strconv.B2S(p))
	return n, nil
}

func LoggerWriter() io.Writer {
	return logger{}
}
