data "aws_route53_zone" "blinders" {
  name         = var.aws_route53_zone_name
  private_zone = false
}

data "aws_acm_certificate" "blinders" {
  domain = var.domain_name_for_certificate
}

resource "aws_apigatewayv2_domain_name" "http_api" {
  domain_name = var.domains.http

  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.blinders.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "http_api" {
  name    = aws_apigatewayv2_domain_name.http_api.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.blinders.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.http_api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.http_api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_domain_name" "websocket_api" {
  domain_name = var.domains.websocket

  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.blinders.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "websocket_api" {
  name    = aws_apigatewayv2_domain_name.websocket_api.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.blinders.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.websocket_api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.websocket_api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}


