# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_remote', 'helm_remote')
compile_opt = 'GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 '

# Get Speedscale API key
speedscale_api_key = os.getenv('SPEEDSCALE_API_KEY')
if not speedscale_api_key:
  print('------------------------------------------------------------')
  print('Please get your personal Speedscale API token, and set it to SPEEDSCALE_API_KEY environment variable')
  print('------------------------------------------------------------')

# Install speedscale operator
helm_remote(
  'speedscale-operator',
  repo_name='speedscale',
  repo_url='https://speedscale.github.io/operator-helm/',
  namespace="speedscale",
  version="v1.0.9",
  create_namespace=True,
  set=[
    "apiKey=" + speedscale_api_key, 
    "clusterName=minikube-" + os.getenv('USER').replace('@nylas.com', '').replace('.', '-'),
    "namespaceSelector=default",
  ]
)

# Label speedscale operator services
speedscale_services = [
  'speedscale-operator-pre-install',
  'speedscale-operator',
]

for speedscale_service in speedscale_services:
  k8s_resource(
    speedscale_service,
    labels="speedscale",
  )

# Spin up redis
helm_remote(
  'redis',
  repo_name="bitnami",
  repo_url='https://charts.bitnami.com/bitnami',
  values=['helm/values-redis-dev.yaml'],
)

# Label redis
k8s_resource(
  'redis-master',
  labels="redis",
)

# Compile example application
local_resource(
  'hello-world-compile',
  compile_opt + 'go build -o bin/hello-world main.go',
  deps=['main.go'],
  labels="example-application",
)

# Build example docker image
docker_build_with_restart(
  'hello-world-image',
  '.',
  entrypoint=['/opt/app/bin/hello-world'],
  dockerfile='Dockerfile',
  only=[
    './bin',
  ],
  live_update=[
    sync('./bin', '/opt/app/bin'),
  ],
)

# Install example helm chart
k8s_yaml(helm('helm'))

# Label and port forwarding example applciation
k8s_resource(
  "hello-world",
  resource_deps=['speedscale-operator', 'redis-master'],
  port_forwards='8090:8090',
  labels="example-application",
)
