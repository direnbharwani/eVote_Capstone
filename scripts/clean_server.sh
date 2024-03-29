#!/bin/bash

# =============================================================================
# Clean Artifacts
# =============================================================================

cd app/server

# =============================================================================
# Read Vote
# =============================================================================

rm read-vote.zip

if test -f read-vote.zip; then
    echo "failed to remove read-vote.zip"
else
    echo "successfully removed read-vote.zip!"
fi

# =============================================================================
# Read Vote
# =============================================================================

rm count-votes.zip

if test -f count-votes.zip; then
    echo "failed to remove count-votes.zip"
else
    echo "successfully removed count-votes.zip!"
fi

# =============================================================================
# Read Vote
# =============================================================================

rm register.zip

if test -f register.zip; then
    echo "failed to remove register.zip"
else
    echo "successfully removed register.zip!"
fi

# Back to root
cd ../../../