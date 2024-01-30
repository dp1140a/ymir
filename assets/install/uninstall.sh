#!/usr/bin/bash

# Determine OS name
OS=$(uname)
VERSION=${PWD##*/}
APP_DIR=/usr/lib/ymir
CURRENT_DIR=$APP_DIR/current
INSTALL_DIR=$APP_DIR/$VERSION

unset REMOVE_VERSION
until [[ $REMOVE_VERSION =~ [YyNn] ]] ; do
  echo -e "You can choose to uninstall only the latest version or all version of ymir."
  echo -e "Uninstalling only the current version will rollback to the previous installed version."
  echo -e "Removing all version will uninstall everything"
  echo -e "NOTE: Either choice WILL NOT DELETE the db, model or printer files!"
  read -p "Do you want to remove only the current version [Y/n]: " REMOVE_VERSION
  REMOVE_VERSION=${REMOVE_VERSION:-Y}
done
if [[ $REMOVE_VERSION =~ [yY] ]]
then
  #Remove only the current version
  rm -rf $CURRENT_DIR
  echo -e "This may break the app link.  If it does cd to the last version installed and run the following:"
  echo -e "sudo ln -s /usr/lib/ymir/${PWD##*/} /usr/lib/ymir/current"

else
  #Remove everything
  systemctl stop ymir
  systemctl disable ymir
  systemctl daemon-reload
  rm -rf /var/log/ymir
  rm -rf $APP_DIR
  rm /usr/bin/ymir
  rm $HOME/.ymir/ymir.toml
  rm $HOME/.local/share/applications/ymir.desktop
  rm $HOME/.local/share/icons/ymir*.png
fi

echo "Done."