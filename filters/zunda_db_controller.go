package filters

import (
	"zundafilter/zunda_mecab"
)

type ZundaDbController interface {
	SelectConvertVerbConjugationTable(baseWord string) (zunda_mecab.ConvertVerbConjugationRow, error)
}
