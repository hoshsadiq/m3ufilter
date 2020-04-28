#!/usr/bin/env bash

# todo need to also do v1, v1.1, v1.1.1 etc

docker tag "$TRAVIS_REPO_SLUG:${TRAVIS_TAG}-$ARCH" "$TRAVIS_REPO_SLUG:latest-$ARCH"

docker push "$TRAVIS_REPO_SLUG:${TRAVIS_TAG}-$ARCH"
docker push "$TRAVIS_REPO_SLUG:latest-$ARCH"

if [[ "$ARCH" == "amd64" ]]; then
    docker tag "$TRAVIS_REPO_SLUG:${TRAVIS_TAG}-$ARCH" "$TRAVIS_REPO_SLUG:latest"
    docker push "$TRAVIS_REPO_SLUG:latest"
fi
