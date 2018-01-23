# Create the Kubernetes worker nodes (e.g. worker1)
resource "digitalocean_droplet" "worker_nodes" {
  count              = "${var.worker_count}"
  image              = "${var.image}"
  name               = "${format("worker%1d", count.index + 1)}"
  region             = "${var.region}"
  size               = "${var.droplet_size}"
  ssh_keys           = ["${digitalocean_ssh_key.default.id}"]
  tags               = ["${digitalocean_tag.cluster_tag.id}"]
  private_networking = true
  volume_ids         = ["${element(digitalocean_volume.worker.*.id, count.index)}"]
}

# Create dedicated block storage for each node created above
resource "digitalocean_volume" "worker" {
  count       = "${var.worker_count}"
  region      = "${var.region}"
  name        = "${format("worker%1d", count.index + 1)}"
  size        = 50
}