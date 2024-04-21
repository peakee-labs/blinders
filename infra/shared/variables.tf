variable "region" {
  type    = string
  default = "ap-south-1"
}

variable "profile" {
  type    = string
  default = "admin.peakee"
}

variable "project_name" {
  type    = string
  default = "blinders"
}

variable "aws_route53_zone_name" {
  type    = string
  default = "peakee.co"
}

variable "domain_name_for_certificate" {
  type    = string
  default = "peakee.co"
}

variable "subject_alternative_names_for_certificate" {
  type    = list(string)
  default = ["api.peakee.co", "*.api.peakee.co", "ws.peakee.co", "*.ws.peakee.co"]
}


variable "mongodb_admin_username" {
  type      = string
  sensitive = false
}

variable "mongodb_admin_password" {
  type      = string
  sensitive = true
}

variable "redis_default_password" {
  type      = string
  sensitive = true
}
