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

