resource "aws_ecs_service" "reports_frontend" {
  name                    = "opg-reports-frontend"
  cluster                 = aws_ecs_cluster.reports.id
  task_definition         = aws_ecs_task_definition.reports_frontend.arn
  desired_count           = 1
  enable_ecs_managed_tags = true
  platform_version        = null
  propagate_tags          = "TASK_DEFINITION"
  wait_for_steady_state   = true
  tags = {
    "Name" : "opg-reports-frontend-ecs-service-${var.environment_name}",
    "Version" : var.reports_frontend_tag
  }

  deployment_controller {
    type = "ECS"
  }

  network_configuration {
    security_groups  = [aws_security_group.reports_frontend.id]
    subnets          = data.aws_subnets.private.ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.reports_frontend.arn
    container_name   = "opg-reports-frontend"
    container_port   = 8000
  }
}

resource "aws_ecs_task_definition" "reports_frontend" {
  family                   = "opg-reports-frontend-${var.environment_name}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  container_definitions    = "[${local.reports_frontend}]"
  task_role_arn            = aws_iam_role.reports_frontend.arn
  execution_role_arn       = aws_iam_role.execution_role.arn
  tags = {
    "Name" : "opg-reports-frontend-ecs-service-${var.environment_name}"
    "Version" : var.reports_frontend_tag
  }
}

locals {
  reports_frontend = jsonencode({
    "name"      = "opg-reports-frontend",
    "cpu"       = 0,
    "essential" = true,
    "image"     = "${data.aws_ecr_repository.reports_frontend.repository_url}:${var.reports_frontend_tag}",
    portMappings = [{
      containerPort = 8000,
      hostPort      = 8000,
      protocol      = "tcp"
    }],
    healthCheck = {
      command     = ["CMD-SHELL", "curl -f http://localhost:8000/overview/ || exit 1"],
      startPeriod = 30,
      interval    = 15,
      timeout     = 10,
      retries     = 3
    },
    logConfiguration = {
      logDriver = "awslogs",
      options = {
        "awslogs-group"         = var.cloudwatch_log_group.name,
        "awslogs-region"        = data.aws_region.current.name,
        "awslogs-stream-prefix" = "opg-reports-frontend"
      }
    },
    "mountPoints" = [],
    "volumesFrom" = [],
  })
}
