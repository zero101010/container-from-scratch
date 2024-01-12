resource "aws_instance" "ec2_instance" {
  ami             = var.ami # Ubuntu20.04
  instance_type   = var.instance_type  # Free tier
  subnet_id       = aws_subnet.public_subnet.id
  security_groups = [aws_security_group.instance_sg.id]
  associate_public_ip_address = true
  key_name = "privateIgor"
  tags = {
    Name = var.ec2_name
  }
  root_block_device {
    volume_size = 30  
    volume_type = "gp2"  
  }
  user_data = <<-EOF
              #!/bin/bash
              sudo apt-get update -y
              sudo apt  install golang-go -y
              EOF
}