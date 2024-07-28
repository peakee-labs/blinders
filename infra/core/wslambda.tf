# important:
# - must use runtime "provided.al2" for go lambdas (not provided.al2023)
# - handler must be "bootstrap" for runtime "provided.al2"
# - must use arc "arm64" for go lambdas

resource "aws_lambda_function" "ws_authorizer" {
  function_name    = "${var.project.name}-ws-authorizer-${var.project.environment}"
  filename         = "../../dist/ws_authorizer-${var.project.environment}.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  architectures    = ["arm64"]
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../../dist/ws_authorizer-${var.project.environment}.zip")

  environment {
    variables = {
      ENVIRONMENT : var.project.environment
      MONGO_DATABASE : local.envs.MONGO_DATABASE
      MONGO_DATABASE_URL : local.envs.MONGO_DATABASE_URL
    }
  }

  tags = {
    project     = var.project.name
    environment = var.project.environment
  }
}

resource "aws_lambda_function" "ws_connect" {
  function_name    = "${var.project.name}-ws-connect-${var.project.environment}"
  filename         = "../../dist/connect-${var.project.environment}.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  architectures    = ["arm64"]
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../../dist/connect-${var.project.environment}.zip")

  environment {
    variables = {
      ENVIRONMENT : var.project.environment
      REDIS_HOST : local.envs.REDIS_HOST
      REDIS_PORT : local.envs.REDIS_PORT
      REDIS_USERNAME : local.envs.REDIS_USERNAME
      REDIS_PASSWORD : local.envs.REDIS_PASSWORD
    }
  }

  tags = {
    project     = var.project.name
    environment = var.project.environment
  }
}


resource "aws_lambda_function" "ws_disconnect" {
  function_name    = "${var.project.name}-ws-disconnect-${var.project.environment}"
  filename         = "../../dist/disconnect-${var.project.environment}.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  architectures    = ["arm64"]
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../../dist/disconnect-${var.project.environment}.zip")

  environment {
    variables = {
      ENVIRONMENT : var.project.environment
      REDIS_HOST : local.envs.REDIS_HOST
      REDIS_PORT : local.envs.REDIS_PORT
      REDIS_USERNAME : local.envs.REDIS_USERNAME
      REDIS_PASSWORD : local.envs.REDIS_PASSWORD
    }
  }

  tags = {
    project     = var.project.name
    environment = var.project.environment
  }
}

resource "aws_lambda_function" "ws_chat" {
  function_name    = "${var.project.name}-ws-chat-${var.project.environment}"
  filename         = "../../dist/wschat-${var.project.environment}.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  architectures    = ["arm64"]
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../../dist/wschat-${var.project.environment}.zip")

  environment {
    variables = {
      ENVIRONMENT : var.project.environment
      REDIS_HOST : local.envs.REDIS_HOST
      REDIS_PORT : local.envs.REDIS_PORT
      REDIS_USERNAME : local.envs.REDIS_USERNAME
      REDIS_PASSWORD : local.envs.REDIS_PASSWORD
      MONGO_DATABASE : local.envs.MONGO_DATABASE
      MONGO_DATABASE_URL : local.envs.MONGO_DATABASE_URL
    }
  }

  tags = {
    project     = var.project.name
    environment = var.project.environment
  }
}
