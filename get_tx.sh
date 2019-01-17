#!/bin/bash

BASE_URL=https://eos.greymass.com

TX_ID=1

if [ -n "$1" ]; then
    TX_ID=$1
fi

if [ -n "$2" ]; then
    BASE_URL=$2
fi

URL=$BASE_URL/v1/history/get_transaction

echo 'TX_ID:'$TX_ID
echo 'URL:'$URL

#curl -i -X POST -H "'Content-type':'application/json'" -d '{"ATime":"'$atime'","BTime":"'$btime'"}' $url
curl -X POST -H "'Content-type':'application/json'" -d '{ "id" : "'$TX_ID'" }' --url $URL
