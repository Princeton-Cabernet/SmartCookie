#!/bin/bash

# Check if SDE_INSTALL and SDE environment variables are set
if [ -z "$SDE_INSTALL" ] || [ -z "$SDE" ]; then
    echo "The SDE_INSTALL and SDE environment variables must be set."
    exit 1
fi

# Compile the P4 program
bf-p4c SmartCookie-HalfSipHash.p4

# Check if the previous command was successful
if [ $? -ne 0 ]; then
    echo "Compilation failed."
    exit 1
fi

# Remove existing configuration file
sudo rm -f "$SDE_INSTALL/share/p4/targets/tofino/SmartCookie-HalfSipHash.conf"

# Move the new configuration file into place
sudo mv SmartCookie-HalfSipHash.tofino/SmartCookie-HalfSipHash.conf "$SDE_INSTALL/share/p4/targets/tofino/SmartCookie-HalfSipHash.conf"

# Remove any existing directory
sudo rm -rf "$SDE_INSTALL/SmartCookie-HalfSipHash.tofino"

# Move the new directory into place
sudo mv SmartCookie-HalfSipHash.tofino "$SDE_INSTALL/"

# Run the switch daemon with the specified program
"$SDE/./run_switchd.sh" -p SmartCookie-HalfSipHash


