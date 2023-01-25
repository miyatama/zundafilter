package filters

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"zundafilter/zunda_mecab"
)

type TestZundaDbAccessor struct{}

func (t TestZundaDbAccessor) SelectConvertVerbConjugationTable(baseWord string) (zunda_mecab.ConvertVerbConjugationRow, error) {
	switch baseWord {
	case "渡す":
		return zunda_mecab.ConvertVerbConjugationRow{
			BaseWord: "渡す",
			Mizen:    "渡さ",
		}, nil
	default:
		return zunda_mecab.ConvertVerbConjugationRow{}, nil
	}
}

func getTestLogger() *zap.Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	myConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := myConfig.Build()
	return logger
}

func TestHonorificFilter(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		expect string
	}{
		{
			name:   "敬語無し",
			text:   "ここからなのだ",
			expect: "ここからなのだ",
		},
		{
			name:   "格助詞+末尾敬語",
			text:   "ここからです",
			expect: "ここから",
		},
		{
			name:   "形容詞+末尾敬語",
			text:   "美しいです",
			expect: "美しい",
		},
		{
			name:   "動詞+末尾敬語",
			text:   "闘います",
			expect: "闘う",
		},
		{
			name:   "名詞+敬語+助詞",
			text:   "良いですか？",
			expect: "良いか？",
		},
		{
			name:   "助詞+敬語+助詞",
			text:   "良いのですか？",
			expect: "良いのか？",
		},
		{
			name:   "動詞+敬語+助詞",
			text:   "そう思いますよ",
			expect: "そう思うよ",
		},
		{
			name:   "接続敬語",
			text:   "ですが部長ならできるよ",
			expect: "だけど部長ならできるよ",
		},
		{
			name:   "動詞-五段活用+敬語(まし)+過去",
			text:   "昨日は話しました",
			expect: "昨日は話した",
		},
		{
			name:   "動詞-上一段活用+敬語(まし)+過去",
			text:   "ここで降りました",
			expect: "ここで降りた",
		},
		{
			name:   "動詞-下一段活用+敬語(まし)+過去",
			text:   "ごみを集めました",
			expect: "ごみを集めた",
		},
		{
			name:   "動詞-カ行変格活用+敬語(まし)+過去",
			text:   "昨日来ました",
			expect: "昨日来た",
		},
		{
			name:   "動詞-サ行変格活用+敬語(まし)+過去",
			text:   "勉強しました",
			expect: "勉強した",
		},
		{
			name:   "動詞-イ音便変化+敬語(まし)+過去",
			text:   "さっき書きました",
			expect: "さっき書いた",
		},
		{
			name:   "動詞-促音便変化+敬語(まし)+過去",
			text:   "手を切りました",
			expect: "手を切った",
		},
		{
			name:   "動詞-撥音便変化+敬語(まし)+過去",
			text:   "昨日読みました",
			expect: "昨日読んだ",
		},
		{
			name:   "名詞+敬語(でし)+過去",
			text:   "体調不良でした",
			expect: "体調不良だった",
		},
		{
			name:   "敬語(否定)",
			text:   "これは渡しません",
			expect: "これは渡さない",
		},
		{
			name:   "敬語(否定)+助詞:「これは渡しませんが」",
			text:   "これは渡しませんが",
			expect: "これは渡さないが",
		},
		{
			name:   "敬語(否定)+助詞+敬語(過去):「結局渡しませんでした」",
			text:   "結局渡しませんでした",
			expect: "結局渡さなかった",
		},
		{
			name:   "敬語(否定)+助詞+敬語(過去)+助詞:「結局渡しませんでしたか？」",
			text:   "結局渡しませんでしたか？",
			expect: "結局渡さなかったか？",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			// t.Parallel()
			zundaDbAccessor := TestZundaDbAccessor{}
			mecabWrapper := zunda_mecab.MecabWrapper{
				Logger: getTestLogger(),
			}
			honorificFilter := HonorificFilter{
				ZundaDb:      zundaDbAccessor,
				MecabWrapper: &mecabWrapper,
				Logger:       getTestLogger(),
			}
			actual, _ := honorificFilter.Convert(testCase.text)
			if actual != testCase.expect {
				t.Fatalf("HonorificFilter.Convert() = %v, expect %v", actual, testCase.expect)
			}
		})
	}
}
