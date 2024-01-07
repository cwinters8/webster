resource "oci_identity_compartment" "main" {
  name          = "webster"
  description   = "Webster resources"
  enable_delete = true
  # compartment_id = var.parent_compartment_ocid
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
