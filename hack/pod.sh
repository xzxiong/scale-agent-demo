#!/bin/bash
self=`readlink -f $0`
basedir=`dirname $self`
bindir="$basedir/../bin"

. $basedir/env.conf

echo_proxy "list all pods in current node"
set -x
$bindir/scale-agent client pod
