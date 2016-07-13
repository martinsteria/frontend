# AzureRM Network Interface

/*
Module Description = This module creates 'n' network interfaces for the Azure cloud.
*/

variable "resourceGroupName" {
  description = "The name of the resource group containing the NIC"
}
variable "location" {
  description = "The geographical location of the NIC"
}
variable "count" {
  default = "1"
  description = "The number of NIC instances to create"
}
variable "name" {
  description = "The name of the NIC"
}
variable "subnetID" {
  description = "The ID of the subnet to put the NIC in"
}
variable "publicIPAddressID" {
  default = ""
  description = "The ID of a public IP-Address associated with the NIC. Leave blank for none"
}

resource "azurerm_network_interface" "nic" {
  resource_group_name = "${var.resourceGroupName}"
  location = "${var.location}"
  count = "${var.count}"
  name = "${var.name}-${count.index}"
  ip_configuration {
    name = "${var.name}-${count.index}-ipconfig"
    subnet_id = "${var.subnetID}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id = "${var.publicIPAddressID}"
  }
}


output "idSplat" {
  /*
Output Description = String of IDs separated by commas.
*/
  value = "${join(",", azurerm_network_interface.nic.*.id)}"
}
