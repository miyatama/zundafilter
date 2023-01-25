package filters

import (
	"go.uber.org/zap"
	"strings"
	"zundafilter/zunda_mecab"
)

type MoodFilter struct {
	MecabWrapper *zunda_mecab.MecabWrapper
	Logger       *zap.Logger
}
type MoodConvertResult struct {
	Features []zunda_mecab.MecabFeature
	Parsed   bool
}

func (m *MoodFilter) Convert(text string) (string, error) {
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

	converters := []func([]zunda_mecab.MecabFeature) MoodConvertResult{
		m.convertNaiAdjectiveMood,
		m.convertPastMood,
		m.convertProhibitionMood,
		m.convertDesireMood,
		m.convertAllowMood,
		m.convertConclusionConversationMood,
		m.convertInvitationMood,
		m.convertOrderTaigenMood,
		m.convertOrderMood,
		m.convertConfirmationMood,
		m.convertQuestionIntentionMood,
		m.convertQuestionMood,
		m.convertUndecisionMood,
		m.convertAgreementMood,
		m.convertRequestMood,
		m.convertConfidenceMood,
		m.convertIntention2Mood,
		m.convertIntentionMood,
		m.convertGuessMood,
		m.convertPossibilityMood,
		m.convertAnxietyMood,
		m.convertAffirmativeMood,
	}
	for _, converter := range converters {
		MoodConvertResult := converter(features)
		if MoodConvertResult.Parsed {
			return m.MecabWrapper.Construct(MoodConvertResult.Features), nil
		}
	}

	return text, nil

}

/*
* 確認のムード
* 条件:  "だろう" + 記号{0..*}
* ex) 昨日、一緒に腹筋しただろう。
 */
func (m *MoodFilter) convertConfirmationMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertConfirmationMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "だろ",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "だ",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "う",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "う",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}

	match, exchangeIndex := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertConfirmationMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextCondition := conditions[2:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextCondition)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:exchangeIndex+2]))
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertConfirmationMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertConfirmationMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 断定のムード
* 条件: 名詞 + "だ" + φ
*   + φ <- [記号 | $ | φ]
* ex) これが正義だ
 */
func (m *MoodFilter) convertAffirmativeMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertAffirmativeMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeNoun,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOneOrNothing,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "だ",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "だ",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertAffirmativeMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// だ+記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+1])) // 名詞まで含める
	texts = append(texts, "なのだ")
	afterTextMatch, afterTextIndex:= m.MecabWrapper.GetMatchIndex(features, conditions[2:])
	if afterTextMatch{
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertAffirmativeMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertAffirmativeMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 可能性のムード
* 条件: "かもしれない" + 記号{0..*}
* ex) 僕は腹筋できるかもしれない。
 */
func (m *MoodFilter) convertPossibilityMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertPossibilityMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "かも",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm: true,
					OriginalForm:      "かも",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "しれ",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm: true,
					OriginalForm:      "しれる",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "ない",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "ない",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertPossibilityMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 3) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+3]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+3:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertPossibilityMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertPossibilityMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 推量のムード
* 条件: "らしい" + 記号{0..*} + EOS
* ex) 僕は腹筋するらしい。
 */
func (m *MoodFilter) convertGuessMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertGuessMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         true,
					Word:              "らしい",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm: true,
					OriginalForm:      "らしい",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     true,
					WordType:          zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:         false,
					Word:              "",
					CheckWordType:     false,
					WordType:          zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm: false,
					OriginalForm:      "",
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertGuessMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+1])) // 「らしい」まで含める
	texts = append(texts, "のだ")
	if (index + 1) <= len(features) {
		texts = append(texts, m.MecabWrapper.Construct(features[index+1:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertGuessMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertGuessMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 意志のムード
* 条件: 動詞(基本形) + 記号{0..*} + EOS
* ex) 僕は腹筋する
 */
func (m *MoodFilter) convertIntentionMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertIntentionMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: true,
					ConjugationForm:      zunda_mecab.MecabConjugationFormKihon,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertIntentionMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+1])) // 動詞まで含める
	texts = append(texts, "のだ")
	if (index + 1) <= len(features) {
		texts = append(texts, m.MecabWrapper.Construct(features[index+1:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertIntentionMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertIntentionMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 確信のムード
* 条件: はず + だ{0..1} + 記号{0..*} + EOS
* ex) 僕は腹筋する
 */
func (m *MoodFilter) convertConfidenceMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertConfidenceMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "はず",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeNoun,
					CheckOriginalForm:    true,
					OriginalForm:         "はず",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOneOrNothing,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "だ",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "だ",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertConfidenceMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	unExchangeConditions := conditions[2:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, unExchangeConditions)
	if !match {
		sugar.Debug("convertConfidenceMood() - after text condition not exists")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+1])) // 「はず」まで含める
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertConfidenceMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertConfidenceMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 勧誘のムード
* 条件: 動詞(連用形) + ましょ + う + 記号{0..*} + EOS
* ex) 僕は腹筋する
 */
func (m *MoodFilter) convertInvitationMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertInvitationMood()")
	// し      動詞,自立,*,*,サ変・スル,連用形,する,シ,シ
	// ましょ  助動詞,*,*,*,特殊・マス,未然ウ接続,ます,マショ,マショ
	// う      助動詞,*,*,*,不変化型,基本形,う,ウ,ウ
	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
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
					CheckWord:            true,
					Word:                 "ましょ",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "ます",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "う",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "う",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
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
		sugar.Debug("convertInvitationMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextConditions := conditions[3:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextConditions)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex+3]))
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertInvitationMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertInvitationMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 依頼のムード
* 条件: 動詞 + て + ください + 記号{0..*} + EOS
* ex) 一緒に腹筋してください
 */
func (m *MoodFilter) convertRequestMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertRequestMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
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
					CheckWord:            true,
					Word:                 "て",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "て",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "ください",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "くださる",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, exchangeIndex := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertRequestMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}
	sugar.Debugf("convertRequestMood() - mood head feature(%d): %s", exchangeIndex, features[exchangeIndex].String())

	// ムード後の文字列
	unExchangeConditions := conditions[3:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, unExchangeConditions)
	sugar.Debugf("convertRequestMood() - after text exists: %v, index: %d", afterTextMatch, afterTextIndex)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:exchangeIndex+3])) // 「ください」まで含める
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}

	sugar.Debugf("convertRequestMood() - exchange text: %v", texts)
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertRequestMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertRequestMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 同意のムード(変換無し)
* 条件: ね + 記号{0..*} + EOS
* ex) 一緒に腹筋したいね
 */
func (m *MoodFilter) convertAgreementMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertAgreementMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "ね",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "ね",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, _ := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertAgreementMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}
	sugar.Infof("convertAgreementMood() - converted: %s", m.MecabWrapper.Construct(features))
	return MoodConvertResult{Features: features, Parsed: true}
}

/*
* 非断定のムード
* 条件: と + 思う + 記号{0..*} + EOS
* ex) 僕は腹筋できると思う
 */
func (m *MoodFilter) convertUndecisionMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertUndecisionMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "と",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "と",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "思う",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "思う",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, _ := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertUndecisionMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextConditions := conditions[2:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextConditions)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:afterTextIndex]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertUndecisionMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertUndecisionMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 質問のムード
* 条件: か + 記号{0..*} + EOS
* ex) 一緒腹筋したか
 */
func (m *MoodFilter) convertQuestionMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertQuestionMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "か",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "か",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
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
		sugar.Debug("convertQuestionMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextConditions := conditions[1:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextConditions)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertQuestionMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertQuestionMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 質問(意志)のムード
* 条件: のか + 記号{0..*} + EOS
* ex) 一緒腹筋したのか
 */
func (m *MoodFilter) convertQuestionIntentionMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertQuestionIntentionMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "の",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeNoun,
					CheckOriginalForm:    true,
					OriginalForm:         "の",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "か",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "か",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
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
		sugar.Debug("convertQuestionIntentionMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (conditionIndex + 2) <= len(features)

	// のか + 記号{0..*} + EOS -> のだ + 記号{0..*} + EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[conditionIndex+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertQuestionIntentionMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertQuestionIntentionMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 命令のムード
* 条件: 動詞 + なさい + 記号{0..*} + EOS
* ex) 一緒に腹筋しなさい
 */
func (m *MoodFilter) convertOrderMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertOrderMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
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
					CheckWord:            true,
					Word:                 "なさい",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "なさる",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
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
		sugar.Debug("convertOrderMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextConditions := conditions[2:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextConditions)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex+2]))
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertOrderMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertOrderMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 命令のムード(体言止め)
* 条件: 動詞(命令) + 記号{0..*} + EOS
* ex) 一緒に闘え
 */
func (m *MoodFilter) convertOrderTaigenMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertOrderTaigenMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
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
		sugar.Debug("convertOrderTaigenMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}
	// 動詞の活用が命令形
	if !features[conditionIndex].ConjugationForm.IsMeirei() {
		sugar.Debug("convertOrderTaigenMood() - not meirei conjugation form")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextConditions := conditions[1:]
	afterTextMatch, afterTextIndex := m.MecabWrapper.GetMatchIndex(features, afterTextConditions)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:conditionIndex+1]))
	texts = append(texts, "なのだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[afterTextIndex:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertOrderTaigenMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertOrderTaigenMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 断定-口頭のムード
* 条件: んだ
* ex) 僕は腹筋できると思う
 */
func (m *MoodFilter) convertConclusionConversationMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertConclusionConversationMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "ん",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeNoun,
					CheckOriginalForm:    true,
					OriginalForm:         "ん",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "だ",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "だ",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertConclusionConversationMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 1) < len(features)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertConclusionConversationMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertConclusionConversationMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 許可のムード
* 条件: してもよい + 記号{0..*} + EOS
* ex) 一緒に腹筋してもよい
 */
func (m *MoodFilter) convertAllowMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertAllowMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "し",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "する",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "て",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "て",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "も",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "も",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "よい",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAdjective,
					CheckOriginalForm:    true,
					OriginalForm:         "よい",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertAllowMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 4) < len(features)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+4]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+4:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertAllowMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertAllowMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 願望のムード
* 条件: したい + 記号{0..*} + EOS
* ex) 僕は腹筋したい
 */
func (m *MoodFilter) convertDesireMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertDesireMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "し",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "する",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "たい",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "たい",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertDesireMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 2) < len(features)

	// 記号{0..*}+EOS -> なのだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+2]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertDesireMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertDesireMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 禁止のムード
* 条件: 動詞 + てはいけない + 記号{0..*} + EOS
* ex) 一緒に腹筋してはいけない。
 */
func (m *MoodFilter) convertProhibitionMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertProhibitionMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
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
					CheckWord:            true,
					Word:                 "て",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "て",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "は",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckOriginalForm:    true,
					OriginalForm:         "は",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "いけ",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "いける",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "ない",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckOriginalForm:    true,
					OriginalForm:         "ない",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertProhibitionMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 5) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+5]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+5:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertProhibitionMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertProhibitionMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 意志のムード
* 条件: (動詞|形容詞|助動詞) + の + 記号{0..*} + EOS
* ex) ここで良いの, ここが大事なの？
 */
func (m *MoodFilter) convertIntention2Mood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertIntention2Mood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAdjective,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
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
					CheckWord:            true,
					Word:                 "の",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeParticle,
					CheckWordSubType1:    true,
					WordSubType1:         zunda_mecab.MecabWordSubType1ParticleTail,
					CheckOriginalForm:    true,
					OriginalForm:         "の",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertIntention2Mood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 2) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+2]))
	texts = append(texts, "だ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertIntention2Mood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertIntention2Mood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 動詞分-過去-た
* 条件: (動詞|形容詞|助動詞) + た + 記号{0..*} + EOS
* ex) ここで良いの, ここが大事なの？
 */
func (m *MoodFilter) convertPastMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertPastMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeVerb,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAdjective,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
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
					CheckWord:            true,
					Word:                 "た",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckWordSubType1:    false,
					WordSubType1:         zunda_mecab.MecabWordSubType1None,
					CheckOriginalForm:    true,
					OriginalForm:         "た",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertPastMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 2) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+2]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertPastMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertPastMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* ナイ形容詞語幹
* 条件: 名詞-ナイ形容詞語幹 + ない + 記号{0..*} + EOS
* ex) それはしょうがない
 */
func (m *MoodFilter) convertNaiAdjectiveMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertNaiAdjectiveMood()")

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
					WordSubType1:         zunda_mecab.MecabWordSubType1NounAdjectiveNai,
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
					CheckWord:            true,
					Word:                 "ない",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeAuxiliaryVerb,
					CheckWordSubType1:    false,
					WordSubType1:         zunda_mecab.MecabWordSubType1None,
					CheckOriginalForm:    true,
					OriginalForm:         "ない",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertNaiAdjectiveMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 2) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+2]))
	texts = append(texts, "のだ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+2:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertNaiAdjectiveMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertNaiAdjectiveMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}

/*
* 不安の「の」
* 条件: の + 記号{0..*} + EOS
* ex) それはいいの
 */
func (m *MoodFilter) convertAnxietyMood(features []zunda_mecab.MecabFeature) MoodConvertResult {
	defer m.Logger.Sync()
	sugar := m.Logger.Sugar()
	sugar.Debug("convertAnxietyMood()")

	conditions := []zunda_mecab.MecabCondition{
		{
			ConditionType: zunda_mecab.MecabConditionTypeOne,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            true,
					Word:                 "の",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckWordSubType1:    false,
					WordSubType1:         zunda_mecab.MecabWordSubType1NounAdjectiveNai,
					CheckOriginalForm:    true,
					OriginalForm:         "の",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeNothingOrContinue,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        true,
					WordType:             zunda_mecab.MecabWordTypeSymbol,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
		{
			ConditionType: zunda_mecab.MecabConditionTypeEOS,
			Features: []zunda_mecab.MecabConditionFeature{

				{
					CheckWord:            false,
					Word:                 "",
					CheckWordType:        false,
					WordType:             zunda_mecab.MecabWordTypeUnknown,
					CheckOriginalForm:    false,
					OriginalForm:         "",
					CheckConjugationForm: false,
					ConjugationForm:      zunda_mecab.MecabConjugationFormNone,
				},
			},
		},
	}

	match, index := m.MecabWrapper.GetMatchIndex(features, conditions)
	if !match {
		sugar.Debug("convertAnxietyMood() - not match")
		return MoodConvertResult{Features: features, Parsed: false}
	}

	// ムード後の文字列
	afterTextMatch := (index + 1) < len(features)

	// 記号{0..*}+EOS -> のだ+記号{0..*}+EOS
	texts := []string{}
	texts = append(texts, m.MecabWrapper.Construct(features[:index+1]))
	texts = append(texts, "だ")
	if afterTextMatch {
		texts = append(texts, m.MecabWrapper.Construct(features[index+1:]))
	}
	exchangeFeatures, err := m.MecabWrapper.ParseToNode(strings.Join(texts, ""))
	if err != nil {
		sugar.Errorf("convertAnxietyMood() - %v", err)
		return MoodConvertResult{Features: features, Parsed: false}
	}

	sugar.Infof("convertAnxietyMood() - converted: %s", m.MecabWrapper.Construct(exchangeFeatures))
	return MoodConvertResult{Features: exchangeFeatures, Parsed: true}
}
