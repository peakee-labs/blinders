resource "aws_instance" "database" {
  ami                    = "ami-072b1c33a2439c226" # ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-arm64-server-20240411
  instance_type          = "t4g.medium"            # 2 CPUs, Arch:arm64, Ram:4GiB
  key_name               = aws_key_pair.tf_ec2_key.key_name
  vpc_security_group_ids = [aws_security_group.ec2_security_group.id]
  depends_on             = [local_file.tf_ec2_key, aws_key_pair.tf_ec2_key]

  root_block_device {
    volume_size = 50
  }

  tags = {
    Name        = "${var.project_name}-database-prod"
    project     = var.project_name,
    environment = "prod"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = tls_private_key.tf_ec2_key.private_key_pem
    }
    script = "../wait_for_instance.sh"
  }

  provisioner "local-exec" {
    command = <<EOT
    ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook ../ansible/ec2_mongodb.ubuntu.ansible.yml \
     -u ubuntu -i '${self.public_ip},' \
     --key-file ./tf_ec2_key.pem \
     --extra-vars 'mongodb_admin_username=${var.mongodb_admin_username} \
      mongodb_admin_password=${var.mongodb_admin_password} \
      arm=true' 
    EOT
  }

  provisioner "local-exec" {
    command = <<EOT
    ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook ../ansible/ec2_redis_stack.ubuntu.ansible.yml \
     -u ubuntu -i '${self.public_ip},' \
     --key-file ./tf_ec2_key.pem \
     --extra-vars 'redis_default_password=${var.redis_default_password}'
    EOT
  }
}


# TODO: need to resolve security group
resource "aws_security_group" "ec2_security_group" {
  name = "${var.project_name}-ec2-security-group-prod"

  # Accept all inbound requests
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "all"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Accept all outbound requests
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "all"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    project     = var.project_name
    environment = "prod"
  }
}

# Create RSA key of size 4096 bits
resource "tls_private_key" "tf_ec2_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "tf_ec2_key" {
  content  = tls_private_key.tf_ec2_key.private_key_pem
  filename = "${path.module}/tf_ec2_key.pem"
}

resource "aws_key_pair" "tf_ec2_key" {
  key_name   = "blinders_prod_tf_ec2_key"
  public_key = tls_private_key.tf_ec2_key.public_key_openssh

  tags = {
    project     = var.project_name
    environment = "prod"
  }
}


output "enable_key_file" {
  value = "chmod 400 ./tf_ec2_key.pem"
}

output "database_public_ip" {
  description = "Public IP address of the database EC2 instance"
  value       = aws_instance.database.public_ip
}


output "ssh_command_to_database" {
  value = "ssh ubuntu@${aws_instance.database.public_ip} -i ./tf_ec2_key.pem"
}


