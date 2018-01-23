cluster:
  networking:
    pod_cidr_block: 172.16.0.0/16
    service_cidr_block: 172.20.0.0/16

  certificates:
    expiry: 17520h
    ca_expiry: 17520h

  ssh:
    user: ubuntu
    ssh_port: 22

# Add-ons are additional components that KET installs on the cluster.
add_ons:
  cni:
    provider: weave
  heapster:
    disable: false
    options:
      heapster:
        replicas: 2

        # Specify kubernetes ServiceType. Defaults to 'ClusterIP'.
        # Options: 'ClusterIP','NodePort','LoadBalancer','ExternalName'.
        service_type: ClusterIP

        # Specify the sink to store heapster data. Defaults to an influxdb pod
        # running on the cluster.
        sink: influxdb:http://heapster-influxdb.kube-system.svc:8086

      influxdb:

        # Provide the name of the persistent volume claim that you will create
        # after installation. If not specified, the data will be stored in
        # ephemeral storage.
        pvc_name: ""

  dashboard:
    disable: false

  package_manager:
    disable: false

    # Options: 'helm'
    provider: helm

  # The rescheduler ensures that critical add-ons remain running on the cluster.
  rescheduler:
    disable: false
