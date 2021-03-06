#!/bin/bash

N=0

while true; do
    VAL="$(head -c 20 /dev/random | head -n 1)"
    ONE="$(echo -n $VAL | lua test/bin.lua $1)"
    TWO="$(echo -n $VAL | "sha$1sum" | cut -f1 -d ' ')"
    if [ "$ONE" != "$TWO" ]; then
        echo WRONG $ONE
        echo $TWO
    else
        echo ok $ONE $N
        echo ..  $TWO
    fi
    N=$(expr $N + 1)
done
