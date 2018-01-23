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
    disable: false

    # Selecting 'custom' will result in a CNI ready cluster, however it is up to
    # you to configure a plugin after the install.
    # Options: 'calico','weave','contiv','custom'.
    provider: calico
    options:
      calico:

        # Options: 'overlay','routed'.
        mode: overlay

        # Options: 'warning','info','debug'.
        log_level: info
