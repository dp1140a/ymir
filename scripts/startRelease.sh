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

VERSION=$(< VERSION)

# ensure you are on latest develop  & master
git checkout develop
git pull origin develop
git checkout -

git checkout master
git pull origin master
git checkout develop
git flow release start $VERSION

# push released version to server
git push
#git checkout develop
echo -e "\n-------------------------------------------------------------------------------"
echo -e "Start your testing on this branch"
echo -e "NOTE:  Only docs and bug fixes go on this branch until it is finished"
echo -e "All new features should be off of develop"
echo -e "When you are ready to release run scripts/finishRelease.sh"
echo -e "-------------------------------------------------------------------------------"