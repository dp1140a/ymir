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
cp bin/ymir $INSTALL_DIR

# create symlink /usr/lib/ymir/current --> current version of Ymir
ln -s $INSTALL_DIR $CURRENT_DIR

#Create symlink in /usr/bin to above
ln -s $CURRENT_DIR/ymir /usr/bin/ymir

# Install config file to /etc/ymir
mkdir -p /etc/ymir
cp ymir.toml /etc/ymir

# Make log dir
mkdir -p /var/log/ymir

cp README.md LICENSE $HOME/.ymir

# Install .desktop file for ui in ~/.local/share/applications
cp lib/ymir-ui.desktop $HOME/.local/share/applications/ymir-ui.desktop

# install icon in ~/.local/share/icons
cp lib/ymir*.png $HOME/.local/share/icons/

# Install .service file in /lib/systemd/system/
cp lib/ymir.service /etc/systemd/system/

echo "It is recommended that you edit the config file at /etc/ymir/ymir.toml before starting the service."
echo "Once you have done that and are ready to start the service run the following:"
echo "systemctl start ymir"
echo "systemctl enable ymir"
echo "systemctl daemon-reload"

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