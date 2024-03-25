# dictionary route
resource "aws_apigatewayv2_integration" "dictionary" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_uri  = aws_lambda_function.dictionary.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_dictionary" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /dictionary"
  target    = "integrations/${aws_apigatewayv2_integration.dictionary.id}"
}

resource "aws_lambda_permission" "dictionary" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dictionary.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_dictionary_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/dictionary"
}

# translate route
resource "aws_apigatewayv2_integration" "translate" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_uri  = aws_lambda_function.translate.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_translate" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /translate"
  target    = "integrations/${aws_apigatewayv2_integration.translate.id}"
}

resource "aws_lambda_permission" "translate" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.translate.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_translate_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/translate"
}

# ws connect
resource "aws_apigatewayv2_integration" "ws_connect" {
  api_id           = aws_apigatewayv2_api.websocket_api.id
  integration_uri  = aws_lambda_function.ws_connect.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "ws_connect" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  route_key          = "$connect"
  target             = "integrations/${aws_apigatewayv2_integration.ws_connect.id}"
  authorization_type = "CUSTOM"
  authorizer_id      = aws_apigatewayv2_authorizer.websocket_authorizer.id
}


# grant invoke lambda permission to api gateway (init trigger for lambda)
resource "aws_lambda_permission" "ws_connect" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ws_connect.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}

# ws disconnect
resource "aws_apigatewayv2_integration" "ws_disconnect" {
  api_id           = aws_apigatewayv2_api.websocket_api.id
  integration_uri  = aws_lambda_function.ws_disconnect.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "ws_disconnect" {
  api_id    = aws_apigatewayv2_api.websocket_api.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.ws_disconnect.id}"
}

resource "aws_lambda_permission" "ws_disconnect" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ws_disconnect.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}


# ws chat
resource "aws_apigatewayv2_integration" "ws_chat" {
  api_id           = aws_apigatewayv2_api.websocket_api.id
  integration_uri  = aws_lambda_function.ws_chat.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "ws_chat" {
  api_id    = aws_apigatewayv2_api.websocket_api.id
  route_key = "chat"
  target    = "integrations/${aws_apigatewayv2_integration.ws_chat.id}"
}

resource "aws_lambda_permission" "ws_chat" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ws_chat.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}


# authorizer
# grant invoke lambda permission to api gateway (init trigger for lambda)
resource "aws_lambda_permission" "ws_authorizer" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ws_authorizer.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}


# rest api
resource "aws_apigatewayv2_integration" "rest" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_uri        = aws_lambda_function.rest.invoke_arn
  integration_type       = "AWS_PROXY"
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "rest" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rest.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_route" "rest" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.rest.id}"
}


resource "aws_apigatewayv2_route" "root" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.rest.id}"
}

output "rest_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/<users|...>"
}
