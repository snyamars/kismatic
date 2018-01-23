cluster:
  networking:
    pod_cidr_block: 172.16.0.0/16
    service_cidr_block: 172.20.0.0/16
    update_hosts_files: true

  certificates:
    expiry: 17520h
    ca_expiry: 17520h

  ssh:
    user: root
    ssh_port: 22

# Add-ons are additional components that KET installs on the cluster.
add_ons:
  cni:
    disable: false
    provider: calico
    options:
      calico:
        mode: overlay
        log_level: info

  dns:
    disable: false

  heapster:
    disable: false
    options:
      heapster:
        replicas: 2
        service_type: ClusterIP
        sink: influxdb:http://heapster-influxdb.kube-system.svc:8086
      influxdb:
        pvc_name: ""

  dashboard:
    disable: false

  package_manager:
    disable: false
    provider: helm

  rescheduler:
    disable: false
