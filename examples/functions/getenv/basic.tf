terraform {
  required_providers {
    local = {
      source = "marcozac/env"
      # version = "..."
    }
  }
  # Provider functions require Terraform 1.8 and later.
  required_version = ">= 1.8.0"
}

output "example_output" {
  value = provider::env::getenv("GETENV_EXAMPLE", true)
}
