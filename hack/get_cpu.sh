#!/bin/bash
self=`readlink -f $0`
basedir=`dirname $self`
bindir="$basedir/../bin"

. $basedir/env.sh

ns=$1
pod=$2
echo_proxy "get cpu config with target pod: $ns . ${pod}"
set -x
$bindir/scale-agent kubelet cpu -n $ns -p $pod
