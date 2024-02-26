// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogMaxSize = 300 // MB
)

// FileLogConfig serializes file log related config in toml/json.
type FileLogConfig struct {
	// Log rootpath
	RootPath string `toml:"rootpath" json:"rootpath"`
	// Log filename, leave empty to disable file log.
	Filename string `toml:"filename" json:"filename"`
	// Max size for a single file, in MB.
	MaxSize int `toml:"max-size" json:"max-size"`
	// Max log keep days, default is never deleting.
	MaxDays int `toml:"max-days" json:"max-days"`
	// Maximum number of old log files to retain.
	MaxBackups int `toml:"max-backups" json:"max-backups"`
	// Compress function for rotated files.
	// Currently only `gzip` and empty are supported, empty means compression disabled.
	Compress string `toml:"compress" json:"compress"`
}

type EncoderConfig struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey     string `json:"messageKey" yaml:"messageKey"`
	LevelKey       string `json:"levelKey" yaml:"levelKey"`
	TimeKey        string `json:"timeKey" yaml:"timeKey"`
	NameKey        string `json:"nameKey" yaml:"nameKey"`
	CallerKey      string `json:"callerKey" yaml:"callerKey"`
	FunctionKey    string `json:"functionKey" yaml:"functionKey"`
	StacktraceKey  string `json:"stacktraceKey" yaml:"stacktraceKey"`
	SkipLineEnding bool   `json:"skipLineEnding" yaml:"skipLineEnding"`
	LineEnding     string `json:"lineEnding" yaml:"lineEnding"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    zapcore.LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     zapcore.TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration zapcore.DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   zapcore.CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName zapcore.NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
	//// Configure the encoder for interface{} type objects.
	//// If not provided, objects are encoded using json.Encoder
	//NewReflectedEncoder func(io.Writer) zapcore.ReflectedEncoder `json:"-" yaml:"-"`
	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
}

// Config serializes log related config in toml/json.
type Config struct {
	// Log level.
	Level string `toml:"level" json:"level"`
	// grpc log level
	GrpcLevel string `toml:"grpc-level" json:"grpc-level"`
	// Log format. one of json, text, or console.
	Format string `toml:"format" json:"format"`
	// Disable automatic timestamps in output.
	DisableTimestamp bool `toml:"disable-timestamp" json:"disable-timestamp"`
	// Stdout enable or not.
	Stdout bool `toml:"stdout" json:"stdout"`
	// File log config.
	File FileLogConfig `toml:"file" json:"file"`
	// Encoder config.
	Encoder EncoderConfig `toml:"encoder" json:"encoder"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `toml:"development" json:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `toml:"disable-caller" json:"disable-caller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `toml:"disable-stacktrace" json:"disable-stacktrace"`
	// DisableErrorVerbose stops annotating logs with the full verbose error
	// message.
	DisableErrorVerbose bool `toml:"disable-error-verbose" json:"disable-error-verbose"`
	// SamplingConfig sets a sampling strategy for the logger. Sampling caps the
	// global CPU and I/O load that logging puts on your process while attempting
	// to preserve a representative subset of your logs.
	//
	// Values configured here are per-second. See zapcore.NewSampler for details.
	Sampling *zap.SamplingConfig `toml:"sampling" json:"sampling"`
}

// ZapProperties records some information about zap.
type ZapProperties struct {
	Core   zapcore.Core
	Syncer zapcore.WriteSyncer
	Level  zap.AtomicLevel
}

func newZapTextEncoder(cfg *Config) zapcore.Encoder {
	return NewTextEncoderByConfig(cfg)
}

func (cfg *Config) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(core, time.Second, cfg.Sampling.Initial, cfg.Sampling.Thereafter, zapcore.SamplerHook(cfg.Sampling.Hook))
		}))
	}
	return opts
}
