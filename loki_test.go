package loki

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var testFilePath = "test.log"

func TestLoggerExample(t *testing.T) {
	SetLevel(DEBUG)
	SetFormatter(NewStandardFormatter())
	msg := "Hi, I'm loki! 你好，我是洛基！"

	Debug("msg: %s", msg)
	Info("msg: %s", msg)
	Warn("msg: %s", msg)
	Error("msg: %s", msg)
	Error(msg)
	Error()
	//Fatal("msg: %s", msg)

	l := New("app:xxx")
	l.Info("msg: %s", msg)
}

func TestLoggerCompile(t *testing.T) {
	l := New("app:xxx")
	l.SetTimeFormatter("")
	assert.Equal(t, l.Compile("%d-%d", 1, 2), "app:xxx 1-2")
}

func TestLoggerCheck(t *testing.T) {
	//LOKI_ENV=app:xxx
	//root
	assert.True(t, logger.Check())

	l1 := New("app:xxx")
	assert.True(t, l1.Check())

	l2 := New("app:xxx1")
	assert.True(t, l2.Check())

	l3 := New("app:x")
	assert.False(t, l3.Check())
}

type ErrFormatter struct {
	Formatter
}

func (f ErrFormatter) format(a ...interface{}) string {
	err := a[0].(error)
	return fmt.Sprintf("Error %v", err)
}

func TestLoggerFormatter(t *testing.T) {
	f := ErrFormatter{}
	SetFormatter(f)
	Info(errors.New("test error"))

	SetFormatter(NewStandardFormatter())
}

func TestLoggerFileHandler(t *testing.T) {
	//timeFormat
	fp, err := os.OpenFile(testFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer fp.Close()
	assert.NoError(t, err)
	SetHandler(NewFileHandler(fp, 10))
	SetLevel(INFO)
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	content, _ := ioutil.ReadFile(testFilePath)
	assert.Equal(t, "", string(content))
	time.Sleep(time.Second)

	content, _ = ioutil.ReadFile(testFilePath)
	assert.NotContains(t, string(content), "debug")
	assert.Contains(t, string(content), "info")
	assert.Contains(t, string(content), "warn")
	assert.Contains(t, string(content), "error")
}

func TestLoggerSetting(t *testing.T) {
	//timeFormat
	l := New("app:xxx")
	assert.Contains(t, l.Compile("xxx"), "app:xxx xxx")
	l.SetTimeFormatter("")
	assert.Equal(t, "app:xxx xxx", l.Compile("xxx"))
	l.SetTimeFormatter(time.RFC822)
	assert.Contains(t, l.Compile("xxx"), "app:xxx xxx")
}
