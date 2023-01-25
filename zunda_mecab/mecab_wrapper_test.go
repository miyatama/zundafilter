package zunda_mecab

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

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

func TestGetMatchIndex(t *testing.T) {
	tests := []struct {
		name        string
		features    []MecabFeature
		conditions  []MecabCondition
		expectMatch bool
		expectIndex int
	}{
		{
			name:        "nothing condition",
			features:    []MecabFeature{},
			conditions:  []MecabCondition{},
			expectMatch: false,
			expectIndex: 0,
		},
		{
			name:     "EOS",
			features: []MecabFeature{},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "Unmatch-名詞+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "動く",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: false,
			expectIndex: 0,
		},
		{
			name: "名詞+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "作戦",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "Unmatch-名詞+助詞+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "動く",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: false,
			expectIndex: 0,
		},
		{
			name: "名詞+助動詞+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "だ",
					WordType:        MecabWordTypeAuxiliaryVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeAuxiliaryVerb,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "文+名詞+記号{0..*}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "、",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "、",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "、",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "、",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "正義",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "正義",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeSymbol,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 4,
		},
		{
			name: "文+名詞+助詞+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "これ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "が",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "だ",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{

						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 2,
		},
		{
			name: "文+名詞+助詞+記号{0..*}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "これ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "が",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "だ",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeSymbol,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 2,
		},
		{
			name: "文+名詞+記号{0}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "、",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "、",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "、",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "、",
				},
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "正義",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "正義",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeSymbol,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 4,
		},
		{
			name: "名詞+助詞{0..1}+記号{0}+EOS-助詞なし",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "正義",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "正義",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOneOrNothing,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "だ",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: true,
							OriginalForm:      "だ",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeSymbol,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "名詞+助詞{0..1}+記号{0}+EOS-助詞あり",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "正義",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "正義",
				},
				{
					EOS:             false,
					Word:            "だ",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "だ",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "正義",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "正義",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOneOrNothing,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "だ",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: true,
							OriginalForm:      "だ",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     true,
							WordType:          MecabWordTypeSymbol,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:         false,
							Word:              "",
							CheckWordType:     false,
							WordType:          MecabWordTypeUnknown,
							CheckOriginalForm: false,
							OriginalForm:      "",
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "動詞(基本形)+記号{0}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "闘う",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "闘う",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             false,
					Word:            "！",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "！",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeOne,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        true,
							WordType:             MecabWordTypeVerb,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: true,
							ConjugationForm:      MecabConjugationFormKihon,
						},
					},
				},
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        true,
							WordType:             MecabWordTypeSymbol,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        false,
							WordType:             MecabWordTypeUnknown,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 0,
		},
		{
			name: "+記号{0..*}+EOS-EOSのみマッチ",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "闘う",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "闘う",
				},
				{
					EOS:             true,
					Word:            "",
					WordType:        MecabWordTypeUnknown,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        true,
							WordType:             MecabWordTypeSymbol,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        false,
							WordType:             MecabWordTypeUnknown,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
			},
			expectMatch: true,
			expectIndex: 1,
		},
		{
			name: "+記号{0..*}+EOS-EOSなしマッチ",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "闘う",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "闘う",
				},
			},
			conditions: []MecabCondition{
				{
					ConditionType: MecabConditionTypeNothingOrContinue,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        true,
							WordType:             MecabWordTypeSymbol,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
				{
					ConditionType: MecabConditionTypeEOS,
					Features: []MecabConditionFeature{
						{
							CheckWord:            false,
							Word:                 "",
							CheckWordType:        false,
							WordType:             MecabWordTypeUnknown,
							CheckOriginalForm:    false,
							OriginalForm:         "",
							CheckConjugationForm: false,
							ConjugationForm:      MecabConjugationFormNone,
						},
					},
				},
			},
			expectMatch: false,
			expectIndex: 0,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			// t.Parallel()
			wrappeer := MecabWrapper{
				Logger: getTestLogger(),
			}
			actualMatch, actualIndex := wrappeer.GetMatchIndex(testCase.features, testCase.conditions)
			if actualMatch != testCase.expectMatch || actualIndex != testCase.expectIndex {
				t.Fatalf("MecabWrapper.GetMatchIndex() = (%v, %v) expect (%v, %v)", actualMatch, actualIndex, testCase.expectMatch, testCase.expectIndex)
			}
		})
	}
}
