#!/usr/bin/env bash

[[ -z "$TAURON_USER" ]] && { echo "Error: TAURON_USER is not set"; exit 1; }
[[ -z "$TAURON_PASS" ]] && { echo "Error: TAURON_PASS is not set"; exit 1; }

COOKIE_JAR=cookie.jar.txt
UA="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
REF="https://elicznik.tauron-dystrybucja.pl/energia"

curl -b $COOKIE_JAR -c $COOKIE_JAR https://logowanie.tauron-dystrybucja.pl/login -L -o /dev/null -H "User-Agent: $UA" -H "Referer: $REF"
curl -b $COOKIE_JAR -c $COOKIE_JAR https://logowanie.tauron-dystrybucja.pl/login -H "User-Agent: $UA" -F username=$TAURON_USER -F password=$TAURON_PASS -F service=https://elicznik.tauron-dystrybucja.pl -H "Referer: $REF" -L -o /dev/null

curl -b $COOKIE_JAR -c $COOKIE_JAR https://elicznik.tauron-dystrybucja.pl/energia/do/dane -H "User-Agent: $UA" -H "Referer: $REF" \
    -F "form[from]=15.10.2021" \
    -F "form[to]=21.10.2021" \
    -F "form[type]=godzin" \
    -F "form[consum]=1" \
    -F "form[oze]=1" \
    -F "form[fileType]=CSV" \
    -o sample.csv

curl -b $COOKIE_JAR -c $COOKIE_JAR https://elicznik.tauron-dystrybucja.pl/energia/do/dane -H "User-Agent: $UA" -H "Referer: $REF" \
    -F "form[from]=30.10.2022" \
    -F "form[to]=30.10.2022" \
    -F "form[type]=godzin" \
    -F "form[consum]=1" \
    -F "form[oze]=1" \
    -F "form[fileType]=CSV" \
    -o sample-CEST-to-CET.csv

curl -b $COOKIE_JAR -c $COOKIE_JAR https://elicznik.tauron-dystrybucja.pl/energia/do/dane -H "User-Agent: $UA" -H "Referer: $REF" \
    -F "form[from]=27.03.2022" \
    -F "form[to]=27.03.2022" \
    -F "form[type]=godzin" \
    -F "form[consum]=1" \
    -F "form[oze]=1" \
    -F "form[fileType]=CSV" \
    -o sample-CET-to-CEST.csv

cat sample.csv | grep -v "2021-10-15 2:00" > sample-missing-data-on-day1.csv
cat sample.csv | grep -v "2021-10-18 14:00" > sample-missing-data-on-day4.csv
cat sample.csv | grep -v "2021-10-21 20:00" > sample-missing-data-on-day7.csv
cat sample.csv | grep -v "2021-10-15" > sample-missing-day1.csv
cat sample.csv | grep -v "2021-10-18" > sample-missing-day4.csv
cat sample.csv | grep -v "2021-10-21" > sample-missing-day7.csv
cat sample.csv | egrep -v "2021-10-16 2:00;.*;pobór;" > sample-different-array-sizes.csv

