#! /bin/sh

# check for $ROOT

if [ -z "$ROOT" ];
then
	echo "\$ROOT is empty"
	exit
fi
exit

# go modules

cd $ROOT
go mod init github.com/SandorMiskey/TEx-Rasa
go mod tidy

# Create virtual python3 environment and activate it

VENV="$ROOT/venv"
python3 -m venv $VENV
source $VENV/bin/activate

# Install packages

$VENV/bin/pip3 install -U pip
$VENV/bin/pip3 install rasa

# self signed tls cert
#
# - gen
# - move to storage/proxy
# - move to storage/engine
