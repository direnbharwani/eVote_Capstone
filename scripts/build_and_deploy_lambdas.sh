#!/bin/bash

# =============================================================================
# Build readVote
# =============================================================================

echo "Building readVote..."

cd app/server/readVote

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip readVote.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv readVote.zip ../readVote.zip

# Check if artifact was built from root level
cd ../../../
if test -f app/server/readVote.zip; then
    echo "readVote built!"
else
    echo "failed to build readVote"
fi

# =============================================================================
# Build submitVote
# =============================================================================

# =============================================================================
# Build countVotes
# =============================================================================

# =============================================================================
# Serverless Deployment
# =============================================================================

cd app/server
# sls config credentials --provider aws --key $AWS_ACCESS_KEY_ID --secret $AWS_SECRET_ACCESS_KEY
# sls deploy

# =============================================================================
# Clean Artifacts
# =============================================================================

rm readVote.zip