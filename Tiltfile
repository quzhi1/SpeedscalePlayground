# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_remote', 'helm_remote')
load('ext://helm_resource', 'helm_resource', 'helm_repo')
compile_opt = 'GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 '

# Get Speedscale API key
speedscale_api_key = os.getenv('SPEEDSCALE_API_KEY')
if not speedscale_api_key:
  print('------------------------------------------------------------')
  print('Please get your personal Speedscale API token, and set it to SPEEDSCALE_API_KEY environment variable')
  print('------------------------------------------------------------')

# Install speedscale operator
helm_repo('speedscale', 'https://speedscale.github.io/operator-helm/')
helm_resource(
  'speedscale-operator',
  'speedscale/speedscale-operator',
  namespace="speedscale",
  flags=[
    "--create-namespace",
    "--set",
    "apiKey=" + speedscale_api_key, 
    "--set",
    "clusterName=minikube-" + os.getenv('USER').replace('@nylas.com', '').replace('.', '-'),
    "--set",
    "namespaceSelector=default",
  ]
)

# Label speedscale operator services
k8s_resource(
  'speedscale-operator',
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
