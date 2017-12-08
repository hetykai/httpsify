#!/bin/sh

if [[ "$OSTYPE" == "linux-gnu" ]]; then
	
elif [[ "$OSTYPE" == "darwin"* ]]; then

else
    echo "Unsupported OS, try to build it from source"
fi
