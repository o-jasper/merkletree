#!/bin/bash
#
# Runs the contract.

# If argument given, that is the contract, otherwise
# from file.
CONTRACT="$1"
if [ "$1" == "" ]; then
    CONTRACT=$(cat contract_address)
fi

fetch_state()
{
    epm query $CONTRACT 0 > state
    I=32
    while [ "$CUR" != "0x" ]; do
        CUR=$(epm query $CONTRACT $I)
        echo $CUR >> state
        I=$(expr $I + 32)
    done
}

if [ "$2" == "query" ]; then
    fetch_state $CONTRACT
    cat state
    exit
fi

GOT=/tmp/got
../../path_chunk_n_root > $GOT  # Check a merkle root
HEXDATA=$(cat $GOT | tail -n +1 |head -n 2 | tr -d '\n')
echo epm transact $CONTRACT $HEXDATA
epm transact $CONTRACT $HEXDATA

tail -n 1 $GOT >> expect

fetch_state $CONTRACT
