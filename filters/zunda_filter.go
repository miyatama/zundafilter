package filters

import (
	"go.uber.org/zap"
	"zundafilter/zunda_mecab"
)

type Converter interface {
	Convert(text string) (string, error)
}

type ZundaFilter struct {
	ZundaDb      ZundaDbController
	MecabWrapper *zunda_mecab.MecabWrapper
	Logger       *zap.Logger
}

func (z *ZundaFilter) Convert(text string) (string, error) {
	defer z.Logger.Sync()
	sugar := z.Logger.Sugar()
	sugar.Debug("ZundaFilter#Convert()")

	converters := []Converter{
		&HonorificFilter{
			ZundaDb:      z.ZundaDb,
			MecabWrapper: z.MecabWrapper,
			Logger:       z.Logger,
		},
		&MoodFilter{
			MecabWrapper: z.MecabWrapper,
			Logger:       z.Logger,
		},
		&PronounFilter{
			MecabWrapper: z.MecabWrapper,
			Logger:       z.Logger,
		},
	}
	var convertedText = text
	for _, converter := range converters {
		resultText, err := converter.Convert(convertedText)
		if err != nil {
			return "", err
		}
		convertedText = resultText
		sugar.Infof("ZundaFilter#filtered() - %s", resultText)
	}
	return convertedText, nil
}
