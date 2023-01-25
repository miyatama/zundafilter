package zunda_mecab

import (
	"fmt"
	"github.com/bluele/mecab-golang"
	"strings"
)

const BOSEOS = "BOS/EOS"

type MecabFeature struct {
	EOS             bool
	Word            string               // 表層形
	WordType        MecabWordType        // 品詞
	WordSubType1    MecabWordSubType1    // 品詞細分類1
	WordSubType2    MecabWordSubType2    // 品詞細分類2
	WordSubType3    MecabWordSubType3    // 品詞細分類3
	ConjugationType MecabConjugationType // 活用型
	ConjugationForm MecabConjugationForm // 活用形
	OriginalForm    string               // 原形
	// 読み
	// 発音
}

type MecabConditionFeature struct {
	CheckWord            bool
	Word                 string
	CheckWordType        bool
	WordType             MecabWordType
	CheckWordSubType1    bool
	WordSubType1         MecabWordSubType1
	CheckOriginalForm    bool
	OriginalForm         string
	CheckConjugationForm bool
	ConjugationForm      MecabConjugationForm
	CheckConjugationType bool
	ConjugationType      MecabConjugationType
}

func (m *MecabConditionFeature) String() string {
	descriptions := []string{}
	if m.CheckWord {
		descriptions = append(descriptions, fmt.Sprintf("単語: %s", m.Word))
	}
	if m.CheckWordType {
		descriptions = append(descriptions, fmt.Sprintf("品詞タイプ: %s", m.WordType.String()))
	}
	if m.CheckOriginalForm {
		descriptions = append(descriptions, fmt.Sprintf("基本形: %s", m.OriginalForm))
	}
	return fmt.Sprintf("{%s}", strings.Join(descriptions, ","))
}

type MecabCondition struct {
	ConditionType MecabConditionType
	Features      []MecabConditionFeature
}

func (m *MecabCondition) String() string {

	descriptions := []string{}
	for _, feature := range m.Features {
		descriptions = append(descriptions, feature.String())
	}
	return fmt.Sprintf("条件: %s, 単語群: [%v]", m.ConditionType.String(), strings.Join(descriptions, ","))

}

// 条件タイプ
type MecabConditionType int

const (
	MecabConditionTypeOne               MecabConditionType = iota // 1つに合致
	MecabConditionTypeOneOrNothing                                // 0..1に合致
	MecabConditionTypeNothingOrContinue                           // 0..*に合致
	MecabConditionTypeEOS                                         // EOSに合致
)

func (m MecabConditionType) String() string {
	switch m {
	case MecabConditionTypeOne:
		return "1つに合致"
	case MecabConditionTypeOneOrNothing:
		return "0..1に合致"
	case MecabConditionTypeNothingOrContinue:
		return "0..*に合致"
	case MecabConditionTypeEOS:
		return "EOSに合致"
	default:
		return "未知"
	}
}

func (m *MecabFeature) String() string {
	if m.EOS {
		return "EOS"
	}
	return fmt.Sprintf("単語: %s(%s) 分類: [%s, %s, %s, %s], 活用: [%s, %s]",
		m.Word,
		m.OriginalForm,
		m.WordType.String(),
		m.WordSubType1.String(),
		m.WordSubType2.String(),
		m.WordSubType3.String(),
		m.ConjugationType.String(),
		m.ConjugationForm.String(),
	)
}

func parseMecabFeatureNode(node *mecab.Node) MecabFeature {
	features := strings.Split(node.Feature(), ",")
	eos := features[0] == BOSEOS
	if eos {
		return MecabFeature{
			EOS: eos,
		}
	}
	return MecabFeature{
		EOS:             eos,
		Word:            node.Surface(),
		WordType:        parseMecabWordType(features[0]),
		WordSubType1:    parseMecabWordSubType1(features[1]),
		WordSubType2:    parseMecabWordSubType2(features[2]),
		WordSubType3:    parseMecabWordSubType3(features[3]),
		ConjugationType: parseMecabConjugationType(features[4]),
		ConjugationForm: parseMecabConjugationForm(features[5]),
		OriginalForm:    features[6],
	}
}

/*
* 品詞タイプのパース
 */
func parseMecabWordType(keyword string) MecabWordType {
	for _, wordType := range []MecabWordType{
		MecabWordTypeUnknown,
		MecabWordTypeNoun,
		MecabWordTypeParticle,
		MecabWordTypeVerb,
		MecabWordTypeAuxiliaryVerb,
		MecabWordTypeAdjective,
		MecabWordTypeAdverb,
		MecabWordTypeSymbol,
		MecabWordTypeAttributive,
		MecabWordTypeConjunction,
	} {
		if wordType.String() == keyword {
			return wordType
		}
	}
	return MecabWordTypeUnknown
}

/*
* 品詞細分類1のパース
 */
func parseMecabWordSubType1(keyword string) MecabWordSubType1 {
	for _, wordSubType := range []MecabWordSubType1{
		MecabWordSubType1None,
		MecabWordSubType1Independence,
		MecabWordSubType1NotIndependence,
		MecabWordSubType1AdverbConnect,
		MecabWordSubType1ParticleCombination,
		MecabWordSubType1ParticleSub,
		MecabWordSubType1ParticleConnect,
		MecabWordSubType1ParticleRank,
		MecabWordSubType1ParticleTail,
		MecabWordSubType1VerbTail,
		MecabWordSubType1NounPopuler,
		MecabWordSubType1NounPronoun,
		MecabWordSubType1NounAdjective,
		MecabWordSubType1SymbolPeriod,
		MecabWordSubType1SymbolComma,
		MecabWordSubType1NounAdjectiveNai,
	} {
		if wordSubType.String() == keyword {
			return wordSubType
		}
	}
	return MecabWordSubType1None
}

/*
* 品詞細分類2のパース
 */
func parseMecabWordSubType2(keyword string) MecabWordSubType2 {
	for _, wordSubType := range []MecabWordSubType2{
		MecabWordSubType2None,
		MecabWordSubType2Popular,
		MecabWordSubType2NumberClassifier,
	} {
		if wordSubType.String() == keyword {
			return wordSubType
		}
	}
	return MecabWordSubType2None
}

/*
* 品詞細分類3のパース
 */
func parseMecabWordSubType3(keyword string) MecabWordSubType3 {
	return MecabWordSubType3None
}

/*
* 活用型のパース
 */
func parseMecabConjugationType(keyword string) MecabConjugationType {
	for _, conjugationType := range []MecabConjugationType{
		MecabConjugationTypeNone,
		MecabConjugationTypeSpTa,
		MecabConjugationTypeSpDa,
		MecabConjugationTypeSpDesu,
		MecabConjugationTypeSahenSuru,
		MecabConjugationTypeIchidan,
		MecabConjugationTypeGodanNa,
		MecabConjugationTypeGodanBa,
		MecabConjugationTypeGodanMa,
		MecabConjugationTypeGodanTa,
		MecabConjugationTypeGodanWaSokuOnbin,
		MecabConjugationTypeGodanRa,
		MecabConjugationTypeGodanRaSp,
		MecabConjugationTypeGodanKaIOnbin,
		MecabConjugationTypeGodanGa,
		MecabConjugationTypeI,
	} {
		if conjugationType.String() == keyword {
			return conjugationType
		}
	}
	return MecabConjugationTypeNone
}

/*
* 活用形のパース
 */
func parseMecabConjugationForm(keyword string) MecabConjugationForm {
	for _, conjugationForm := range []MecabConjugationForm{
		MecabConjugationFormNone,
		MecabConjugationFormKihon,
		MecabConjugationFormTaigenSetsuzoku,
		MecabConjugationFormKatei,
		MecabConjugationFormKateiShukuYaku1,
		MecabConjugationFormTaigenSeatsuzokuSp,
		MecabConjugationFormTaigenSeatsuzokuSp2,
		MecabConjugationFormMeireiE,
		MecabConjugationFormMeireiI,
		MecabConjugationFormMeireiRo,
		MecabConjugationFormMeireiYo,
		MecabConjugationFormBungoKihon,
		MecabConjugationFormMizen,
		MecabConjugationFormMizenUSetsuzoku,
		MecabConjugationFormMizenNuSetsuzoku,
		MecabConjugationFormMizenReruSetsuzoku,
		MecabConjugationFormMizenSp,
		MecabConjugationFormGendaiKihon,
		MecabConjugationFormRenyou,
		MecabConjugationFormRenyouTaSetsuzoku,
	} {
		if conjugationForm.String() == keyword {
			return conjugationForm
		}
	}
	return MecabConjugationFormNone
}

// 品詞タイプ
type MecabWordType int

const (
	MecabWordTypeUnknown       MecabWordType = iota
	MecabWordTypeNoun                        // 名詞
	MecabWordTypeParticle                    // 助詞
	MecabWordTypeVerb                        // 動詞
	MecabWordTypeAuxiliaryVerb               // 助動詞
	MecabWordTypeAdjective                   // 形容詞
	MecabWordTypeAdverb                      // 副詞
	MecabWordTypeSymbol                      // 記号
	MecabWordTypeAttributive                 // 連体詞
	MecabWordTypeConjunction                 // 接続詞
)

func (m MecabWordType) String() string {
	switch m {
	case MecabWordTypeNoun:
		return "名詞"
	case MecabWordTypeParticle:
		return "助詞"
	case MecabWordTypeVerb:
		return "動詞"
	case MecabWordTypeAuxiliaryVerb:
		return "助動詞"
	case MecabWordTypeAdjective:
		return "形容詞"
	case MecabWordTypeAdverb:
		return "副詞"
	case MecabWordTypeSymbol:
		return "記号"
	case MecabWordTypeAttributive:
		return "連体詞"
	case MecabWordTypeConjunction:
		return "接続詞"
	default:
		return "未知"
	}
}

// 品詞細分類1
type MecabWordSubType1 int

const (
	// 共通
	MecabWordSubType1None            MecabWordSubType1 = iota // なし
	MecabWordSubType1Independence                             // 自立
	MecabWordSubType1NotIndependence                          // 非自立
	// 副詞
	MecabWordSubType1AdverbConnect // 助詞類接続
	// 助動詞
	// 助詞
	MecabWordSubType1ParticleCombination // 係助詞
	MecabWordSubType1ParticleSub         // 副助詞
	MecabWordSubType1ParticleConnect     // 接続助詞
	MecabWordSubType1ParticleRank        // 格助詞
	MecabWordSubType1ParticleTail        // 終助詞
	// 動詞
	MecabWordSubType1VerbTail // 接尾
	// 名詞
	MecabWordSubType1NounPopuler   // 一般
	MecabWordSubType1NounPronoun   // 代名詞
	MecabWordSubType1NounAdjective // 形容動詞語幹
	MecabWordSubType1NounAdjectiveNai // ナイ形容詞語幹
	// 形容詞
	// 記号
	MecabWordSubType1SymbolPeriod // 句点
	MecabWordSubType1SymbolComma  // 読点
	// 連体詞
)

func (m MecabWordSubType1) String() string {
	switch m {
	case MecabWordSubType1None:
		return "なし"
	case MecabWordSubType1Independence:
		return "自立"
	case MecabWordSubType1NotIndependence:
		return "非自立"
	case MecabWordSubType1AdverbConnect:
		return "助詞類接続"
	case MecabWordSubType1ParticleCombination:
		return "係助詞"
	case MecabWordSubType1ParticleSub:
		return "副助詞"
	case MecabWordSubType1ParticleConnect:
		return "接続助詞"
	case MecabWordSubType1ParticleRank:
		return "格助詞"
	case MecabWordSubType1ParticleTail:
		return "終助詞"
	case MecabWordSubType1VerbTail:
		return "接尾"
	case MecabWordSubType1NounPopuler:
		return "一般"
	case MecabWordSubType1NounPronoun:
		return "代名詞"
	case MecabWordSubType1NounAdjective:
		return "形容動詞語幹"
	case MecabWordSubType1SymbolPeriod:
		return "句点"
	case MecabWordSubType1SymbolComma:
		return "読点"
	case MecabWordSubType1NounAdjectiveNai :
		return "ナイ形容詞語幹"
	default:
		return "未知"
	}
}

// 品詞細分類2
type MecabWordSubType2 int

const (
	// 共通
	MecabWordSubType2None             MecabWordSubType2 = iota // なし
	MecabWordSubType2Popular                                   // 一般
	MecabWordSubType2NumberClassifier                          // 助数詞
)

func (m MecabWordSubType2) String() string {
	switch m {
	case MecabWordSubType2None:
		return "*"
	case MecabWordSubType2Popular:
		return "一般"
	case MecabWordSubType2NumberClassifier:
		return "助数詞"
	default:
		return "未知"
	}
}

// 品詞細分類3
type MecabWordSubType3 int

const (
	// 共通
	MecabWordSubType3None MecabWordSubType3 = iota // なし
)

func (m MecabWordSubType3) String() string {
	return "*"
}

// 活用型
type MecabConjugationType int

const (
	// 共通
	MecabConjugationTypeNone MecabConjugationType = iota // なし
	// 副詞
	// 助動詞
	MecabConjugationTypeSpTa      // 特殊・タ
	MecabConjugationTypeSpDa      // 特殊・ダ
	MecabConjugationTypeSpDesu    // 特殊・デス
	MecabConjugationTypeSahenSuru // サ変・スル
	// 助詞
	// 動詞
	MecabConjugationTypeIchidan          // 一段
	MecabConjugationTypeGodanNa          // 五段・ナ行
	MecabConjugationTypeGodanBa          // 五段・バ行
	MecabConjugationTypeGodanMa          // 五段・マ行
	MecabConjugationTypeGodanTa          // 五段・タ行
	MecabConjugationTypeGodanWaSokuOnbin // 五段・ワ行促音便
	MecabConjugationTypeGodanRa          // 五段・ラ行
	MecabConjugationTypeGodanRaSp        // 五段・ラ行特殊
	MecabConjugationTypeGodanKaIOnbin    // 五段・カ行イ音便
	MecabConjugationTypeGodanGa          // 五段・ガ行
	// 名詞
	// 形容詞
	MecabConjugationTypeI // 形容詞・イ段
	// 記号
	// 連体詞
)

func (m MecabConjugationType) String() string {
	switch m {
	case MecabConjugationTypeNone:
		return "*"
	case MecabConjugationTypeSpTa:
		return "特殊・タ"
	case MecabConjugationTypeSpDa:
		return "特殊・ダ"
	case MecabConjugationTypeSpDesu:
		return "特殊・デス"
	case MecabConjugationTypeSahenSuru:
		return "サ変・スル"
	case MecabConjugationTypeIchidan:
		return "一段"
	case MecabConjugationTypeGodanNa:
		return "五段・ナ行"
	case MecabConjugationTypeGodanBa:
		return "五段・バ行"
	case MecabConjugationTypeGodanMa:
		return "五段・マ行"
	case MecabConjugationTypeGodanTa:
		return "五段・タ行"
	case MecabConjugationTypeGodanWaSokuOnbin:
		return "五段・ワ行促音便"
	case MecabConjugationTypeGodanRa:
		return "五段・ラ行"
	case MecabConjugationTypeGodanRaSp:
		return "五段・ラ行特殊"
	case MecabConjugationTypeGodanKaIOnbin:
		return "五段・カ行イ音便"
	case MecabConjugationTypeGodanGa:
		return "五段・ガ行"
	case MecabConjugationTypeI:
		return "形容詞・イ段"
	default:
		return "未知"
	}
}

// 活用形
type MecabConjugationForm int

const (
	// 共通
	MecabConjugationFormNone  MecabConjugationForm = iota // なし
	MecabConjugationFormKihon                             // 基本形
	// 副詞
	// 助動詞
	MecabConjugationFormTaigenSetsuzoku // 体言接続
	// 助詞
	// 動詞
	MecabConjugationFormKatei               // 仮定形
	MecabConjugationFormKateiShukuYaku1     // 仮定縮約１
	MecabConjugationFormTaigenSeatsuzokuSp  // 体言接続特殊
	MecabConjugationFormTaigenSeatsuzokuSp2 // 体言接続特殊２
	MecabConjugationFormMeireiE             // 命令ｅ
	MecabConjugationFormMeireiI             // 命令ｉ
	MecabConjugationFormMeireiRo            // 命令ｒｏ
	MecabConjugationFormMeireiYo            // 命令ｙｏ
	MecabConjugationFormBungoKihon          // 文語基本形
	MecabConjugationFormMizen               // 未然形
	MecabConjugationFormMizenUSetsuzoku     // 未然ウ接続
	MecabConjugationFormMizenNuSetsuzoku    // 未然ヌ接続
	MecabConjugationFormMizenReruSetsuzoku  // 未然レル接続
	MecabConjugationFormMizenSp             // 未然特殊
	MecabConjugationFormGendaiKihon         // 現代基本形
	MecabConjugationFormRenyou              // 連用形
	MecabConjugationFormRenyouTaSetsuzoku   // 連用タ接続
	// 名詞
	// 形容詞
	// 記号
	// 連体詞
)

func (m MecabConjugationForm) IsMeirei() bool {
	if m == MecabConjugationFormMeireiE || // 命令ｅ
		m == MecabConjugationFormMeireiI || // 命令ｉ
		m == MecabConjugationFormMeireiRo || // 命令ｒｏ
		m == MecabConjugationFormMeireiYo { // 命令ｙｏ
		return true
	}
	return false
}
func (m MecabConjugationForm) String() string {
	switch m {
	case MecabConjugationFormNone:
		return "*"
	case MecabConjugationFormKihon:
		return "基本形"
	case MecabConjugationFormTaigenSetsuzoku:
		return "体言接続"
	case MecabConjugationFormKatei:
		return "仮定形"
	case MecabConjugationFormKateiShukuYaku1:
		return "仮定縮約１"
	case MecabConjugationFormTaigenSeatsuzokuSp:
		return "体言接続特殊"
	case MecabConjugationFormTaigenSeatsuzokuSp2:
		return "体言接続特殊２"
	case MecabConjugationFormMeireiE:
		return "命令ｅ"
	case MecabConjugationFormMeireiI:
		return "命令ｉ"
	case MecabConjugationFormMeireiRo:
		return "命令ｒｏ"
	case MecabConjugationFormMeireiYo:
		return "命令ｙｏ"
	case MecabConjugationFormBungoKihon:
		return "文語基本形"
	case MecabConjugationFormMizen:
		return "未然形"
	case MecabConjugationFormMizenUSetsuzoku:
		return "未然ウ接続"
	case MecabConjugationFormMizenNuSetsuzoku:
		return "未然ヌ接続"
	case MecabConjugationFormMizenReruSetsuzoku:
		return "未然レル接続"
	case MecabConjugationFormMizenSp:
		return "未然特殊"
	case MecabConjugationFormGendaiKihon:
		return "現代基本形"
	case MecabConjugationFormRenyou:
		return "連用形"
	case MecabConjugationFormRenyouTaSetsuzoku:
		return "連用タ接続"
	default:
		return "未知"
	}
}
