#!/bin/bash
##################################################
# Variables
##################################################
PROGNAME=$(basename $0)
VERSION="1.0"
VERB_FILE="Verb.csv"
ADJ_FILE="Adj.csv"
IMPORT_TEMP_CSV="temp.csv"
DB_NAME="zunda.db"

##################################################
# functions
##################################################
function usage() {
  echo "Usage: ${PROGNAME} [OPTIONS]"
  echo "  need sqlite3."
  echo "Options:"
  echo "  -h, --help"
  echo "  -v, --version"
}

function outputInfo() {
  echo "[INFO] ${1}"
}

function cleanupBeforeRun() {
  outputInfo "cleanupBeforeRun()"
  removeCsvfiles
  if [ -e ${DB_NAME} ] ; then
    rm -f ${DB_NAME}
  fi
  removeImportTempCsv
}

function cleanupAfterRun() {
  outputInfo "cleanupAfterRun()"
  removeCsvfiles
  removeImportTempCsv
}

function removeCsvfiles() {
  rm -f ./*.csv
  rm -f ./*.csv.eucjp
}

function copyAndConvert() {
  outputInfo "copyAndConvert()"
  mecab_dictionary_version="2.7.0-20070801"
  curl -fsSL 'https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7MWVlSDBCSXZMTXM' -o mecab-ipadic-${mecab_dictionary_version}.tar.gz
  tar zxfv mecab-ipadic-${mecab_dictionary_version}.tar.gz
  rm mecab-ipadic-${mecab_dictionary_version}.tar.gz
  files=("${VERB_FILE}" "${ADJ_FILE}")
  for file in "${files[@]}" ; do
    cp "./mecab-ipadic-${mecab_dictionary_version}/${file}" ./${file}.eucjp
    iconv -f EUCJP -t UTF8 ./${file}.eucjp > ./${file}
  done
  rm -fr mecab-ipadic-${mecab_dictionary_version}
}

function createDatabase() {
  outputInfo "createDatabase()"
  echo ".open ${DB_NAME}" | sqlite3
}

function createTables() {
  outputInfo "createTables()"
  # 動詞の未然形変換テーブル
  createVerbConvertConjugationTable
}

function createVerbConvertConjugationTable() {
  SQL=`cat << EOF
CREATE TABLE ConvertVerbConjugationTable(
  base_word TEXT, 
  mizen TEXT
);
EOF
`
  echo "${SQL}" | sqlite3 ${DB_NAME}
}

function importData() {
  outputInfo "importData()"
  # 動詞の未然形変換テーブル
  importVerbConvertConjugationTable
}

function importVerbConvertConjugationTable() {
  removeImportTempCsv

  cat ${VERB_FILE} | \
    grep ',未然形,' | \
    awk -F, '{printf("%s,%s\n", $11, $1)}' >> ./${IMPORT_TEMP_CSV}

  sqlite3 -separator , ${DB_NAME} ".import ./${IMPORT_TEMP_CSV} ConvertVerbConjugationTable"
}

function removeImportTempCsv() {
  if [ -e ${IMPORT_TEMP_CSV} ] ; then
    rm -f ${IMPORT_TEMP_CSV}
  fi
}

function main() {
  cleanupBeforeRun
  copyAndConvert
  createDatabase
  createTables
  importData
  cleanupAfterRun
}

##################################################
# Main
##################################################
for OPT in "${@}"
do
  case ${OPT} in
    -h | --help)
      usage
      exit 1
      ;;
    -v | --version)
      echo ${VERSION}
      exit 1
      ;;
    -- | -)
      shift 1
      param+=( "$@" )
      break
      ;;
    -*)
      echo "${PROGNAME}: illegal option -- '$(echo $1 | sed 's/^-*//')'" 1>$2
      exit 1
      ;;
    *)
      if [ [ ! -z "$1"] && [ ! "$1" =~ ^-+ ] ] ; then
        param+=( "$1" )
        shift 1
      fi
      ;;
  esac
done

# sqlite3 exists check
exists=`type sqlite3 > /dev/null; echo ${?}`
if [ ${exists} -ne 0 ] ; then
  echo "need sqlite3 command.please install sqlite3."
  exit 1
fi
main
