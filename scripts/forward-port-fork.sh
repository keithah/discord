#!/usr/bin/env bash
set -euo pipefail

upstream_remote="${UPSTREAM_REMOTE:-upstream}"
upstream_url="${UPSTREAM_URL:-https://github.com/mautrix/discord.git}"
upstream_branch="${UPSTREAM_BRANCH:-main}"
target_branch="${TARGET_BRANCH:-main}"
sync_branch="${SYNC_BRANCH:-sync/upstream-main}"

if ! git remote get-url "${upstream_remote}" >/dev/null 2>&1; then
	git remote add "${upstream_remote}" "${upstream_url}"
fi

git fetch origin "${target_branch}"
git fetch "${upstream_remote}" --tags --prune

upstream_ref="${upstream_remote}/${upstream_branch}"
target_ref="origin/${target_branch}"
base="$(git merge-base "${target_ref}" "${upstream_ref}")"
mapfile -t patches < <(git rev-list --reverse "${base}..${target_ref}")

if [[ "${#patches[@]}" -eq 0 ]]; then
	echo "No fork-only commits found on ${target_ref}; nothing to forward-port."
	exit 0
fi

git checkout -B "${sync_branch}" "${upstream_ref}"

for commit in "${patches[@]}"; do
	if git cherry-pick -x "${commit}"; then
		continue
	fi

	if git diff --quiet && git diff --cached --quiet; then
		git cherry-pick --skip
		continue
	fi

	echo "Conflict while cherry-picking $(git rev-parse --short "${commit}")" >&2
	echo "Resolve the conflict, run 'git cherry-pick --continue', then continue the stack manually." >&2
	exit 1
done

echo "Forward-ported ${#patches[@]} fork-only commits onto ${upstream_ref}."
echo "Review with: git log --oneline --decorate ${upstream_ref}..HEAD"
echo "Verify with: go test ./..."
