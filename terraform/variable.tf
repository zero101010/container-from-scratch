## EC2 Variables
variable "ami"{
    default = "ami-0c7217cdde317cfec"
}
variable "instance_type"{
    default = "t2.micro"
}
variable "ec2_name"{
    default = "go-instance"
}
variable "volume_size"{
    default = "30"
}

## Security Group(sg) Variables
variable "ports_allow"{
    default = [
        22,8000,3000
    ]
}
