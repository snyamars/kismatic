variable "private_ssh_key_path" {}

variable "public_ssh_key_path" {}

variable "ssh_user" {
  default = "vagrant"
}

variable "kismatic_version" {}

variable "cluster_name" {}

variable "cluster_owner" {}

variable "instance_size" {}

variable master_count {}

variable etcd_count {}

variable worker_count {}

variable ingress_count {}

variable storage_count {}
