output "master_pub_ips" {
  value = ["${digitalocean_droplet.master_nodes.*.ipv4_address}"]
}

output "etcd_pub_ips" {
  value = ["${digitalocean_droplet.etcd_nodes.*.ipv4_address}"]
}

output "worker_pub_ips" {
  value = ["${digitalocean_droplet.worker_nodes.*.ipv4_address}"]
}

output "ingress_pub_ips" {
  value = ["${digitalocean_droplet.ingress_nodes.*.ipv4_address}"]
}

output "storage_pub_ips" {
  value = ["${digitalocean_droplet.storage_nodes.*.ipv4_address}"]
}

output "master_priv_ips" {
  value = ["${digitalocean_droplet.master_nodes.*.ipv4_address_private}"]
}

output "etcd_priv_ips" {
  value = ["${digitalocean_droplet.etcd_nodes.*.ipv4_address_private}"]
}

output "worker_priv_ips" {
  value = ["${digitalocean_droplet.worker_nodes.*.ipv4_address_private}"]
}

output "ingress_priv_ips" {
  value = ["${digitalocean_droplet.ingress_nodes.*.ipv4_address_private}"]
}

output "storage_priv_ips" {
  value = ["${digitalocean_droplet.storage_nodes.*.ipv4_address_private}"]
}

output "master_hosts" {
  value = ["${digitalocean_droplet.master_nodes.*.name}"]
}

output "etcd_hosts" {
  value = ["${digitalocean_droplet.etcd_nodes.*.name}"]
}

output "worker_hosts" {
  value = ["${digitalocean_droplet.worker_nodes.*.name}"]
}

output "ingress_hosts" {
  value = ["${digitalocean_droplet.ingress_nodes.*.name}"]
}

output "storage_hosts" {
  value = ["${digitalocean_droplet.storage_nodes.*.name}"]
}

output "master_lb" {
  value = ["${digitalocean_loadbalancer.master_lb.ip}"]
}

output "ingress_lb" {
  value = ["${digitalocean_loadbalancer.ingress_lb.ip}"]
}
