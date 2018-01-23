
# possible options
# New York - nyc1, nyc2, nyc3
# San Fransico - sfo1, sfo2
# Amsterdam - ams2, ams3
# Singapore - sgp1
# London - lon1
# Frankfurt - fra1
# Toronto - tor1
# Bangalore - blr1
variable region {
  description = "Region to launch in"
  default     = "lon1"
}

variable "public_ssh_key_path" {
  description = "The path to the SSH public key"
}

variable image {
  description = "Name of the image to use"
  default     = "centos-7-x64"
}

variable "cluster_name" {
  description = "The name of the cluster"
}

variable droplet_size {
  description = "Size of the droplets"
  default     = "4gb"
}

variable master_count {
  description = "Number of k8s master droplets"
}

variable etcd_count {
  description = "Number of etcd droplets"
}

variable worker_count {
  description = "Number of k8s worker droplets"
}

variable ingress_count {
  description = "Number of k8s ingress droplets"
}

variable storage_count {
  description = "Number of GlusterFS droplets"
}