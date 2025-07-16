resource "aws_ecs_service" "reports_frontend" {
  name                    = "opg-reports-frontend"
  cluster                 = aws_ecs_cluster.reports.id
  task_definition         = aws_ecs_task_definition.reports_frontend.arn
  desired_count           = 1
  enable_ecs_managed_tags = true
  launch_type             = "FARGATE"
  platform_version        = "1.4.0"
  propagate_tags          = "TASK_DEFINITION"
  wait_for_steady_state   = true
  depends_on              = [aws_alb_target_group.reports_frontend]
  tags = {
    "Name" : "opg-reports-frontend-ecs-service-${var.environment_name}",
    "Version" : var.reports_frontend_tag
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = false
  }

  network_configuration {
    security_groups  = [aws_security_group.reports_frontend.id]
    subnets          = data.aws_subnets.private.ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.reports_frontend.arn
    container_name   = "opg-reports-frontend"
    container_port   = 8080
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
      containerPort = 8080,
      hostPort      = 8080,
      protocol      = "tcp"
    }],
    # healthCheck = {
    #   command     = ["CMD-SHELL", "wget -O /dev/null -S http://localhost:8080/ 2>&1 || exit 1    "],
    #   startPeriod = 30,
    #   interval    = 15,
    #   timeout     = 10,
    #   retries     = 3
    # },
    environment = [
      {
        name  = "SERVERS_API_ADDR",
        value = "${aws_service_discovery_service.reports_api.name}.${aws_service_discovery_private_dns_namespace.reports.name}:8081"
      },
      {
        name  = "SERVERS_FRONT_ADDR",
        value = ":8080"
      },
      {
        name  = "SERVERS_FRONT_DIRECTORY",
        value = "./"
      }
    ],
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
