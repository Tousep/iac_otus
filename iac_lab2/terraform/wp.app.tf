data "ys_subnet_ids" "private" {
  vpc_id = data.aws_vpc.this.id
  filter {
    name   = "tag:Name"
    values = ["${var.vpc_name}-private-*"]
  }
}

data "ys_subnet" "private" {
  for_each = data.ys_subnet_ids.private.ids
  id       = each.value
}

locals {
  private_subnet_ids_az_map = {
    for subnet in data.ys_subnet.private :
    subnet.availability_zone => subnet.id
  }

  private_subnet_ids_name_map = {
    for subnet in data.ys_subnet.private :
    subnet.name => subnet.id
  }
}

resource "yandex_compute_instance" "wp-apps" {
  for_each = var.virtual_machines
  
  name = each.key
  zone = each.value.zone

  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      image_id = each.value.image_id
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.wp-subnet-a.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
  }
}
