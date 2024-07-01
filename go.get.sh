#!/bin/bash
#
## for v0.28.0

go_get_update() {
    local pkg=$1
    local ver=$2
    echo "go get -u $pkg@$ver"
    go get -u $pkg@$ver
}
go_get() {
    local pkg=$1
    local ver=$2
    echo "go get $pkg@$ver"
    go get $pkg@$ver
}

## for v0.28.0
go_get_update github.com/google/cel-go v0.16.1
go_get github.com/google/cel-go/parser v0.16.1

go_get github.com/golang/protobuf/proto v1.5.4
go_get github.com/prometheus/common v0.44.0

go_get_update github.com/opencontainers/runc v1.1.7
go_get_update k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9

go_get k8s.io/kubernetes v1.28.4

## for github.com/opencontainers/runc v1.1.7
go_get github.com/cilium/ebpf v0.7.0
