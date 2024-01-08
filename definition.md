# How Docker works ?
## Chroot
- It's used to create a root dir, isolating the root dir. To do this run the bash commands bellow to create isolate root files:
```
container_dir="container"
echo container_dir
mkdir container_dir && cd container_dir
mkdir bin lib lib64 proc
cp /bin/bash bin && cp /bin/ls bin && cp /bin/ps bin && cp /bin/mount bin 
cp -r /lib/x86_64-linux-gnu/ lib
cp -r /lib64/ld-linux-x86-64.so.2 lib64/
sudo chroot .
mount -t proc proc /proc
```
## Namespaces
- Divide kernel and isolate processes. This allow us to use differents process in differents spaces of kernel. In our case our container will use only the process that we want.
```
# man unshare
sudo unshare -p -f --mount-proc /bin/bash
```


## Cgroups
- Make the isolation of CPU and memory. To combine every part that I explained you can run this command:

```
sudo cgexec -g cpu,memory,pids:/container_dir unshare --pid --uts --mount --fork chroot container_dir
```

- If you need can run the `container.bash` that I create to explain this concepts. To run this bash only need pass the name of the new dir as parameter. For example:

```
bash container.bash igor_container
```

If need more details watch this video: https://www.youtube.com/watch?v=S7Hv2CdNmuA