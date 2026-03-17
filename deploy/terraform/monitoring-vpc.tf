module "monitoring_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.monitoring_cluster_name}-vpc"
  cidr = "10.1.0.0/16"

  azs             = slice(data.aws_availability_zones.available.names, 0, 2)
  private_subnets = ["10.1.1.0/24", "10.1.2.0/24"]
  public_subnets  = ["10.1.101.0/24", "10.1.102.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true

  enable_dns_hostnames = true
  enable_dns_support   = true

  public_subnet_tags = {
    "kubernetes.io/role/elb"                                  = 1
    "kubernetes.io/cluster/${var.monitoring_cluster_name}" = "owned"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb"                         = 1
    "kubernetes.io/cluster/${var.monitoring_cluster_name}" = "owned"
  }

  tags = {
    Project     = "sck-workshop"
    Environment = "monitoring"
  }
}
