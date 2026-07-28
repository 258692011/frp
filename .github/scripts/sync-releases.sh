#!/usr/bin/env bash
set -euo pipefail

: "${GH_TOKEN:?GH_TOKEN is required}"
: "${UPSTREAM_REPO:?UPSTREAM_REPO is required}"
: "${FORK_REPO:?FORK_REPO is required}"

WORK_DIR="${RUNNER_TEMP:-/tmp}/release-sync"
rm -rf "$WORK_DIR"
mkdir -p "$WORK_DIR"

mapfile -t RELEASES < <(
  gh release list --repo "$UPSTREAM_REPO" --limit 100 \
    --json tagName,name,isLatest,isPrerelease,createdAt \
    --jq 'sort_by(.createdAt) | .[] | @base64'
)

for encoded in "${RELEASES[@]}"; do
  release_json=$(printf '%s' "$encoded" | base64 -d)
  tag=$(jq -r '.tagName' <<< "$release_json")
  name=$(jq -r '.name // .tagName' <<< "$release_json")
  is_latest=$(jq -r '.isLatest' <<< "$release_json")
  is_prerelease=$(jq -r '.isPrerelease' <<< "$release_json")

  if ! gh release view "$tag" --repo "$FORK_REPO" >/dev/null 2>&1; then
    echo "Creating release $tag"
    notes=$(gh release view "$tag" --repo "$UPSTREAM_REPO" --json body --jq '.body // ""')
    flags=(--latest=false)
    [[ "$is_latest" == "true" ]] && flags=(--latest)
    [[ "$is_prerelease" == "true" ]] && flags+=(--prerelease)

    gh release create "$tag" --repo "$FORK_REPO" --title "$name" \
      --notes "${notes}

---
Auto-synced from https://github.com/${UPSTREAM_REPO}" \
      "${flags[@]}"
  else
    echo "Release $tag already exists"
  fi

  fork_assets=$(gh release view "$tag" --repo "$FORK_REPO" \
    --json assets --jq '[.assets[].name]')
  mapfile -t upstream_assets < <(
    gh release view "$tag" --repo "$UPSTREAM_REPO" \
      --json assets --jq '.assets[].name'
  )

  asset_dir="$WORK_DIR/${tag//\//_}"
  mkdir -p "$asset_dir"

  for asset in "${upstream_assets[@]}"; do
    if jq -e --arg name "$asset" 'index($name) != null' <<< "$fork_assets" >/dev/null; then
      echo "Asset $tag/$asset already exists"
      continue
    fi

    echo "Copying asset $tag/$asset"
    gh release download "$tag" --repo "$UPSTREAM_REPO" \
      --pattern "$asset" --dir "$asset_dir" --clobber
    gh release upload "$tag" "$asset_dir/$asset" --repo "$FORK_REPO"
  done
done

echo "Release synchronization complete"
