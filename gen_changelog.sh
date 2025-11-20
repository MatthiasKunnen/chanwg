#!/usr/bin/env bash

set -euo pipefail

root_dir=$(git rev-parse --show-toplevel)
docker run -t \
	--mount type=bind,src="${root_dir}/.git",dst=/app/.git,readonly \
	--mount type=bind,src="${root_dir}/cliff.toml",dst=/app/cliff.toml,readonly \
	--mount type=bind,src="${root_dir}/CHANGELOG.md",dst=/app/CHANGELOG.md \
	"orhunp/git-cliff:${TAG:-latest}" --latest --prepend CHANGELOG.md "$@"

docker run -t \
	--mount type=bind,src="${root_dir}/.git",dst=/app/.git,readonly \
	--mount type=bind,src="${root_dir}/cliff.toml",dst=/app/cliff.toml,readonly \
	--mount type=bind,src="${root_dir}/CHANGELOG.md",dst=/app/CHANGELOG.md \
	"orhunp/git-cliff:${TAG:-latest}" --unreleased --prepend CHANGELOG.md "$@"
