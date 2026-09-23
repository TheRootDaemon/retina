package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSugaredLogger(t *testing.T) {
	var buf bytes.Buffer

	sugar := newSugaredLogger(&buf, LevelInfo)
	sugar.Infof("hello %s", "world")

	assert.Contains(t, buf.String(), "hello world")
}

func TestNew(t *testing.T) {
	l := New(LevelWarn, &bytes.Buffer{})
	assert.NotNil(t, l.sugar)
	assert.Equal(t, LevelWarn, l.level)
}
