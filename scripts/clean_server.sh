#!/bin/bash

cd app/server

# =============================================================================
# read-vote
# =============================================================================

rm read-vote.zip

if test -f read-vote.zip; then
    echo "failed to remove read-vote.zip"
else
    echo "successfully removed read-vote.zip!"
fi

# =============================================================================
# count-votes
# =============================================================================

rm count-votes.zip

if test -f count-votes.zip; then
    echo "failed to remove count-votes.zip"
else
    echo "successfully removed count-votes.zip!"
fi

# =============================================================================
# register
# =============================================================================

rm register.zip

if test -f register.zip; then
    echo "failed to remove register.zip"
else
    echo "successfully removed register.zip!"
fi

# =============================================================================
# submit-vote
# =============================================================================

rm submit-vote.zip

if test -f submit-vote.zip; then
    echo "failed to remove submit-vote.zip"
else
    echo "successfully removed submit-vote.zip!"
fi

# =============================================================================
# create-election
# =============================================================================

rm create-election.zip

if test -f create-election.zip; then
    echo "failed to remove create-election.zip"
else
    echo "successfully removed create-election.zip!"
fi

# =============================================================================
# get-election
# =============================================================================

rm get-election.zip

if test -f get-election.zip; then
    echo "failed to remove get-election.zip"
else
    echo "successfully removed get-election.zip!"
fi


# =============================================================================
# Back to root
cd ../../../