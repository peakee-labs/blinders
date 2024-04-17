data "aws_route53_zone" "blinders" {
  name         = "peakee.co"
  private_zone = false
}

resource "aws_acm_certificate" "blinders" {
  domain_name               = "peakee.co"
  subject_alternative_names = ["api.peakee.co", "*.api.peakee.co", "ws.peakee.co", "*.ws.peakee.co"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "validation" {
  for_each = {
    for dvo in aws_acm_certificate.blinders.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.blinders.zone_id
}

resource "aws_acm_certificate_validation" "blinders" {
  certificate_arn         = aws_acm_certificate.blinders.arn
  validation_record_fqdns = [for record in aws_route53_record.validation : record.fqdn]
}

