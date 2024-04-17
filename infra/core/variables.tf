variable "region" {
  type    = string
  default = "ap-south-1"
}

variable "profile" {
  type    = string
  default = "admin.peakee"
}

variable "project" {
  type = object({
    name           = string
    environment    = string
    default_region = string
  })

}

variable "domains" {
  type = object({
    http      = string
    websocket = string
  })
}

variable "env_filename" {
  type = string
}

locals {
  envs = { for tuple in regexall("(.*)=(.*)", file(var.env_filename)) : tuple[0] => sensitive(tuple[1]) }
}

variable "aws_route53_zone_name" {
  type    = string
  default = "peakee.co"
}

variable "domain_name_for_certificate" {
  type    = string
  default = "peakee.co"
}
