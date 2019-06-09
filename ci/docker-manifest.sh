#!/usr/bin/env bash

set -eux

token="$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:$TRAVIS_REPO_SLUG:pull" | jq -r '.token')"
allTags="$(curl -Ls "https://index.docker.io/v2/$TRAVIS_REPO_SLUG/tags/list" -H "Authorization: Bearer $token" | jq '.tags')"

versions=("${TRAVIS_TAG}" "latest")
for version in "${versions[@]}"; do
    tags=($(echo "$allTags" | jq -r --arg version "$version" '.[] | select(. | startswith($version))'))
    manifests=()
    for tag in "${tags[@]}"; do
        arch=${tag#"$version-"}
        manifests+=("$TRAVIS_REPO_SLUG:$tag")
#        echo "tag=$tag; version=$version; arch=$arch"
    done
    docker manifest create "$TRAVIS_REPO_SLUG:$version" "${manifests[@]}"
done
