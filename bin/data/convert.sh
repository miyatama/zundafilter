#!/bin/bash

rm -f convert_result.txt
while read line; do
	echo ${line} | ../zundafilter >> ./convert_result.txt
	echo "" >> ./convert_result.txt
done < "source.txt"

cat ./convert_result.txt
