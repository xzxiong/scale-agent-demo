#!/bin/bash
self=`readlink -f $0`
basedir=`dirname $self`
bindir="$basedir/../bin"

. $basedir/env.sh

echo_proxy "show key config"
set -x
$bindir/scale-agent kubelet --show
