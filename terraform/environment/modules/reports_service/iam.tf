resource "aws_iam_role" "reports_api" {
  name               = "opg-reports-${var.environment_name}-api"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_policy.json
}

resource "aws_iam_role" "reports_frontend" {
  name               = "opg-reports-${var.environment_name}-frontend"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_policy.json
}

resource "aws_iam_role" "execution_role" {
  name               = "opg-reports-${var.environment_name}-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_policy.json
}

resource "aws_iam_role_policy" "execution_role" {
  role   = aws_iam_role.execution_role.id
  policy = data.aws_iam_policy_document.execution_role.json
}

data "aws_iam_policy_document" "execution_role" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "ssm:GetParameters",
      "secretsmanager:GetSecretValue",
    ]
  }
}

data "aws_iam_policy_document" "ecs_tasks_assume_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      identifiers = ["ecs-tasks.amazonaws.com"]
      type        = "Service"
    }
  }
}