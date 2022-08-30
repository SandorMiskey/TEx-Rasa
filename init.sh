#! /bin/bash

VENV="$PWD/venv"

# Create virtual python3 environment and activate it

python3 -m venv $VENV
source $VENV/bin/activate

# Install packages

$VENV/bin/pip3 install -U pip
$VENV/bin/pip3 install rasa

# init

# rasa init

