#! /bin/bash

# go module

cd cmd
go mod init github.com/SandorMiskey/TEx-Rasa/cmd
go mod tidy

# Create virtual python3 environment and activate it

VENV="$PWD/venv"
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
