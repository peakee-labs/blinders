
resource "aws_apigatewayv2_api" "http_api" {
  name          = "blinders-http-api"
  protocol_type = "HTTP"

  # for lambda-proxy integration, cors is enabled by backend
  # cors_configuration {
  #   allow_origins = ["https://app.peakee.co", "http://localhost:3000"]
  # }
}

resource "aws_apigatewayv2_stage" "http_staging" {
  name        = "staging"
  api_id      = aws_apigatewayv2_api.http_api.id
  auto_deploy = true
}

resource "aws_apigatewayv2_api_mapping" "http_api_v1" {
  api_id          = aws_apigatewayv2_api.http_api.id
  domain_name     = aws_apigatewayv2_domain_name.http_api.id
  stage           = aws_apigatewayv2_stage.http_staging.id
  api_mapping_key = "v1"
}

resource "aws_apigatewayv2_api" "websocket_api" {
  name                       = "blinders-websocket-api"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "websocket_staging" {
  name        = "staging"
  api_id      = aws_apigatewayv2_api.websocket_api.id
  auto_deploy = true
}

resource "aws_apigatewayv2_api_mapping" "websocket_api_v1" {
  api_id          = aws_apigatewayv2_api.websocket_api.id
  domain_name     = aws_apigatewayv2_domain_name.websocket_api.id
  stage           = aws_apigatewayv2_stage.websocket_staging.id
  api_mapping_key = "v1"
}

resource "aws_apigatewayv2_authorizer" "websocket_authorizer" {
  name             = "blinders-websocket-authorizer"
  api_id           = aws_apigatewayv2_api.websocket_api.id
  authorizer_type  = "REQUEST"
  authorizer_uri   = aws_lambda_function.ws_authorizer.invoke_arn
  identity_sources = ["route.request.querystring.token"]
}

output "http-api-endpoint" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}"
}

output "websocket-api-endpoint" {
  value = "wss://${aws_apigatewayv2_api_mapping.websocket_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.websocket_api_v1.api_mapping_key}"
}
