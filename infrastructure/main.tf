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
  size     = "s-1vcpu-2gb"
  ssh_keys = [var.ssh_key_fingerprint]

  connection {
    user        = "root"
    host        = "${digitalocean_droplet.ams3.ipv4_address}"
    private_key = "${file("~/.ssh/id_rsa")}"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root",
      "mkdir install",
    ]
  }

  provisioner "file" {
    source      = "./scripts/install_minikube.sh"
    destination = "/root/install/install_minikube.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root/install",
      "chmod +x ./install_minikube.sh",
      "./install_minikube.sh",
    ]
  }
}

resource "digitalocean_droplet" "lon1" {
  name     = "federation-lon1"
  image    = "ubuntu-18-04-x64"
  region   = "ams3"
  size     = "s-1vcpu-2gb"
  ssh_keys = [var.ssh_key_fingerprint]

  connection {
    user        = "root"
    host        = "${digitalocean_droplet.lon1.ipv4_address}"
    private_key = "${file("~/.ssh/id_rsa")}"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root",
      "mkdir install",
    ]
  }

  provisioner "file" {
    source      = "./scripts/install_minikube.sh"
    destination = "/root/install/install_minikube.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root/install",
      "chmod +x ./install_minikube.sh",
      "./install_minikube.sh",
    ]
  }
}

output "ips" {
  value = "\n ams3: ${digitalocean_droplet.ams3.ipv4_address}\n lon1: ${digitalocean_droplet.lon1.ipv4_address}"
}
