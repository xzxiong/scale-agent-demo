#!/bin/bash

## const
################
node="10.128.0.124"
#rootfs="/"
rootfs="/rootfs"

## function
################

get_last_pod() {
	local res=`kubectl get pod -A -o wide | grep ${node} | tail -n 1`
	#local content=($res)
	ns=${content[0]}
    pod=${content[1]}
}

echo_proxy() {
	echo "[`date '+%F %T'`] $@"
}


if [ -z "$HOSTNAME" ] ; then
    get_last_pod
    export POD_NAMESPACE="$ns"
    export HOSTNAME="$pod"
fi
export ROOTFS="$rootfs"
