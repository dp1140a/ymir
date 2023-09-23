#!/usr/bin/env bash
set -e
CUR_DIR=`pwd`
WD=`pwd`

if [[ "$WD" = */scripts ]];
then
  echo "Not in Project root.  Changing WD to project root"
  cd ..
  WD=`pwd`

else
  echo "In Project Root"
fi
echo "Working Dir: $WD"

rm -rf .git
rm VERSION
rm go.mod
gh repo delete $(basename $WD) --yes

echo -e "DONE"