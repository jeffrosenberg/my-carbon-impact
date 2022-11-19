package logging

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	once   sync.Once
	logger *zerolog.Logger

	// Inject via environment variable
	logLevel zerolog.Level = zerolog.DebugLevel // zerolog.Level: Trace = -1, Debug = 0, Info = 1, Error = 3, Disabled = 7

	// TODO: Inject at compile-time?
	CommitID string = "unknown"
)

func LogFunction(function string, start time.Time, msg string, state map[string]interface{}) {
	if logger == nil {
		_ = GetLogger() // Initialize logger but don't worry about the return since it's tied to this module
	}

	logger.Info().
		Str("function", function).
		Str("start_time", start.Format(time.RFC822)).
		Int64("duration_ms", time.Now().Sub(start).Milliseconds()).
		Fields(state).
		Msg(msg)
}

func GetLogger() *zerolog.Logger {
	return GetLoggerWithContext(context.Background())
}

func GetLoggerWithContext(ctx context.Context) *zerolog.Logger {
	if logger == nil {
		once.Do(
			func() {
				zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
				lg := log.Logger

				// Read LogLevel from environment variables
				env := os.Getenv("zerolog_level")
				level, err := strconv.ParseInt(env, 0, 8)
				if err != nil {
					lg.Warn().Str("zerolog_level", env).Msg("Unable to set zerolog.Level")
				} else {
					logLevel = zerolog.Level(level)
				}

				// If possible, enhance the logger with info from the lambda context
				if lc, ok := lambdacontext.FromContext(ctx); ok {
					lg = log.With().
						Str("commit_id", CommitID).
						Str("request_id", lc.AwsRequestID).
						Str("lambda_function", lambdacontext.FunctionName).
						Logger().
						Level(zerolog.Level(level))
				} else {
					lg = log.With().
						Str("commit_id", CommitID).
						Logger().
						Level(zerolog.Level(level))
					lg.Warn().Msg("Lambda context not found")
				}

				// zerolog usage note: must use Msg() or Send() to trigger logs to actually send
				logger = &lg
				logger.Trace().Str("log_level", logger.GetLevel().String()).Msg("Logging initialized")
			},
		)
	}

	return logger
}

func AppendContext(logger *zerolog.Logger, ctx context.Context) *zerolog.Logger {
	// If possible, enhance the logger with info from the lambda context
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		lg := logger.With().
			Str("request_id", lc.AwsRequestID).
			Str("lambda_function", lambdacontext.FunctionName).
			Logger()
		return &lg
	}

	return logger
}
