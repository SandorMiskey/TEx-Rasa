#! /bin/bash

# go module



# Create virtual python3 environment and activate it

VENV="$PWD/venv"
python3 -m venv $VENV
source $VENV/bin/activate

# Install packages

$VENV/bin/pip3 install -U pip
$VENV/bin/pip3 install rasa
