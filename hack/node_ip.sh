#!/usr/bin/bash
self=`readlink -f $0`
basedir=`dirname $self`
bindir="$basedir/../bin"

. $basedir/env.sh
set -x
$bindir/scale-agent client node
