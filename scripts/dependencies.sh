#!/bin/bash

which uncover &> /dev/null
if [ "$?" -ne "0" ]; then
	go install github.com/gregoryv/uncover/cmd/uncover@v0.9.0
fi
