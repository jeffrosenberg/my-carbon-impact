package logging

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		Int64("duration_ms", time.Since(start).Milliseconds()).
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
				// TODO: Am I initializing Zerolog correctly/optimally?
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
						Level(zerolog.Level(logLevel))
				} else {
					lg = log.With().
						Str("commit_id", CommitID).
						Logger().
						Level(zerolog.Level(logLevel))
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

// GetTestLogger returns a console logger used for testing
// Because this isn't tied to the Mutex used for application logging,
// it returns a direct struct rather than a pointer.
func GetTestLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	return zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}
