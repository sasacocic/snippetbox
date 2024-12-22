#!/bin/bash



# there seems to be some weirdness around Command and using reflex together
# using bash script is just way easier
# go build ./cmd/run && ./run


reflex --decoration=none -r '\.go$' -s -- bash -c 'go build ./cmd/web/ && ./web'

echo "./run.sh completing"
