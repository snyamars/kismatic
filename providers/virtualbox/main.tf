provider "virtualbox" {}

resource "virtualbox_vm" "master" {
  # count = "${var.master_count}"
  count = 1
  name  = "${var.cluster_name}-master-${count.index}"

  image = "../vendor-terraform/ubuntu1604.box"

  cpus   = 1
  memory = "2GB"

  network_adapter {
    type = "nat"
  }
}

resource "virtualbox_vm" "etcd" {
  # count = "${var.etcd_count}"
  count = 1
  name  = "${var.cluster_name}-etcd-${count.index}"

  image = "../vendor-terraform/ubuntu1604.box"

  cpus   = 1
  memory = "2GB"

  network_adapter {
    type = "nat"
  }
}

resource "virtualbox_vm" "worker" {
  # count = "${var.worker_count}"
  count = 1
  name  = "${var.cluster_name}-worker-${count.index}"

  image = "../vendor-terraform/ubuntu1604.box"

  cpus   = 1
  memory = "2GB"

  network_adapter {
    type = "nat"
  }
}

resource "virtualbox_vm" "ingress" {
  # count = "${var.ingress_count}"
  count = 1
  name  = "${var.cluster_name}-ingress-${count.index}"

  image = "../vendor-terraform/ubuntu1604.box"

  cpus   = 1
  memory = "2GB"

  network_adapter {
    type = "nat"
  }
}

resource "virtualbox_vm" "storage" {
  # count = "${var.master_count}"
  count = 1
  name  = "${var.cluster_name}-bastion-${count.index}"

  image = "../vendor-terraform/ubuntu1604.box"

  cpus   = 1
  memory = "2GB"

  network_adapter {
    type = "nat"
  }
}
