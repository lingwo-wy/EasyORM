package log

import (
	"os"
	"testing"
)

func TestSetLevel(t *testing.T) {
		SetLevel(ErrorLevel)
		if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
			t.Fatal("设置日志层级错误")
		}
		SetLevel(Disabled)
		if infoLog.Writer() == os.Stdout || errorLog.Writer() == os.Stdout {
			t.Fatal("设置日志层级错误")
		}
}
