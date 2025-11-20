#!/usr/bin/env bash

set -euo pipefail

touch cliff.toml
touch CHANGELOG.md

root_dir=$(git rev-parse --show-toplevel)
docker run -t \
	--mount type=bind,src="${root_dir}/.git",dst=/app/.git,readonly \
	--mount type=bind,src="${root_dir}/cliff.toml",dst=/app/cliff.toml,readonly \
	--mount type=bind,src="${root_dir}/CHANGELOG.md",dst=/app/CHANGELOG.md \
	"orhunp/git-cliff:${TAG:-latest}" -o CHANGELOG.md "$@"