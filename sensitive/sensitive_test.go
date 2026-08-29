package sensitive

import (
	"testing"
)

func TestInitWords(t *testing.T) {
	InitWords([]string{"敏感词", "违禁词"})
}

func TestCensorIsPass(t *testing.T) {
	InitWords([]string{"敏感词", "违禁词"})

	if !CensorIsPass("这是正常文本") {
		t.Error("normal text should pass")
	}
}

func TestCensorAndReplace(t *testing.T) {
	InitWords([]string{"敏感词", "违禁词"})

	pass, result := CensorAndReplace("这是正常文本")
	if !pass {
		t.Error("normal text should pass")
	}
	if result != "这是正常文本" {
		t.Errorf("result = %q, want original text", result)
	}
}

func TestCensorWithSensitiveWord(t *testing.T) {
	InitWords([]string{"敏感词"})

	pass, _ := CensorAndReplace("包含敏感词的文本")
	if pass {
		t.Error("text with sensitive word should not pass")
	}
}
