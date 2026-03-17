module "monitoring_eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.monitoring_cluster_name
  cluster_version = var.cluster_version

  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  enable_cluster_creator_admin_permissions = true

  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      most_recent              = true
      service_account_role_arn = module.monitoring_ebs_csi_irsa.iam_role_arn
    }
  }

  vpc_id     = module.monitoring_vpc.vpc_id
  subnet_ids = module.monitoring_vpc.private_subnets

  eks_managed_node_groups = {
    monitoring = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = var.monitoring_node_instance_types
      capacity_type  = var.use_spot_instances ? "SPOT" : "ON_DEMAND"

      min_size     = var.monitoring_node_min_size
      max_size     = var.monitoring_node_max_size
      desired_size = var.monitoring_node_desired_size
    }
  }

  tags = {
    Project     = "sck-workshop"
    Environment = "monitoring"
  }
}

module "monitoring_ebs_csi_irsa" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "~> 5.0"

  role_name             = "${var.monitoring_cluster_name}-ebs-csi"
  attach_ebs_csi_policy = true

  oidc_providers = {
    main = {
      provider_arn               = module.monitoring_eks.oidc_provider_arn
      namespace_service_accounts = ["kube-system:ebs-csi-controller-sa"]
    }
  }
}
