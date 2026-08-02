#!/usr/bin/env bash
set -euo pipefail

readonly BASE_CHART="oci://ghcr.io/sjoon99/charts/cells"
readonly BASE_VERSION="0.1.5"
readonly BASE_DIGEST="sha256:fdfc0ba1c32145edb4e6f6bfa4452b727d15e89ced402fbf506e9f78adac6a64"

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
output_dir="${1:-${script_dir}/dist}"
build_dir="$(mktemp -d)"
trap 'rm -rf -- "${build_dir}"' EXIT

pull_output="$(helm pull "${BASE_CHART}" \
  --version "${BASE_VERSION}" \
  --untar \
  --untardir "${build_dir}" 2>&1)"
printf '%s\n' "${pull_output}"
grep -Fq "Digest: ${BASE_DIGEST}" <<<"${pull_output}"

sed -i -E 's/^  frontendPassword: .+$/  frontendPassword: ""/' \
  "${build_dir}/cells/values.yaml"
sed -i -E 's@^    frontendpassword: .*@    frontendpassword: {{ .Values.install.frontendPassword | quote }}@' \
  "${build_dir}/cells/templates/configmap.yaml"
rm -- "${build_dir}/cells/LOCAL_CHART_PATCH_NOTES.md"

patch \
  --directory "${build_dir}/cells" \
  --fuzz 0 \
  --strip 1 \
  < "${script_dir}/cells-0.1.6.patch"

helm lint "${build_dir}/cells"
mkdir -p -- "${output_dir}"
helm package "${build_dir}/cells" --destination "${output_dir}"
