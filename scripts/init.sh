#!/usr/bin/bash
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
REPO_NAME=$(basename $WD)

#Initialize Git
echo -e "\nInitializing Git"
git init

unset CREATE_REPO
until [[ $CREATE_REPO =~ [YyNn] ]] ; do
  read -p "Do you want to create a GitHub repo for this project [N/y]: " CREATE_REPO
  CREATE_REPO=${CREATE_REPO:-N}
done

if [[ $CREATE_REPO =~ [yY] ]]
then
  echo -e "Creating Repo"
  echo -e "Repo Name: $REPO_NAME"
  RESULT=$(gh repo create $REPO_NAME --public)
  wait
  echo -e "Created Repo $RESULT"
  git remote add origin git@github.com:dp1140a/$REPO_NAME.git
  wait
  git remote -v
fi

#Initialize GitFlow
echo -e "\nInitializing Git Flow"
git flow init

#Initialize semver
echo -e "\nInitializing semver"
$(which semver) init
wait

#Initialize go mod
go mod init github.com/dp1140a/$REPO_NAME

if [[ $CREATE_REPO =~ [yY] ]]
then
  #Initial Commit
  git checkout develop
  git add --all
  git commit -m "Project Initialization"
  git push origin develop
fi


echo -e "\n********************************************************************************************************"
echo -e "                                  Project Initialization complete"
echo -e "********************************************************************************************************"