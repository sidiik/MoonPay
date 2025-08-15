# API Gateway service
k8s_yaml([
  'services/api_gateway/zarf/k8s/deployment.yaml',
])

docker_build('api-gateway', 'services/api_gateway', dockerfile='services/api_gateway/Dockerfile')

k8s_resource(workload='api-gateway', port_forwards=8080)

# Auth service service
k8s_yaml([
  'services/auth_service/zarf/k8s/deployment.yaml',
  'services/auth_service/zarf/k8s/configmap.yaml',
  'services/auth_service/zarf/k8s/secret.yaml',
])

docker_build('auth-service', 'services/auth_service', dockerfile='services/auth_service/Dockerfile')

k8s_resource(workload='auth-service')

