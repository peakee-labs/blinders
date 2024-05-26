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


# suggest route
resource "aws_apigatewayv2_integration" "suggest" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_uri  = aws_lambda_function.pysuggest.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_suggest" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /suggest"
  target    = "integrations/${aws_apigatewayv2_integration.suggest.id}"
}

resource "aws_lambda_permission" "suggest" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.pysuggest.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_suggest_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/suggest"
}


# suggest v2 with gosuggest route
resource "aws_apigatewayv2_integration" "suggestv2" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_uri  = aws_lambda_function.gosuggest.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_suggest_v2" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /suggest/v2"
  target    = "integrations/${aws_apigatewayv2_integration.suggestv2.id}"
}

resource "aws_lambda_permission" "suggestv2" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gosuggest.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_suggest_v2_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/suggest/v2"
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

# explore api
resource "aws_apigatewayv2_integration" "explore" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_uri        = aws_lambda_function.explore.invoke_arn
  integration_type       = "AWS_PROXY"
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "explore" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.explore.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_route" "explore" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /explore/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.explore.id}"
}

# explore api
resource "aws_apigatewayv2_integration" "practice" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_uri        = aws_lambda_function.practice.invoke_arn
  integration_type       = "AWS_PROXY"
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "practice" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.practice.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_route" "practice" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /practice/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.practice.id}"
}

resource "aws_apigatewayv2_integration" "collecting-get" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_uri  = aws_lambda_function.collecting-get.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_collecting-get" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /collecting/get"
  target    = "integrations/${aws_apigatewayv2_integration.collecting-get.id}"
}

resource "aws_lambda_permission" "collecting-get" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.collecting-get.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_collecting_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/collecting/get"
}
