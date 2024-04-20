resource "aws_instance" "database" {
  ami                    = "ami-02a2af70a66af6dfb"
  instance_type          = "t2.micro"
  key_name               = aws_key_pair.tf_ec2_key.key_name
  vpc_security_group_ids = [aws_security_group.ec2_security_group.id]

  tags = { Name = "${var.project_name}-database-shared" }

  provisioner "local-exec" {
    command = <<EOT
    ansible-playbook ../ec2_mongodb.ansible.yml -u ec2-user -i '${self.public_ip},' \
     --key-file ./tf_ec2_key.pem \
     --extra-vars 'mongodb_admin_username=${var.mongodb_admin_username} \
      mongodb_admin_password=${var.mongodb_admin_password}'
    EOT
  }

  provisioner "local-exec" {
    command = <<EOT
    ansible-playbook ../ec2_redis_stack.ansible.yml -u ec2-user -i '${self.public_ip},' \
     --key-file ./tf_ec2_key.pem \
     --extra-vars 'redis_default_password=${var.redis_default_password}'
    EOT
  }
}


# consider use external embedding api
resource "aws_instance" "services" {
  ami                    = "ami-02a2af70a66af6dfb"
  instance_type          = "t2.micro"
  key_name               = aws_key_pair.tf_ec2_key.key_name
  vpc_security_group_ids = [aws_security_group.ec2_security_group.id]

  tags = merge(
    { Name = "${var.project_name}-services-shared" },
  )
}

# TODO: need to resolve security group
resource "aws_security_group" "ec2_security_group" {
  name = "${var.project_name}-ec2-security-group-shared"

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
  key_name   = "tf_ec2_key"
  public_key = tls_private_key.tf_ec2_key.public_key_openssh
}


output "enable_key_file" {
  value = "chmod 400 ./tf_ec2_key.pem"
}

output "database_public_ip" {
  description = "Public IP address of the database EC2 instance"
  value       = aws_instance.database.public_ip
}


output "ssh_command_to_database" {
  value = "ssh ec2-user@${aws_instance.database.public_ip} -i ./tf_ec2_key.pem"
}

output "services_public_ip" {
  description = "Public IP address of the services EC2 instance"
  value       = aws_instance.services.public_ip
}


output "ssh_command_to_sevices" {
  value = "ssh ec2-user@${aws_instance.services.public_ip} -i ./tf_ec2_key.pem"
}
