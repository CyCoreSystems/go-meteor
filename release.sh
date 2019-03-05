#!/bin/bash -e

# Only run if we are inside CI
if [ ! -e $CI ]; then
   # Create the release
   echo "Creating release..."
   goreleaser release
fi
