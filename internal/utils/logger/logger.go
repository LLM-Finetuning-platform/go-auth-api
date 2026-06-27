package logger

import (
	"log/slog"
	"os"
)

type Slogger struct {
	logger slog.Logger
}

func (s *Slogger) GetLogger () (*slog.Logger){
	return  &s.logger
}

func (s* Slogger) NewLogger () (*slog.Logger, error){
	s.logger = *slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &s.logger, nil
}

func GetAuthSlogger() (*Slogger) {
	slogger := &Slogger{logger: *slog.New(slog.NewTextHandler(os.Stdout, nil))}
	return  slogger
}

var AuthSlogger *Slogger = GetAuthSlogger()