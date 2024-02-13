
resource "aws_apigatewayv2_api" "blinders" {
  name          = "blinders-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "staging" {
  name        = "staging"
  api_id      = aws_apigatewayv2_api.blinders.id
  auto_deploy = true
}

resource "aws_apigatewayv2_api_mapping" "blinders_v1" {
  api_id          = aws_apigatewayv2_api.blinders.id
  domain_name     = aws_apigatewayv2_domain_name.blinders.id
  stage           = aws_apigatewayv2_stage.staging.id
  api_mapping_key = "v1"
}

resource "aws_apigatewayv2_api" "blinders_websocket" {
  name                       = "blinders-websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_api_mapping" "blinders_websocket_v1" {
  api_id          = aws_apigatewayv2_api.blinders_websocket.id
  domain_name     = aws_apigatewayv2_domain_name.blinders.id
  stage           = aws_apigatewayv2_stage.staging.id
  api_mapping_key = "v1"
}

output "custom_domain_api_v1" {
  value = "https://${aws_apigatewayv2_api_mapping.blinders_v1.domain_name}/${aws_apigatewayv2_api_mapping.blinders_v1.api_mapping_key}"
}
