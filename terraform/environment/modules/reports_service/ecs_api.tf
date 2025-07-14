resource "aws_ecs_service" "reports_api" {
  name                    = "opg-reports-api"
  cluster                 = aws_ecs_cluster.reports.id
  task_definition         = aws_ecs_task_definition.reports_api.arn
  desired_count           = 1
  enable_ecs_managed_tags = true
  launch_type             = "FARGATE"
  platform_version        = "1.4.0"
  propagate_tags          = "TASK_DEFINITION"
  wait_for_steady_state   = true
  tags = {
    "Name" : "opg-reports-api-ecs-service-${var.environment_name}",
    "Version" : var.reports_api_tag
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = false
  }

  network_configuration {
    security_groups  = [aws_security_group.reports_api.id]
    subnets          = data.aws_subnets.private.ids
    assign_public_ip = false
  }

  service_registries {
    registry_arn = aws_service_discovery_service.reports_api.arn
  }
}

resource "aws_ecs_task_definition" "reports_api" {
  family                   = "opg-reports-api-${var.environment_name}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048
  container_definitions    = "[${local.reports_api}]"
  task_role_arn            = aws_iam_role.reports_api.arn
  execution_role_arn       = aws_iam_role.execution_role.arn
  tags = {
    "Name" : "opg-reports-api-ecs-service-${var.environment_name}"
    "Version" : var.reports_api_tag
  }
}

locals {
  reports_api = jsonencode({
    "name"      = "opg-reports-api",
    "cpu"       = 0,
    "essential" = true,
    "image"     = "${data.aws_ecr_repository.reports_api.repository_url}:${var.reports_api_tag}",
    portMappings = [{
      containerPort = 8081,
      hostPort      = 8081,
      protocol      = "tcp"
    }],
    logConfiguration = {
      logDriver = "awslogs",
      options = {
        "awslogs-group"         = var.cloudwatch_log_group.name,
        "awslogs-region"        = data.aws_region.current.name,
        "awslogs-stream-prefix" = "opg-reports-api"
      }
    },
    environment = [
      {
        name  = "SERVERS_API_ADDR",
        value = ":8081"
      },
      {
        name  = "DATABASE_PATH",
        value = "./databases/api.db"
      }
    ],
    "mountPoints" = [],
    "volumesFrom" = [],
  })
}

resource "aws_service_discovery_service" "reports_api" {
  name = "opg-reports-api"

  dns_config {
    namespace_id   = aws_service_discovery_private_dns_namespace.reports.id
    routing_policy = "MULTIVALUE"
    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}
