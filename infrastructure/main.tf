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
  size     = "s-2vcpu-2gb"
  ssh_keys = [var.ssh_key_fingerprint]

  connection {
    user        = "root"
    host        = "${digitalocean_droplet.ams3.ipv4_address}"
    private_key = "${file("~/.ssh/id_rsa")}"
    timeout     = "1m"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /root && mkdir install"
    ]
  }

  provisioner "file" {
    source      = "./scripts/install/"
    destination = "/root/install/"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/install/install_docker.sh",
      "chmod +x /root/install/install_minikube.sh",
      "/root/install/install_docker.sh",
      "/root/install/install_minikube.sh"
    ]
  }

  provisioner "local-exec" {
    command = "./scripts/prepare_remote_kube_access.sh ${digitalocean_droplet.ams3.ipv4_address}"
  }
}

output "ips" {
  value = "\n ams3: ${digitalocean_droplet.ams3.ipv4_address} \n"
}
