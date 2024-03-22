#!/bin/bash

# =============================================================================
# Build readVote
# =============================================================================

echo "Building read-vote..."

cd app/server/read-vote

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip read-vote.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv read-vote.zip ../read-vote.zip

# Check if artifact was built from root level
cd ../../../
if test -f app/server/read-vote.zip; then
    echo "read-vote built!"
else
    echo "failed to build read-vote"
fi

# =============================================================================
# Build submitVote
# =============================================================================

# =============================================================================
# Build count-votes
# =============================================================================

echo "Building count-votes..."

cd app/server/count-votes

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip count-votes.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv count-votes.zip ../count-votes.zip

# Check if artifact was built from root level
cd ../../../
if test -f app/server/count-votes.zip; then
    echo "count-votes built!"
else
    echo "failed to build count-votes"
fi

# =============================================================================
# Serverless Deployment
# =============================================================================

cd app/server
sls config credentials --provider aws --key $AWS_ACCESS_KEY_ID --secret $AWS_SECRET_ACCESS_KEY
sls deploy

# =============================================================================
# Clean Artifacts
# =============================================================================

rm read-vote.zip
rm count-votes.zip