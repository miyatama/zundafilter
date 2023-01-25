package zunda_mecab

import (
	"github.com/bluele/mecab-golang"
	"go.uber.org/zap"
	"strings"
)

type MecabWrapper struct {
	Logger *zap.Logger
}

func (w *MecabWrapper) ParseToNode(text string) ([]MecabFeature, error) {
	defer w.Logger.Sync()
	sugar := w.Logger.Sugar()

	mecabFeatures := []MecabFeature{}
	m, err := mecab.New("-Owakati")
	if err != nil {
		panic(err)
	}
	defer m.Destroy()

	tg, err := m.NewTagger()
	if err != nil {
		return mecabFeatures, err
	}
	defer tg.Destroy()

	lt, err := m.NewLattice(text)
	if err != nil {
		return mecabFeatures, err
	}
	defer lt.Destroy()

	node := tg.ParseToNode(lt)
	for {
		feature := parseMecabFeatureNode(node)
		sugar.Debug(feature.String())
		if len(mecabFeatures) != 0 || !feature.EOS {
			mecabFeatures = append(mecabFeatures, feature)
		}
		if node.Next() != nil {
			break
		}
	}
	return mecabFeatures, nil
}

func (w MecabWrapper) ParseToNodeWithoutEos(text string) ([]MecabFeature, error) {
	features, err := w.ParseToNode(text)
	if err != nil {
		return features, err
	}

	eosRemovedFeatures := []MecabFeature{}
	for _, feature := range features {
		if feature.EOS {
			continue
		}
		eosRemovedFeatures = append(eosRemovedFeatures, feature)
	}
	return eosRemovedFeatures, nil
}

func (w MecabWrapper) Construct(features []MecabFeature) string {
	words := []string{}
	for _, feature := range features {
		words = append(words, feature.Word)
	}
	return strings.Join(words, "")
}

/*
* Featureのパターン一致
* 条件に合致したインデックスを返す
* return
*   [0]: 合致
*   [0]: 合致した先頭インデックス。合致しない場合は-1
 */
func (w MecabWrapper) GetMatchIndex(features []MecabFeature, conditions []MecabCondition) (bool, int) {
	defer w.Logger.Sync()
	sugar := w.Logger.Sugar()
	sugar.Debugf("GetMatchIndex()")

	// EOSのみの指定
	if len(features) <= 0 && len(conditions) == 1 && conditions[0].ConditionType == MecabConditionTypeEOS {
		return true, 0
	}

	for i := 0; i < len(features); i++ {
		match := w.getMatchIndex(features[i:], conditions)
		if !match {
			continue
		}
		return match, i
	}
	return false, 0
}

func (w MecabWrapper) getMatchIndex(features []MecabFeature, conditions []MecabCondition) bool {
	defer w.Logger.Sync()
	sugar := w.Logger.Sugar()
	sugar.Debugf("getMatchIndex()")
	// 条件なしなら無条件に合致
	if len(conditions) <= 0 {
		sugar.Debugf("getMatchIndex() - nothing conditions")
		return true
	}
	match, length := w.matchFeaturesWithCondition(
		features,
		conditions[0])
	if !match {
		sugar.Debugf("getMatchIndex() - unmatch")
		return false
	}
	return w.getMatchIndex(features[length:], conditions[1:])
}

/*
* 条件に一致したインデックスを返す
* return
*   [0]: 合致
*   [1]: 合致した要素数
 */
func (w MecabWrapper) matchFeaturesWithCondition(features []MecabFeature, condition MecabCondition) (bool, int) {
	defer w.Logger.Sync()
	sugar := w.Logger.Sugar()
	sugar.Debugf("matchFeaturesWithCondition() - condition type: %s", condition.ConditionType.String())
	switch condition.ConditionType {
	case MecabConditionTypeOne: // 1つに合致
		// 検証要素無し
		if len(features) <= 0 {
			sugar.Debug("matchFeaturesWithCondition() - feature not found")
			return false, 0
		}
		if !w.matchFeatureWithCondition(features[0], condition) {
			sugar.Debug("matchFeaturesWithCondition() - unmatch")
			return false, 0
		}
		sugar.Debug("matchFeaturesWithCondition() - match")
		return true, 1

	case MecabConditionTypeOneOrNothing: // 0..1に合致
		// 検証要素無し
		if len(features) <= 0 {
			sugar.Debug("matchFeaturesWithCondition() - feature not found")
			return false, 0
		}
		// 0以上なので常にマッチ
		if w.matchFeatureWithCondition(features[0], condition) {
			sugar.Debug("matchFeaturesWithCondition() - match.length: 1")
			return true, 1
		}
		sugar.Debug("matchFeaturesWithCondition() - match.length: 0")
		return true, 0
	case MecabConditionTypeNothingOrContinue: // 0..*に合致
		// 0以上なので常にマッチ
		for i, feature := range features {
			if !w.matchFeatureWithCondition(feature, condition) {
				sugar.Debugf("matchFeaturesWithCondition() - match.length: %d", i)
				return true, i
			}
		}

		sugar.Debugf("matchFeaturesWithCondition() - retain all match.length: %d", len(features))
		return true, len(features)

	case MecabConditionTypeEOS: // EOSに合致
		for _, feature := range features {
			if !w.matchFeatureWithCondition(feature, condition) {

				sugar.Debug("matchFeaturesWithCondition() - unmatch")
				return false, 0
			}
		}
		sugar.Debugf("matchFeaturesWithCondition() - match.length: %d", len(features))
		return true, len(features)
	default:
		// 条件不備
		return false, 0
	}
}

func (w MecabWrapper) matchFeatureWithCondition(feature MecabFeature, condition MecabCondition) bool {
	defer w.Logger.Sync()
	sugar := w.Logger.Sugar()
	sugar.Debugf("matchFeatureWithCondition() - Feature: %s, condition type: %s", feature.String(), condition.ConditionType.String())
	switch condition.ConditionType {
	case MecabConditionTypeOne, MecabConditionTypeOneOrNothing, MecabConditionTypeNothingOrContinue: // 1つに合致, 0..1に合致, 0..*に合致
		for _, conditionFeature := range condition.Features {
			sugar.Debugf("matchFeatureWithCondition() - condition feature: %s", conditionFeature.String())
			if conditionFeature.CheckWord && feature.Word != conditionFeature.Word {
				sugar.Debugf("matchFeatureWithCondition() - unmatch word.feature: %s, condition: %s", feature.Word, conditionFeature.Word)
				continue
			}
			if conditionFeature.CheckWordType && feature.WordType != conditionFeature.WordType {
				sugar.Debugf("matchFeatureWithCondition() - unmatch word type.feature: %s, condition: %s", feature.WordType.String(), conditionFeature.WordType.String())
				continue
			}
			if conditionFeature.CheckWordSubType1 && feature.WordSubType1 != conditionFeature.WordSubType1 {
				sugar.Debugf("matchFeatureWithCondition() - unmatch word sub type1.feature: %s, condition: %s", feature.WordSubType1.String(), conditionFeature.WordSubType1.String())
				continue
			}
			if conditionFeature.CheckOriginalForm && feature.OriginalForm != conditionFeature.OriginalForm {
				sugar.Debugf("matchFeatureWithCondition() - unmatch original form.feature: %s, condition: %s", feature.OriginalForm, conditionFeature.OriginalForm)
				continue
			}
			if conditionFeature.CheckConjugationType && feature.ConjugationType != conditionFeature.ConjugationType {
				sugar.Debugf("matchFeatureWithCondition() - unmatch conjugation type.feature: %s, condition: %s", feature.ConjugationType, conditionFeature.ConjugationType)
				continue
			}
			if conditionFeature.CheckConjugationForm && feature.ConjugationForm != conditionFeature.ConjugationForm {
				sugar.Debugf("matchFeatureWithCondition() - unmatch conjugation form.feature: %s, condition: %s", feature.ConjugationForm, conditionFeature.ConjugationForm)
				continue
			}
			sugar.Debug("matchFeatureWithCondition() - match")
			return true
		}
		sugar.Debug("matchFeatureWithCondition() - unmatch")
		return false

	case MecabConditionTypeEOS: // EOSに合致
		if feature.EOS {
			return true
		}
	}
	return false
}
