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
  image    = "docker-18-04"
  region   = "ams3"
  size     = "s-2vcpu-2gb"
  ssh_keys = [var.ssh_key_fingerprint]

  connection {
    user        = "root"
    host        = "${digitalocean_droplet.ams3.ipv4_address}"
    private_key = "${file("~/.ssh/id_rsa")}"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root && mkdir install"
    ]
  }

  # TODO copy whole dir with installs scripts
  provisioner "file" {
    source      = "./scripts/install_minikube.sh"
    destination = "/root/install/install_minikube.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/install/install_minikube.sh",
      "/root/install/install_minikube.sh"
    ]
  }
}

output "ips" {
  value = "\n ams3: ${digitalocean_droplet.ams3.ipv4_address} \n"
}
