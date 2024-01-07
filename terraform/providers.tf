terraform {
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 5.23"
    }
  }

  backend "pg" {}
}

provider "oci" {
  region = var.region
}
