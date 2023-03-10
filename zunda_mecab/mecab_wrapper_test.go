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
			name: "Unmatch-ĺčŠ+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ĺă",
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
			name: "ĺčŠ+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ä˝ćŚ",
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
			name: "Unmatch-ĺčŠ+ĺŠčŠ+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
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
					Word:            "ĺă",
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
			name: "ĺčŠ+ĺŠĺčŠ+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
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
					Word:            "ă ",
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
			name: "ć+ĺčŠ+č¨ĺˇ{0..*}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ă",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ă",
				},
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ă",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ă",
				},
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ďź",
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
					Word:            "ďź",
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
					Word:            "ďź",
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
							Word:              "ć­ŁçžŠ",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "ć­ŁçžŠ",
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
			name: "ć+ĺčŠ+ĺŠčŠ+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ăă",
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
					Word:            "ă",
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
					Word:            "ć­ŁçžŠ",
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
					Word:            "ă ",
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
			name: "ć+ĺčŠ+ĺŠčŠ+č¨ĺˇ{0..*}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ăă",
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
					Word:            "ă",
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
					Word:            "ć­ŁçžŠ",
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
					Word:            "ă ",
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
					Word:            "ďź",
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
					Word:            "ďź",
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
					Word:            "ďź",
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
			name: "ć+ĺčŠ+č¨ĺˇ{0}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ă",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ă",
				},
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ă",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ă",
				},
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
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
							Word:              "ć­ŁçžŠ",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "ć­ŁçžŠ",
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
			name: "ĺčŠ+ĺŠčŠ{0..1}+č¨ĺˇ{0}+EOS-ĺŠčŠăŞă",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
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
							Word:              "ć­ŁçžŠ",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "ć­ŁçžŠ",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOneOrNothing,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "ă ",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: true,
							OriginalForm:      "ă ",
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
			name: "ĺčŠ+ĺŠčŠ{0..1}+č¨ĺˇ{0}+EOS-ĺŠčŠăă",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "ć­ŁçžŠ",
					WordType:        MecabWordTypeNoun,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ć­ŁçžŠ",
				},
				{
					EOS:             false,
					Word:            "ă ",
					WordType:        MecabWordTypeParticle,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ă ",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
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
							Word:              "ć­ŁçžŠ",
							CheckWordType:     true,
							WordType:          MecabWordTypeNoun,
							CheckOriginalForm: true,
							OriginalForm:      "ć­ŁçžŠ",
						},
					},
				},
				{
					ConditionType: MecabConditionTypeOneOrNothing,
					Features: []MecabConditionFeature{
						{
							CheckWord:         true,
							Word:              "ă ",
							CheckWordType:     true,
							WordType:          MecabWordTypeParticle,
							CheckOriginalForm: true,
							OriginalForm:      "ă ",
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
			name: "ĺčŠ(ĺşćŹĺ˝˘)+č¨ĺˇ{0}+EOS",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "éă",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "éă",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
				},
				{
					EOS:             false,
					Word:            "ďź",
					WordType:        MecabWordTypeSymbol,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormNone,
					OriginalForm:    "ďź",
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
			name: "+č¨ĺˇ{0..*}+EOS-EOSăŽăżăăă",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "éă",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "éă",
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
			name: "+č¨ĺˇ{0..*}+EOS-EOSăŞăăăă",
			features: []MecabFeature{
				{
					EOS:             false,
					Word:            "éă",
					WordType:        MecabWordTypeVerb,
					WordSubType1:    MecabWordSubType1None,
					WordSubType2:    MecabWordSubType2None,
					WordSubType3:    MecabWordSubType3None,
					ConjugationType: MecabConjugationTypeNone,
					ConjugationForm: MecabConjugationFormKihon,
					OriginalForm:    "éă",
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
