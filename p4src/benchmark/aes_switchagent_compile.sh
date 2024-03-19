#!/bin/bash

# Check if SDE_INSTALL and SDE environment variables are set
if [ -z "$SDE_INSTALL" ] || [ -z "$SDE" ]; then
    echo "The SDE_INSTALL and SDE environment variables must be set."
    exit 1
fi

# Compile the P4 program
bf-p4c SmartCookie-AES.p4

# Check if the previous command was successful
if [ $? -ne 0 ]; then
    echo "Compilation failed."
    exit 1
fi

# Remove existing configuration file
sudo rm -f "$SDE_INSTALL/share/p4/targets/tofino/SmartCookie-AES.conf"

# Move the new configuration file into place
sudo mv SmartCookie-AES.tofino/SmartCookie-AES.conf "$SDE_INSTALL/share/p4/targets/tofino/SmartCookie-AES.conf"

# Remove any existing directory
sudo rm -rf "$SDE_INSTALL/SmartCookie-AES.tofino"

# Move the new directory into place
sudo mv SmartCookie-AES.tofino "$SDE_INSTALL/"


echo "Compilation complete."
