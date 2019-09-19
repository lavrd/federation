variable "digitalocean_token" {
  type = string
}

variable "ssh_key_fingerprint" {
  type = string
}

provider "digitalocean" {
  token = var.digitalocean_token
}

resource "digitalocean_droplet" "ams3" {
  name     = "federation-ams3"
  image    = "ubuntu-18-04-x64"
  region   = "ams3"
  size     = "512mb"
  ssh_keys = [var.ssh_key_fingerprint]
}

resource "digitalocean_droplet" "lon1" {
  name     = "federation-lon1"
  image    = "ubuntu-18-04-x64"
  region   = "ams3"
  size     = "512mb"
  ssh_keys = [var.ssh_key_fingerprint]
}

output "ips" {
  value = "\n ams3: ${digitalocean_droplet.ams3.ipv4_address}\n lon1: ${digitalocean_droplet.lon1.ipv4_address}"
}
