package filters

import (
	"go.uber.org/zap"
	"strings"
	"zundafilter/zunda_mecab"
)

type HonorificFilter struct {
	ZundaDb      ZundaDbController
	MecabWrapper *zunda_mecab.MecabWrapper
	Logger       *zap.Logger
}

func (h *HonorificFilter) Convert(text string) (string, error) {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("HonorificFilter#Convert()")

	features, err := h.MecabWrapper.ParseToNode(text)
	if err != nil {
		return "", nil
	}

	convertedFeatures := convert(
		features,
		[]func([]zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature{
			h.convertVerbBeforePastHonorificNegative,
			h.convertSahenVerbBeforePastHorific,
			h.convertVerbBeforePastHorificHatsuOnbin,
			h.convertVerbBeforePastHorificIOnbin,
			h.convertVerbBeforePastHorificSokuOnbin,
			h.convertNounBeforePastHorific,
			h.convertVerbBeforeHonorificNegative,
			h.removePastHonorificWord,
			h.convertVerbBeforeHonorific,
			h.convertSpecials,
			h.removeHonorificWord,
		})

	return h.MecabWrapper.Construct(convertedFeatures), nil

}

func convert(features []zunda_mecab.MecabFeature, converters []func([]zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	if len(converters) == 0 {
		return features
	}
	return convert(
		converters[0](features),
		converters[1:],
	)
}

/*
* 敬語判断条件(現在)
 */
func getHonorificWords() []zunda_mecab.MecabCondition {
	return []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "です",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
				{
					CheckWord:         true,
					Word:              "ます",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
				{
					CheckWord:         true,
					Word:              "ですが",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeConjunction,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}
}

/*
* 動詞 + 敬語(否定)の対応
* ex) ここから動きません -> ここから動かない
 */
func (h *HonorificFilter) convertVerbBeforeHonorificNegative(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()

	sugar.Debug("convertVerbBeforeHonorificNegative()")
	conditions := []zunda_mecab.MecabCondition{

		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "ませ",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "ん",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ん",
				},
			},
		},
	}
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	// 条件:
	//  + 敬語の一つ前が動詞
	//  + 敬語の一つ後が「ん」
	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index]))
	conjugationForm, err := h.ZundaDb.SelectConvertVerbConjugationTable(features[index].OriginalForm)
	if err != nil {
		sugar.Errorf("convertVerbBeforeHonorificNegative() - error: %v", err)
		return features
	}
	texts = append(texts, conjugationForm.Mizen)
	texts = append(texts, "ない")
	if (index + len(conditions)) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[index+len(conditions):]))
	}

	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertVerbBeforeHonorificNegative() - error: %v", err)
		return features
	}
	sugar.Infof("convertVerbBeforeHonorificNegative() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures

}

/*
* 動詞変換を含む敬語解除(現在)
 */
func (h *HonorificFilter) convertVerbBeforeHonorific(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()

	sugar.Debug("convertVerbBeforeHonorific()")

	conditions := getHonorificWords()
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	// 条件:
	// + 敬語の一つ前が動詞
	// + 動詞が原形ではない
	// + 敬語が現在形
	if index < 1 {
		return features
	}
	if features[index-1].WordType != zunda_mecab.MecabWordTypeVerb {
		return features
	}
	if features[index-1].Word == features[index-1].OriginalForm {
		return features
	}
	exchangedFeatures := []zunda_mecab.MecabFeature{}
	exchangedFeatures = append(exchangedFeatures, features[:index-1]...)
	exchangedFeatures = append(exchangedFeatures, zunda_mecab.MecabFeature{
		EOS:             features[index-1].EOS,
		Word:            features[index-1].OriginalForm,
		WordType:        features[index-1].WordType,
		WordSubType1:    features[index-1].WordSubType1,
		WordSubType2:    features[index-1].WordSubType2,
		WordSubType3:    features[index-1].WordSubType3,
		ConjugationType: features[index-1].ConjugationType,
		ConjugationForm: features[index-1].ConjugationForm,
		OriginalForm:    features[index-1].OriginalForm,
	})
	exchangedFeatures = append(exchangedFeatures, features[index:]...)
	sugar.Infof("convertVerbBeforeHonorific() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 特定敬語(現在)の変換(「ですが」など)
 */
func (h *HonorificFilter) convertSpecials(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertSpecials()")

	conditions := []struct {
		Condition []zunda_mecab.MecabCondition
		DistWord  string
	}{
		{
			Condition: []zunda_mecab.MecabCondition{
				{
					ConditionType: zunda_mecab.MecabConditionTypeOne,
					Features: []zunda_mecab.MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "ですが",
							CheckWordType:     true,
							WordType:          zunda_mecab.MecabWordTypeConjunction,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			DistWord: "だけど",
		},
	}

	var exchangedFeatures = features
	for _, condition := range conditions {
		sugar.Debugf("convertSpecials() special: %s", condition.Condition[0].String())
		match, index := h.MecabWrapper.GetMatchIndex(features, condition.Condition)
		if !match {
			continue
		}
		var texts = []string{}
		texts = append(texts, h.MecabWrapper.Construct(exchangedFeatures[:index]))
		texts = append(texts, condition.DistWord)
		if (index + 1) <= len(exchangedFeatures) {
			texts = append(texts, h.MecabWrapper.Construct(exchangedFeatures[index+1:]))
		}

		sugar.Debugf("convertSpecials(): exchange")

		parseResult, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
		if err != nil {
			sugar.Errorf("convertSpecials() - error: %v", err)
			return features
		}
		exchangedFeatures = parseResult
	}

	sugar.Infof("convertSpecials() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 敬語削除(現在)
 */
func (h *HonorificFilter) removeHonorificWord(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("removeHonorificWord()")

	conditions := getHonorificWords()
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index]))
	if (index + 1) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[index+1:]))
	}
	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("removeHonorificWord() - error: %v", err)
		return features
	}
	sugar.Infof("removeHonorificWord() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 敬語判断条件(過去)
 */
func getPastHonorificWords() []zunda_mecab.MecabCondition {
	return []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "でし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "です",
				},
				{
					CheckWord:         true,
					Word:              "まし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
	}
}

/*
* 動詞 + 敬語(現在) + ん + 敬語(過去) + た
* ex) 彼は動きませんでした -> 彼は動かなかった
 */
func (h *HonorificFilter) convertVerbBeforePastHonorificNegative(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertVerbBeforePastHonorificNegative()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "ませ",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "ん",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ん",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "でし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "です",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "た",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "た",
				},
			},
		},
	}
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	// 条件:
	//  + 敬語の一つ前が動詞
	conjugationForm, err := h.ZundaDb.SelectConvertVerbConjugationTable(features[index].OriginalForm)
	if err != nil {
		sugar.Errorf("convertVerbBeforePastHonorificNegative() - can not fetch verb conjugation form(%v) - %v", features[index-1].OriginalForm, err)
		return features
	}

	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index]))
	texts = append(texts, conjugationForm.Mizen+"なかった")
	if (index + len(conditions)) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[(index+len(conditions)):]))
	}

	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertVerbBeforePastHonorificNegative() - error: %v", err)
		return features
	}

	sugar.Infof("convertVerbBeforePastHonorificNegative() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 動詞(サ行変格活用) + 敬語(過去)の変換
* 「する」が五段活用と認識される為、別途変換を実施
 */
func (h *HonorificFilter) convertSahenVerbBeforePastHorific(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()

	sugar.Debug("convertSahenVerbBeforePastHorific()")
	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckConjugationType: true,
					ConjugationType:      zunda_mecab.MecabConjugationTypeSahenSuru,
					CheckOriginalForm:    false,
					OriginalForm:         "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "でし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "です",
				},
				{
					CheckWord:         true,
					Word:              "まし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "た",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: true,
					OriginalForm:      "た",
				},
			},
		},
	}

	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index+1]))
	texts = append(texts, h.MecabWrapper.Construct(features[index+2:]))

	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertVerbBeforePastHonorificNegative() - error: %v", err)
		return features
	}

	sugar.Infof("convertSahenVerbBeforePastHorific() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 名詞 + 敬語(過去)の変換
 */
func (h *HonorificFilter) convertNounBeforePastHorific(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertNounBeforePastHorific()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeNoun,
					CheckConjugationType: false,
					ConjugationType:      zunda_mecab.MecabConjugationTypeNone,
					CheckOriginalForm:    false,
					OriginalForm:         "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "でし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "です",
				},
				{
					CheckWord:         true,
					Word:              "まし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "た",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: true,
					OriginalForm:      "た",
				},
			},
		},
	}

	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		return features
	}

	// 条件
	// 並びが 名詞+敬語(過去)+"た"
	// ex) そこ は 雪国 でし た
	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index+1]))
	texts = append(texts, "だった")
	if (index + 3) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[index+3:]))
	}

	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertNounBeforePastHorific() - error: %v", err)
		return features
	}

	sugar.Infof("convertNounBeforePastHorific() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 動詞(撥音便) + 敬語(過去)の変換
 */
func (h *HonorificFilter) convertVerbBeforePastHorificHatsuOnbin(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertVerbBeforePastHorificHatsuOnbin()")
	// 条件
	// 敬語の直前が動詞
	// 敬語の直後が「た」
	// 動詞の活用が五段・ナ行,五段・バ行,五段・マ行
	exchangedFeatures := h.convertVerbBeforePastHorificOnbin(features,
		[]zunda_mecab.MecabConjugationType{
			zunda_mecab.MecabConjugationTypeGodanNa,
			zunda_mecab.MecabConjugationTypeGodanBa,
			zunda_mecab.MecabConjugationTypeGodanMa,
		},
		"た",
		"んだ")
	sugar.Infof("convertVerbBeforePastHorificHatsuOnbin() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 動詞(イ音便) + 敬語(過去)の変換
 */
func (h *HonorificFilter) convertVerbBeforePastHorificIOnbin(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertVerbBeforePastHorificIOnbin()")
	// 条件
	// 敬語の直前が動詞
	// 敬語の直後が「た」
	// 動詞の活用が五段・カ行イ音便, 五段・ガ行
	exchangedFeatures := h.convertVerbBeforePastHorificOnbin(features,
		[]zunda_mecab.MecabConjugationType{
			zunda_mecab.MecabConjugationTypeGodanKaIOnbin,
			zunda_mecab.MecabConjugationTypeGodanGa,
		},
		"た",
		"いた")
	sugar.Infof("convertVerbBeforePastHorificIOnbin() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 動詞(促音便) + 敬語(過去)の変換
 */
func (h *HonorificFilter) convertVerbBeforePastHorificSokuOnbin(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("convertVerbBeforePastHorificSokuOnbin()")

	// 条件
	// 敬語の直前が動詞
	// 敬語の直後が「た」
	// 動詞の活用が五段・タ行, 五段・ワ行促音便, 五段・ラ行
	exchangedFeatures := h.convertVerbBeforePastHorificOnbin(features,
		[]zunda_mecab.MecabConjugationType{
			zunda_mecab.MecabConjugationTypeGodanTa,
			zunda_mecab.MecabConjugationTypeGodanWaSokuOnbin,
			zunda_mecab.MecabConjugationTypeGodanRa,
		},
		"た",
		"った")
	sugar.Infof("convertVerbBeforePastHorificSokuOnbin() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures

}

/*
* 音便+敬語(過去)の対応
 */
func (h *HonorificFilter) convertVerbBeforePastHorificOnbin(features []zunda_mecab.MecabFeature, conjugationTypes []zunda_mecab.MecabConjugationType, particleAfterHonorific string, replacedAfterText string) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()

	sugar.Debug("convertVerbBeforePastHorificOnbin()")
	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              "でし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "です",
				},
				{
					CheckWord:         true,
					Word:              "まし",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ます",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{
				{
					CheckWord:         true,
					Word:              particleAfterHonorific,
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertVerbBeforePastHorificOnbin() - has not past honorific words")
		return features
	}

	verbFeatures, err := h.MecabWrapper.ParseToNodeWithoutEos(features[index].OriginalForm)
	if err != nil {
		sugar.Errorf("convertVerbBeforePastHorificOnbin() - can not parse %s, %v", err)
		return features
	}
	for _, feature := range verbFeatures {
		sugar.Debug(feature.String())
	}
	var isExchangeConjugationType = false
	for _, conjugationType := range conjugationTypes {
		if verbFeatures[0].ConjugationType == conjugationType {
			isExchangeConjugationType = true
			break
		}
	}
	if !isExchangeConjugationType {
		sugar.Debugf("convertVerbBeforePastHorificOnbin() - invalid verb conjugation type(%s - %s)", verbFeatures[0].Word, verbFeatures[0].ConjugationType.String())
		return features
	}
	slice := []rune(features[index].OriginalForm)
	replacedVerbText := string(slice[0:len(slice)-1]) + replacedAfterText

	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index]))
	texts = append(texts, replacedVerbText)
	if (index + 3) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[index+3:]))
	}

	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertVerbBeforePastHorificOnbin() - error: %v", err)
		return features
	}
	sugar.Infof("convertVerbBeforePastHorificOnbin() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}

/*
* 敬語削除(過去)
 */
func (h *HonorificFilter) removePastHonorificWord(features []zunda_mecab.MecabFeature) []zunda_mecab.MecabFeature {
	defer h.Logger.Sync()
	sugar := h.Logger.Sugar()
	sugar.Debug("removePastHonorificWord()")

	conditions := getPastHonorificWords()
	match, index := h.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("removePastHonorificWord() - has not past honorific words")
		return features
	}

	texts := []string{}
	texts = append(texts, h.MecabWrapper.Construct(features[:index]))
	if (index + 1) <= len(features) {
		texts = append(texts, h.MecabWrapper.Construct(features[index+1:]))
	}
	exchangedFeatures, err := h.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("removePastHonorificWord() - error: %v", err)
		return features
	}
	sugar.Infof("removePastHonorificWord() - converted: %s", h.MecabWrapper.Construct(exchangedFeatures))
	return exchangedFeatures
}
