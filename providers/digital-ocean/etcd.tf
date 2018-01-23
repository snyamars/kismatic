# Create the Kubernetes Etcd nodes (e.g. etcd1)
resource "digitalocean_droplet" "etcd_nodes" {
  count              = "${var.etcd_count}"
  image              = "${var.image}"
  name               = "${format("etcd%1d", count.index + 1)}"
  region             = "${var.region}"
  size               = "${var.droplet_size}"
  ssh_keys           = ["${digitalocean_ssh_key.default.id}"]
  tags               = ["${digitalocean_tag.cluster_tag.id}"]
  private_networking = true
  volume_ids         = ["${element(digitalocean_volume.etcd.*.id, count.index)}"]
}

# Create dedicated block storage for each node created above
resource "digitalocean_volume" "etcd" {
  count       = "${var.etcd_count}"
  region      = "${var.region}"
  name        = "${format("etcd%1d", count.index + 1)}"
  size        = 50
}