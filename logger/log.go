package logger

// Trace logs a message at LevelTrace.
func (l *Logger) Trace(format string, args ...any) {
	if !l.Enabled(LevelTrace) {
		return
	}
	l.sugar.Debugf(format, args...)
}

// Debug logs a message at LevelDebug.
func (l *Logger) Debug(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

// Info logs a message at LevelInfo.
func (l *Logger) Info(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

// Warn logs a message at LevelWarn.
func (l *Logger) Warn(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}

// Error logs a message at LevelError.
func (l *Logger) Error(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

// Log logs a message at the given level.
func (l *Logger) Log(level Level, format string, args ...any) {
	switch level {
	case LevelTrace, LevelDebug:
		l.sugar.Debugf(format, args...)
	case LevelInfo:
		l.sugar.Infof(format, args...)
	case LevelWarn:
		l.sugar.Warnf(format, args...)
	case LevelError:
		l.sugar.Errorf(format, args...)
	}
}
