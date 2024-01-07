resource "oci_identity_compartment" "main" {
  name          = "webster"
  description   = "Webster resources"
  enable_delete = true
}

resource "oci_objectstorage_bucket" "output" {
  name           = "webster-client-output"
  namespace      = data.oci_objectstorage_namespace.main.namespace
  compartment_id = oci_identity_compartment.main.id
  auto_tiering   = "InfrequentAccess"
  versioning     = "Enabled"
}

data "oci_objectstorage_namespace" "main" {
  compartment_id = oci_identity_compartment.main.id
}

output "storage_endpoint" {
  value = "https://${data.oci_objectstorage_namespace.main.namespace}.compat.objectstorage.${var.region}.oraclecloud.com"
}
