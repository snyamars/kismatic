# Create the GlusterFS storage nodes (e.g. storage1)
resource "digitalocean_droplet" "storage_nodes" {
  count              = "${var.storage_count}"
  image              = "${var.image}"
  name               = "${format("storage%1d", count.index + 1)}"
  region             = "${var.region}"
  size               = "${var.droplet_size}"
  ssh_keys           = ["${digitalocean_ssh_key.default.id}"]
  tags               = ["${digitalocean_tag.cluster_tag.id}"]
  private_networking = true
  volume_ids         = ["${element(digitalocean_volume.storage.*.id, count.index)}"]
}

# Create dedicated block storage for each node created above
resource "digitalocean_volume" "storage" {
  count       = "${var.storage_count}"
  region      = "${var.region}"
  name        = "${format("storage%1d", count.index + 1)}"
  size        = 50
}