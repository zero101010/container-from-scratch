#!/bin/bash
## Create the CHROOT, Need to improve to use in linux env
## I tried in mac env and didn't work as I expected
export container_dir="$1"
echo $container_dir
mkdir $container_dir && cd $container_dir
mkdir bin lib lib64 proc
## Copy bins that will be using inside of our container
cp /bin/bash bin && cp /bin/ls bin && cp /bin/ps bin && cp /bin/mount bin 
## Copy libs that will be used from binaries of our container
cp -r /lib/x86_64-linux-gnu/ lib
cp -r /lib64/ld-linux-x86-64.so.2 lib64/
## Create CGROUPS of memory, cpu and pids
sudo mkdir /sys/fs/cgroup/cpu/$container_dir
sudo mkdir /sys/fs/cgroup/memory/$container_dir 
sudo mkdir /sys/fs/cgroup/pids/$container_dir
pwd
### Create the container setting the CHROOT Created, namespaces to isolate the processes and the CGROUP
sudo cgexec -g cpu,memory,pids:/$container_dir unshare --pid --uts --mount --fork chroot .