docker_build('taskmanager:dev', '.', dockerfile='Dockerfile')
k8s_yaml(['deploy/deployment.yaml', 'deploy/service.yaml'])
k8s_resource('taskmanager', port_forwards=8080)
