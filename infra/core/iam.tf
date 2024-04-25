resource "aws_iam_role" "lambda_role" {
  name               = "${var.project.name}-lambda-role-${var.project.environment}"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Action": "sts:AssumeRole",
        "Principal": {
            "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
    }]
}
EOF
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "aws_iam_policy" "iam_policy_for_lambda" {
  name        = "${var.project.name}-lambda-policy-${var.project.environment}"
  path        = "/"
  description = "AWS IAM Policy for managing aws lambda role"
  policy      = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Action": [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
        ],
        "Resource": "arn:aws:logs:*:*:*",
        "Effect": "Allow"
    },
    {
        "Effect": "Allow",
        "Action": [
            "execute-api:ManageConnections" 
        ],
        "Resource": [
          "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.websocket_api.id}/${aws_apigatewayv2_stage.websocket_default.name}/*"
        ]
    },
    {
        "Effect": "Allow",
        "Action": "lambda:InvokeFunction",
        "Resource": [
          "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${aws_lambda_function.notification.function_name}",
          "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${aws_lambda_function.collecting-push.function_name}",
          "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${aws_lambda_function.collecting-get.function_name}"
        ]
    }]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda.arn
}

