# My own container
## How Containers works ?
### Chroot
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
### Namespaces
- Divide kernel and isolate processes. This allow us to use differents process in differents spaces of kernel. In our case our container will use only the process that we want.
```
# man unshare
sudo unshare -p -f --mount-proc /bin/bash
```


### Cgroups
- Make the isolation of CPU and memory. To combine every part that I explained you can run this command:

```
sudo cgexec -g cpu,memory,pids:/container_dir unshare --pid --uts --mount --fork chroot container_dir
```

## Terraform 
- If you don't use a linux environment as me, I created terraform directory to create a simple ec2 instance in AWS to help you to run this labor. It's only necessary be logged in AWS cli, if you don't know how to do this follow this [ link](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-quickstart.html).
- After configure aws cli create a Certificate in AWS console to use to make ssh connections. After create go to `terraform/terraform.tfvars` and change the `key_name` value to the name that was created for you. After this it's only enter inside of `terraform` directory and run this commands:
```
terraform init
terraform plan
terraform apply
```
## Create the container with Bash
- Observation, bash script only work in linux environment
- To run this bash only need pass the name of the new dir as parameter. For example, in this case I choose the name `igor_container`, but could be any name that you want:

```
bash container.bash igor_container
```
- If need more details about how the container works watch this video: https://www.youtube.com/watch?v=S7Hv2CdNmuA
## Create the container with Golang
- Observation, go script only works in linux environment
- This application was a example to explain how works the container. To do This I understand the concept of that running in bash script and isolating the process using namespace, creating CHROOT and the CGROUP. If do you want to check how this work only with bash script fell free to check `container.bash`.
- I love golang, and I thought that would be a good chance to use my learned about container and use golang.
- It's important pay attention about the concepts from creating of containers. Understand about how  containers works is the goal of this project.
- To run our code em golang run the command:
```
sudo go run main.go run /bin/sh
```

- If do you want to check how this works it's only necessary enter in our github [actions](https://github.com/zero101010/container-from-scratch/actions/runs/7523306008/job/20476500493) that have the CI's that show the difference about run command `ps` inside of our container and outside. The steps are called `Run ps in our OS without container` and `Run PS inside of the container`

