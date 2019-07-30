#!/usr/bin/env bash

set -euo pipefail

# todo utilise the API to get all the tags and match up the correct ones on the same line, finally, combine
#   below is some initial work done on this
#
# latest, latest-386, latest-amd64, latest-arm, latest-arm64, v0-beta, v0-beta-386, v0-beta-amd64, v0-beta-arm, v0-beta-arm64, v0.1-beta, v0.1-beta-386, v0.1-beta-amd64, v0.1-beta-arm, v0.1-beta-arm64, v0.1.0-beta, v0.1.0-beta-386, v0.1.0-beta-amd64, v0.1.0-beta-arm, v0.1.0-beta-arm64
# stable, stable-386, stable-amd64, stable-arm, stable-arm64, v0, v0-386, v0-amd64, v0-arm, v0-arm64, v0.1, v0.1-386, v0.1-amd64, v0.1-arm, v0.1-arm64, v0.1.0, v0.1.0-386, v0.1.0-amd64, v0.1.0-arm, v0.1.0-arm64
#get_tags() {
#    # get authorization token
#    token="$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:$TRAVIS_REPO_SLUG:pull" | jq -r '.token')"
#
#    # find all tags
#    tags=($(curl -s -H "Authorization: Bearer $TOKEN" "https://index.docker.io/v2/$TRAVIS_REPO_SLUG/tags/list" | jq -r '.tags[]'))
#
#    # for each tags
#    for tag in "${tags[@]}"; do
#      digest=$(curl -fsSI -H "Authorization: Bearer $TOKEN" -H "Accept: application/vnd.docker.distribution.manifest.v2+json" "https://index.docker.io/v2/$TRAVIS_REPO_SLUG/manifests/$tag" | awk -F: '/^Docker-Content-Digest/{print $2}')
#
#        echo "$tag $digest"
#    done
#}

push_readme() {
  declare -r image="${1}"
  shift
  text=""
  for file in "$@"; do
    text="$text$(<"$file")"
  done

  token=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${DOCKER_USERNAME}'", "password": "'${DOCKER_PASSWORD}'"}' https://cloud.docker.com/v2/users/login/ | jq -r .token)

  local code=$(jq -n --arg msg "$text" \
    '{"registry":"registry-1.docker.io","full_description": $msg }' | \
        curl -s -o /dev/null  -L -w "%{http_code}" \
           https://cloud.docker.com/v2/repositories/"${image}"/ \
           -d @- -X PATCH \
           -H "Content-Type: application/json" \
           -H "Authorization: JWT ${token}")

  if [[ "${code}" = "200" ]]; then
    printf "Successfully pushed README to Docker Hub"
  else
    printf "Unable to push README to Docker Hub, response code: %s\n" "${code}"
    exit 1
  fi
}

push_readme "$@"
