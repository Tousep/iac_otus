output "load_balancer_public_ip" {
  description = "Public IP address of load balancer"
  value = yandex_lb_network_load_balancer.wp_lb.listener.*.external_address_spec[0].*.address
}

output "database_host_fqdn" {
  description = "DB hostname"
  value = local.dbhosts
  # value = yandex_compute_instance.wp-apps.id
  # value = yandex_vpc_network.wp-network.subnet_ids
  # value = element(yandex_vpc_network.wp-network.labels, "wp-subnet-a")  #local.dbhosts
}