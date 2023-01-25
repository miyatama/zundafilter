package zunda_mecab

import (
	"os"
	"fmt"
	"path/filepath"
	"go.uber.org/zap"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_NAME = "data/zunda.db"
)

type ZundaDbRepository struct {
}
type ConvertVerbConjugationRow struct {
	BaseWord string
	Mizen    string
}

func (z ZundaDbRepository) SelectConvertVerbConjugationTable(baseWord string) (ConvertVerbConjugationRow, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Debug("ZundaDbRepository#SelectConvertVerbConjugationTable()")

	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir := filepath.Dir(path)
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s", binDir, DB_NAME))
	if err != nil {
		return ConvertVerbConjugationRow{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT mizen FROM ConvertVerbConjugationTable WHERE base_word = ?")
	if err != nil {
		return ConvertVerbConjugationRow{}, err
	}
	defer stmt.Close()

	var mizen string
	err = stmt.QueryRow(baseWord).Scan(&mizen)
	if err != nil {
		return ConvertVerbConjugationRow{}, err
	}

	sugar.Debugf("ZundaDbRepository#SelectConvertVerbConjugationTable() - return (%v, %v)", baseWord, mizen)
	return ConvertVerbConjugationRow{
		BaseWord: baseWord,
		Mizen:    mizen,
	}, err
}
