#!/bin/sh

find . -name "go.mod" | xargs -n 1 dirname | xargs -I{} sh -c 'cd "{}" && go mod tidy'
