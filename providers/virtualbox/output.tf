output "etcd_pub_ips" {
  value = ["${virtualbox_vm.etcd.*.network_adapter.0.ipv4_address}"]
}

output "master_pub_ips" {
  value = ["${virtualbox_vm.master.*.network_adapter.0.ipv4_address}"]
}

output "worker_pub_ips" {
  value = ["${virtualbox_vm.worker.*.network_adapter.0.ipv4_address}"]
}

output "ingress_pub_ips" {
  value = ["${virtualbox_vm.ingress.*.network_adapter.0.ipv4_address}"]
}

output "storage_pub_ips" {
  value = ["${virtualbox_vm.storage.*.network_adapter.0.ipv4_address}"]
}

output "etcd_priv_ips" {
  value = ["${virtualbox_vm.etcd.*.network_adapter.0.ipv4_address}"]
}

output "master_priv_ips" {
  value = ["${virtualbox_vm.master.*.network_adapter.0.ipv4_address}"]
}

output "worker_priv_ips" {
  value = ["${virtualbox_vm.worker.*.network_adapter.0.ipv4_address}"]
}

output "ingress_priv_ips" {
  value = ["${virtualbox_vm.ingress.*.network_adapter.0.ipv4_address}"]
}

output "storage_priv_ips" {
  value = ["${virtualbox_vm.storage.*.network_adapter.0.ipv4_address}"]
}

output "etcd_hosts" {
  value = ["${virtualbox_vm.etcd.*.network_adapter.0.ipv4_address}"]
}

output "master_hosts" {
  value = ["${virtualbox_vm.master.*.network_adapter.0.ipv4_address}"]
}

output "worker_hosts" {
  value = ["${virtualbox_vm.worker.*.network_adapter.0.ipv4_address}"]
}

output "ingress_hosts" {
  value = ["${virtualbox_vm.ingress.*.network_adapter.0.ipv4_address}"]
}

output "storage_hosts" {
  value = ["${virtualbox_vm.storage.*.network_adapter.0.ipv4_address}"]
}

output "master_lb" {
  value = ["${virtualbox_vm.master.0.network_adapter.0.ipv4_address}"]
}

output "ingress_lb" {
  value = ["${virtualbox_vm.ingress.0.network_adapter.0.ipv4_address}"]
}
