#!/bin/sh
set -euo pipefail

## ref, https://cloud.tencent.com/developer/article/1839411
#############################################################
## exampel result
## $ go mod download -json k8s.io/api@kubernetes-1.28.11
# {
# 	"Path": "k8s.io/api",
# 	"Version": "v0.28.11",
# 	"Query": "kubernetes-1.28.11",
# 	"Info": "/Users/jacksonxie/go/pkg/mod/cache/download/k8s.io/api/@v/v0.28.11.info",
# 	"GoMod": "/Users/jacksonxie/go/pkg/mod/cache/download/k8s.io/api/@v/v0.28.11.mod",
# 	"Zip": "/Users/jacksonxie/go/pkg/mod/cache/download/k8s.io/api/@v/v0.28.11.zip",
# 	"Dir": "/Users/jacksonxie/go/pkg/mod/k8s.io/api@v0.28.11",
# 	"Sum": "h1:2qFr3jSpjy/9QirmlRP0LZeomexuwyRlE8CWUn9hPNY=",
# 	"GoModSum": "h1:nQSGyxQ2sbS73i1zEJyaktFvFfD72z/7nU+LqxzNnXk="
# }

VERSION=${1#"v"}
if [ -z "$VERSION" ]; then
    echo "Must specify version!"
    echo "like: $0 v1.28.11"
    exit 1
fi
MODS=($(
    curl -sS https://raw.githubusercontent.com/kubernetes/kubernetes/v${VERSION}/go.mod |
    sed -n 's|.*k8s.io/\(.*\) => ./staging/src/k8s.io/.*|k8s.io/\1|p'
))
total=${#MODS[@]}
idx=1
for MOD in "${MODS[@]}"; do
    V=$(
        go mod download -json "${MOD}@kubernetes-${VERSION}" |
        sed -n 's|.*"Version": "\(.*\)".*|\1|p'
    )
    echo "get ${idx}/${total} ${MOD}@${V}"
    let idx=${idx}+1
    go mod edit "-replace=${MOD}=${MOD}@${V}"
done
go get "k8s.io/kubernetes@v${VERSION}"

echo "ALL Done."
