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

if test -f ../read-vote.zip; then
    echo "read-vote built!"
else
    echo "failed to build read-vote"
fi

# =============================================================================
# Build count-votes
# =============================================================================

echo "Building count-votes..."

cd ../count-votes

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip count-votes.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv count-votes.zip ../count-votes.zip

# Check if artifact was built from root level
if test -f ../count-votes.zip; then
    echo "count-votes built!"
else
    echo "failed to build count-votes"
fi

# =============================================================================
# Build register
# =============================================================================

echo "Building register..."

cd ../register

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip register.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv register.zip ../register.zip

# Check if artifact was built from root level
if test -f ../register.zip; then
    echo "register built!"
else
    echo "failed to build register"
fi

# =============================================================================
# Build submit-vote
# =============================================================================

echo "Building submit-vote..."

cd ../submit-vote

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip submit-vote.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv submit-vote.zip ../submit-vote.zip

# Check if artifact was built from root level
if test -f ../submit-vote.zip; then
    echo "submit-vote built!"
else
    echo "failed to build submit-vote"
fi

# =============================================================================
# Build get-election
# =============================================================================

echo "Building get-election..."

cd ../get-election

# build go binary
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go

# zip as build artifact for serverless deployment
zip get-election.zip bootstrap

# delete built binary & move readVote.zip to root level for deployment
rm bootstrap
mv get-election.zip ../get-election.zip

# Check if artifact was built from root level
if test -f ../get-election.zip; then
    echo "get-election built!"
else
    echo "failed to build get-election"
fi


# =============================================================================
# Back to root
cd ../../../