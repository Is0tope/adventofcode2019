#!/bin/bash
mkdir $1
mkdir $1/a
cp template.go $1/a/$1a.go
mkdir $1/b
touch $1/b/$1b.go
touch $1/input.txt
touch $1/test.txt