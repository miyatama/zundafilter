package filters

import (
	"go.uber.org/zap"
	"strings"
	"zundafilter/zunda_mecab"
)

type PronounFilter struct {
	MecabWrapper *zunda_mecab.MecabWrapper
	Logger       *zap.Logger
}
type PronounConvertResult struct {
	Features []zunda_mecab.MecabFeature
	Parsed   bool
}

func (m *PronounFilter) Convert(text string) (string, error) {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()

	features, err := m.MecabWrapper.ParseToNode(text)
	if err != nil {
		return "", nil
	}

	// パース結果の出力
	for _, feature := range features {
		sugar.Debug(feature.String())
	}

	converters := []func([]zunda_mecab.MecabFeature) PronounConvertResult{
		m.convertPronoun,
	}
	for _, converter := range converters {
		PronounConvertResult := converter(features)
		if PronounConvertResult.Parsed {
			return m.MecabWrapper.Construct(PronounConvertResult.Features), nil
		}
	}

	return text, nil

}

/*
* 代名詞の変換
* 条件: 代名詞 + 助詞
* ex) 私は野球が好きです
 */
func (m *PronounFilter) convertPronoun(features []zunda_mecab.MecabFeature) PronounConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertPronoun()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeNoun,
					CheckWordSubType1:    true,
					WordSubType1:         zunda_mecab.MecabWordSubType1NounPronoun,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckWordSubType1:    false,
					WordSubType1:         zunda_mecab.MecabWordSubType1None,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, conditionIndex := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertPronoun() - not match")
		return PronounConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (conditionIndex + 1) < len(features)

	// 代名詞 + 助詞 -> ぼく + 助詞
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex]))
	texts = append(texts, "ぼく")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[conditionIndex + 1:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertPronoun() - %v", err)
		return PronounConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertPronoun() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return PronounConvertResult{Features: exchangeFeatures, Parsed: true}
}
