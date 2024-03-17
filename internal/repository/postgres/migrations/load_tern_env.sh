#!/bin/bash

# Load the .env file
export $(egrep -v '^#' ../../../../.env | xargs)

# Replace the placeholders in the tern.conf file
envsubst < tern.conf.template > tern.conf
