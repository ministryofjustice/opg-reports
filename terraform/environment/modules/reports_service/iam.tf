resource "aws_iam_role" "reports_api" {
  name               = "${var.tags.application}-api"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_policy.json
}

resource "aws_iam_role" "reports_frontend" {
  name               = "${var.tags.application}-frontend"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_policy.json
}

resource "aws_iam_role" "execution_role" {
  name               = "${var.tags.application}-execution-role"
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