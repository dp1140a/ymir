#!/usr/bin/bash

# Determine OS name
OS=$(uname)
VERSION=${PWD##*/}
APP_DIR=/usr/lib/ymir
CURRENT_DIR=$APP_DIR/current
INSTALL_DIR=$APP_DIR/$VERSION
# Install git
if [ "$OS" = "Linux" ]; then
echo "This is a Linux Machine"
#Install binary, README and LICENSE in /usr/lib/ymir/[version]
mkdir -p $INSTALL_DIR
cp {README.md, LICENSE, ymir} $INSTALL_DIR

# create symlink /usr/lib/ymir/current --> current version of Ymir
ln -s $INSTALL_DIR $CURRENT_DIR

#Create symlink in /usr/bin to above
ln -s /usr/bin/ymir $CURRENT_DIR/bin/ymir

# Install config file to ~/.ymir
mkdir -p $HOME/.ymir
cp ymir.toml $HOME/.ymir

# Install .desktop file for ui in ~/.local/share/applications
cp $CURRENT_DIR/lib/ymir.desktop $HOME/.local/share/applications/

# install icon in ~/.local/share/icons
cp $CURRENT_DIR/lib/ymir*.png $HOME/.local/share/icons/

# Install .service file in /lib/systemd/system/
cp $CURRENT_DIR/lib/ymir.service /lib/systemd/system/
systemctl start ymir
systemctl enable ymir
systemctl daemon-reload

else
  echo "Unsupported OS"
  exit 1
fi

#test bin install
echo -e "Checking ymir installation"
if ymir version >/dev/null 2>&1; then
  echo "ymir is configured correctly."
else
  echo "ymir install test failed. Please check the installation."
  exit 1
fi

# Test the service
echo "Checking service installation"
if systemctl is-active --quiet service;
then
  echo -e "Ymir service installed"
else
  echo -e "Oops.  Something happened installing ymir as a service.  Check the logs and rerun."
  exit 1
fi

echo "Congratulations, Ymir has now successfully been installed!"