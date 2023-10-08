variable "yc_cloud" {
  type = string
  description = "Yandex Cloud ID"
}

variable "yc_folder" {
  type = string
  description = "Yandex Cloud folder"
}

variable "yc_key_file" {
  type = string
  description = "Yandex Cloud OAuth token"
}

variable "db_password" {
  description = "MySQL user pasword"
}

# variable "subnets" {
#   type = map(map())
#   default = {}
# }
variable "virtual_machines" {
  type = map(map(string))
  default = {
    "wp-app-1" = { subnet = "wp-subnet-a", image_id = "fd80viupr3qjr5g6g9du", zone = "ru-central1-a" }
    "wp-app-2" = { subnet = "wp-subnet-b", image_id = "fd80viupr3qjr5g6g9du", zone = "ru-central1-a" }
    "wp-app-3" = { subnet = "wp-subnet-c", image_id = "fd80viupr3qjr5g6g9du", zone = "ru-central1-a" }
  }
}