#!/usr/bin/env bash

docker tag $TRAVIS_REPO_SLUG:${TRAVIS_TAG}-$ARCH $TRAVIS_REPO_SLUG:latest-$ARCH

docker push $TRAVIS_REPO_SLUG:${TRAVIS_TAG}-$ARCH
docker push $TRAVIS_REPO_SLUG:latest-$ARCH
