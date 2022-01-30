#!/bin/sh

URLS="
http://localhost:8080/ping
"

for url in $(echo $URLS | tr ";" "\n")
do
    status_code=$(curl --write-out %{http_code} --silent --output /dev/null "$url")

    if [ "$status_code" -ne 200 ] ; then
       echo "Site status changed to $status_code"
    fi
done
