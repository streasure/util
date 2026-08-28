package sensitive

import textcensor "github.com/kai1987/go-text-censor"

func CensorIsPass(text string) bool {
	return textcensor.IsPass(text, true)
}

func CensorAndReplace(text string) (bool, string) {
	return textcensor.CheckAndReplace(text, true, '*')
}

func InitWords(words []string) {
	textcensor.InitWords(words, true)
	defaultPunctuation := "0123456789abcdefghijklmnopqrstuvwxyz !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：\"￥（）——、！……"
	textcensor.SetPunctuation(defaultPunctuation)
}
