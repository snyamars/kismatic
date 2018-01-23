provider "digitalocean" {
}

# Create a new tag
resource "digitalocean_tag" "cluster_tag" {
  name = "${var.cluster_name}"
}

# Upload SSH key which all instances will use.
resource "digitalocean_ssh_key" "default" {
  name       = "KET SSH Key"
  public_key = "${file("${var.public_ssh_key_path}")}"
}