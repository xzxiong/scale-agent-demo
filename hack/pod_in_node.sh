#!/bin/bash
self=`readlink -f $0`
basedir=`dirname $self`
bindir="$basedir/../bin"

. $basedir/env.sh

node=$1
echo_proxy "list all pods in current node: '$node'"
set -x
$bindir/scale-agent client pod -n $node
