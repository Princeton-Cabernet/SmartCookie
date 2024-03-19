#!/bin/bash

# Check if SDE_INSTALL and SDE environment variables are set
if [ -z "$SDE_INSTALL" ] || [ -z "$SDE" ]; then
    echo "The SDE_INSTALL and SDE environment variables must be set."
    exit 1
fi

# Compile the P4 program
bf-p4c synflood_assist.p4

# Check if the previous command was successful
if [ $? -ne 0 ]; then
    echo "Compilation failed."
    exit 1
fi

# Remove existing configuration file
sudo rm -f "$SDE_INSTALL/share/p4/targets/tofino/synflood_assist.conf"

# Move the new configuration file into place
sudo mv synflood_assist.tofino/synflood_assist.conf "$SDE_INSTALL/share/p4/targets/tofino/synflood_assist.conf"

# Remove any existing directory
sudo rm -rf "$SDE_INSTALL/synflood_assist.tofino"

# Move the new directory into place
sudo mv synflood_assist.tofino "$SDE_INSTALL/"


echo "Compilation complete."
