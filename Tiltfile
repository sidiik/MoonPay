# API Gateway service
k8s_yaml([
  'services/api_gateway/zarf/k8s/deployment.yaml',
  'services/api_gateway/zarf/k8s/secret.yaml',
])

docker_build('api-gateway', 'services/api_gateway', dockerfile='services/api_gateway/Dockerfile')

k8s_resource(workload='api-gateway', port_forwards=8080)

# Auth service 
k8s_yaml([
  'services/auth_service/zarf/k8s/deployment.yaml',
  'services/auth_service/zarf/k8s/configmap.yaml',
  'services/auth_service/zarf/k8s/secret.yaml',
])

docker_build('auth-service', 'services/auth_service', dockerfile='services/auth_service/Dockerfile')

k8s_resource(workload='auth-service')

# Notification service
k8s_yaml([
  'services/notification_service/zarf/k8s/deployment.yaml',
  'services/notification_service/zarf/k8s/configmap.yaml',
  'services/notification_service/zarf/k8s/secret.yaml',
])

docker_build('notification-service', 'services/notification_service', dockerfile='services/notification_service/Dockerfile')

k8s_resource(workload='notification-service')

# Wallet service 
k8s_yaml([
  'services/wallet_service/zarf/k8s/deployment.yaml',
  'services/wallet_service/zarf/k8s/configmap.yaml',
  'services/wallet_service/zarf/k8s/secret.yaml',
])

docker_build('wallet-service', 'services/wallet_service', dockerfile='services/wallet_service/Dockerfile')

k8s_resource(workload='wallet-service')