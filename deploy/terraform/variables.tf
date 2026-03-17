variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-7"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "sck-workshop"
}

variable "cluster_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.33"
}

variable "node_instance_types" {
  description = "EC2 instance types for worker nodes"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "node_desired_size" {
  description = "Desired number of worker nodes"
  type        = number
  default     = 2
}

variable "node_min_size" {
  description = "Minimum number of worker nodes"
  type        = number
  default     = 1
}

variable "node_max_size" {
  description = "Maximum number of worker nodes"
  type        = number
  default     = 3
}

variable "use_spot_instances" {
  description = "Use spot instances to save ~50% on compute costs"
  type        = bool
  default     = false
}

# Monitoring cluster variables
variable "monitoring_cluster_name" {
  description = "EKS monitoring cluster name"
  type        = string
  default     = "sck-monitoring"
}

variable "monitoring_node_instance_types" {
  description = "EC2 instance types for monitoring worker nodes"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "monitoring_node_desired_size" {
  description = "Desired number of monitoring worker nodes"
  type        = number
  default     = 2
}

variable "monitoring_node_min_size" {
  description = "Minimum number of monitoring worker nodes"
  type        = number
  default     = 1
}

variable "monitoring_node_max_size" {
  description = "Maximum number of monitoring worker nodes"
  type        = number
  default     = 3
}
