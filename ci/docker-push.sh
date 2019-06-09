#!/usr/bin/env bash

docker tag $TRAVIS_REPO_SLUG:${TRAVIS_TAG:-v0.0.0}-$ARCH $TRAVIS_REPO_SLUG:latest-$ARCH

docker push $TRAVIS_REPO_SLUG:${TRAVIS_TAG:-v0.0.0}-$ARCH
docker push $TRAVIS_REPO_SLUG:latest-$ARCH
