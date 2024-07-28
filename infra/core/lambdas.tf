# resource "aws_lambda_function" "rest" {
#   function_name    = "${var.project.name}-rest-api-${var.project.environment}"
#   filename         = "../../dist/rest-${var.project.environment}.zip"
#   handler          = "bootstrap"
#   role             = aws_iam_role.lambda_role.arn
#   depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
#   runtime          = "provided.al2"
#   architectures    = ["arm64"]
#   source_code_hash = filebase64sha256("../../dist/rest-${var.project.environment}.zip")

#   environment {
#     variables = {
#       ENVIRONMENT : var.project.environment
#       MONGO_DATABASE : local.envs.MONGO_DATABASE
#       MONGO_DATABASE_URL : local.envs.MONGO_DATABASE_URL
#     }
#   }

#   tags = {
#     project     = var.project.name
#     environment = var.project.environment
#   }
# }

resource "aws_lambda_function" "practice" {
  for_each         = data.external.lambdas
  function_name    = "${var.project.name}-practice-${var.project.environment}"
  filename         = "../../dist/practice-${var.project.environment}.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  architectures    = ["arm64"]
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../../dist/practice-${var.project.environment}.zip")

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


data "external" "lambdas" {
  program = ["sh", "../../scripts/lookup_lambdas.sh"]
}

output "lambdas" {
  value = data.external.lambdas.result
}
