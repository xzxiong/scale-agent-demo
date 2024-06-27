#!/bin/sh
set -euo pipefail

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

while read MOD VERSION; do
    V=$(
        go mod download -json "${MOD}@${VERSION}" |
        sed -n 's|.*"Version": "\(.*\)".*|\1|p'
    )
    #echo "go get ${MOD}@${V}"
    echo "cmd: go mod edit \"-replace=${MOD}=${MOD}@${V}\""
done << EOF
k8s.io/cri-client kubernetes-1.28.0-alpha.4
k8s.io/cri-client kubernetes-1.31.0-alpha.1
k8s.io/legacy-cloud-providers kubernetes-1.30.0-rc.0
EOF

echo "ALL Done."
