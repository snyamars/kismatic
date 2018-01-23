# Create the Kubernetes ingress nodes (e.g. ingress1)
resource "digitalocean_droplet" "ingress_nodes" {
  count              = "${var.ingress_count}"
  image              = "${var.image}"
  name               = "${format("ingress%1d", count.index + 1)}"
  region             = "${var.region}"
  size               = "${var.droplet_size}"
  ssh_keys           = ["${digitalocean_ssh_key.default.id}"]
  tags               = ["${digitalocean_tag.cluster_tag.id}"]
  private_networking = true
  volume_ids         = ["${element(digitalocean_volume.ingress.*.id, count.index)}"]
}

# Create dedicated block storage for each node created above
resource "digitalocean_volume" "ingress" {
  count       = "${var.ingress_count}"
  region      = "${var.region}"
  name        = "${format("ingress%1d", count.index + 1)}"
  size        = 50
}

# Create the Load Balancer for the ingress nodes
resource "digitalocean_loadbalancer" "ingress_lb" {
  name = "ingress-loadbalancer"
  region = "${var.region}"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  forwarding_rule {
    entry_port = 443
    entry_protocol = "https"

    target_port = 443
    target_protocol = "https"

    tls_passthrough = true
  }

  healthcheck {
    port = 80
    protocol = "http"
    path = "/healthz"
  }

  droplet_ids = ["${digitalocean_droplet.ingress_nodes.*.id}"]
}