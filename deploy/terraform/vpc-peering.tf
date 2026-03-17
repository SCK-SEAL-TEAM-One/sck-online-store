# VPC Peering between sck-workshop and sck-monitoring
resource "aws_vpc_peering_connection" "workshop_to_monitoring" {
  vpc_id      = module.vpc.vpc_id
  peer_vpc_id = module.monitoring_vpc.vpc_id
  auto_accept = true

  accepter {
    allow_remote_vpc_dns_resolution = true
  }

  requester {
    allow_remote_vpc_dns_resolution = true
  }

  tags = {
    Name        = "workshop-to-monitoring"
    Project     = "sck-workshop"
    Environment = "workshop"
  }
}

# Routes from workshop VPC → monitoring VPC
resource "aws_route" "workshop_private_to_monitoring" {
  count                     = length(module.vpc.private_route_table_ids)
  route_table_id            = module.vpc.private_route_table_ids[count.index]
  destination_cidr_block    = "10.1.0.0/16"
  vpc_peering_connection_id = aws_vpc_peering_connection.workshop_to_monitoring.id
}

resource "aws_route" "workshop_public_to_monitoring" {
  count                     = length(module.vpc.public_route_table_ids)
  route_table_id            = module.vpc.public_route_table_ids[count.index]
  destination_cidr_block    = "10.1.0.0/16"
  vpc_peering_connection_id = aws_vpc_peering_connection.workshop_to_monitoring.id
}

# Routes from monitoring VPC → workshop VPC
resource "aws_route" "monitoring_private_to_workshop" {
  count                     = length(module.monitoring_vpc.private_route_table_ids)
  route_table_id            = module.monitoring_vpc.private_route_table_ids[count.index]
  destination_cidr_block    = "10.0.0.0/16"
  vpc_peering_connection_id = aws_vpc_peering_connection.workshop_to_monitoring.id
}

resource "aws_route" "monitoring_public_to_workshop" {
  count                     = length(module.monitoring_vpc.public_route_table_ids)
  route_table_id            = module.monitoring_vpc.public_route_table_ids[count.index]
  destination_cidr_block    = "10.0.0.0/16"
  vpc_peering_connection_id = aws_vpc_peering_connection.workshop_to_monitoring.id
}

# Security group rules: allow OTLP traffic from workshop VPC to monitoring cluster
resource "aws_security_group_rule" "monitoring_node_allow_otlp_from_workshop" {
  type              = "ingress"
  from_port         = 4317
  to_port           = 4318
  protocol          = "tcp"
  cidr_blocks       = ["10.0.0.0/16"]
  security_group_id = module.monitoring_eks.node_security_group_id
  description       = "Allow OTLP gRPC/HTTP from workshop VPC (node SG)"
}

resource "aws_security_group_rule" "monitoring_cluster_allow_otlp_from_workshop" {
  type              = "ingress"
  from_port         = 4317
  to_port           = 4318
  protocol          = "tcp"
  cidr_blocks       = ["10.0.0.0/16"]
  security_group_id = module.monitoring_eks.cluster_primary_security_group_id
  description       = "Allow OTLP gRPC/HTTP from workshop VPC (cluster primary SG)"
}

# Allow all TCP from workshop VPC to monitoring nodes (for NLB health checks and NodePort traffic)
resource "aws_security_group_rule" "monitoring_node_allow_all_from_workshop" {
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = ["10.0.0.0/16"]
  security_group_id = module.monitoring_eks.node_security_group_id
  description       = "Allow all TCP from workshop VPC (node SG)"
}

resource "aws_security_group_rule" "monitoring_cluster_allow_all_from_workshop" {
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = ["10.0.0.0/16"]
  security_group_id = module.monitoring_eks.cluster_primary_security_group_id
  description       = "Allow all TCP from workshop VPC (cluster primary SG)"
}
