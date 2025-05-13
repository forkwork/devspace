#!/bin/zsh

export NUM_WORKSPACES=60

# Start the workspaces
for i in $(seq 1 $NUM_WORKSPACES);
do
    devspace delete --force "loadtest$i" &
    sleep 2
done

wait
