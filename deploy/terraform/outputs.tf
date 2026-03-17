output "cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data for the cluster"
  value       = module.eks.cluster_certificate_authority_data
  sensitive   = true
}

output "configure_kubectl" {
  description = "Command to configure kubectl"
  value       = "aws eks update-kubeconfig --region ${var.region} --name ${module.eks.cluster_name}"
}

output "get_ingress_url" {
  description = "Command to get the Ingress external URL"
  value       = "kubectl -n ingress-nginx get svc ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'"
}

# Monitoring cluster outputs
output "monitoring_cluster_endpoint" {
  description = "Monitoring EKS cluster endpoint"
  value       = module.monitoring_eks.cluster_endpoint
}

output "monitoring_cluster_name" {
  description = "Monitoring EKS cluster name"
  value       = module.monitoring_eks.cluster_name
}

output "configure_monitoring_kubectl" {
  description = "Command to configure kubectl for monitoring cluster"
  value       = "aws eks update-kubeconfig --region ${var.region} --name ${module.monitoring_eks.cluster_name}"
}

output "get_grafana_url" {
  description = "Command to get Grafana external URL"
  value       = "kubectl --context arn:aws:eks:${var.region}:$(aws sts get-caller-identity --query Account --output text):cluster/${var.monitoring_cluster_name} -n monitoring get svc grafana -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'"
}

output "get_otel_collector_endpoint" {
  description = "Command to get OTEL collector internal NLB endpoint"
  value       = "kubectl --context arn:aws:eks:${var.region}:$(aws sts get-caller-identity --query Account --output text):cluster/${var.monitoring_cluster_name} -n monitoring get svc otel-collector-opentelemetry-collector -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'"
}
