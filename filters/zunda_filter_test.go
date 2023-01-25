package filters

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"zundafilter/zunda_mecab"
)

type TestZundaFilterDbAccessor struct{}

func (t TestZundaFilterDbAccessor) SelectConvertVerbConjugationTable(baseWord string) (zunda_mecab.ConvertVerbConjugationRow, error) {
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

func getZundaFilterTestLogger() *zap.Logger {
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

func TestZundaFilter(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		expect string
	}{
		{
			name:   "例文001",
			text:   "この暮になってひどいよ、おれにとっちゃあ一時間が何万円にもつくときだからね",
			expect:   "この暮になってひどいよ、ぼくにとっちゃあ一時間が何万円にもつくときだからね",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			zundaDbAccessor := TestZundaFilterDbAccessor{}
			mecabWrapper := zunda_mecab.MecabWrapper{
				Logger: getZundaFilterTestLogger(),
			}
			filter := ZundaFilter{
				ZundaDb:      zundaDbAccessor,
				MecabWrapper: &mecabWrapper,
				Logger:       getTestLogger(),
			}
			actual, _ := filter.Convert(testCase.text)
			if actual != testCase.expect {
				t.Fatalf("ZundaFilter.Convert() = %v, expect %v", actual, testCase.expect)
			}
		})
	}
}
