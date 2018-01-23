# Create the Kubernetes master nodes (e.g. master1)
resource "digitalocean_droplet" "master_nodes" {
  count              = "${var.master_count}"
  image              = "${var.image}"
  name               = "${format("master%1d", count.index + 1)}"
  region             = "${var.region}"
  size               = "${var.droplet_size}"
  ssh_keys           = ["${digitalocean_ssh_key.default.id}"]
  tags               = ["${digitalocean_tag.cluster_tag.id}"]
  private_networking = true
  volume_ids         = ["${element(digitalocean_volume.master.*.id, count.index)}"]
}

# Create dedicated block storage for each node created above
resource "digitalocean_volume" "master" {
  count       = "${var.master_count}"
  region      = "${var.region}"
  name        = "${format("master%1d", count.index + 1)}"
  size        = 50
}

# Create the Load Balancer for the master nodes
resource "digitalocean_loadbalancer" "master_lb" {
  name = "master-loadbalancer"
  region = "${var.region}"

  forwarding_rule {
    entry_port = 6443
    entry_protocol = "https"

    target_port = 6443
    target_protocol = "https"

    tls_passthrough = true
  }

  healthcheck {
    port = 6443
    protocol = "tcp"
  }

  droplet_ids = ["${digitalocean_droplet.master_nodes.*.id}"]
}