terraform {
  required_providers {
    local = {
      source = "marcozac/env"
      # version = "..."
    }
  }
}

data "env_file" "example_env" {
  path = "${path.module}/.env"
}

output "example_output" {
  value = data.env_file.example_env.result["MY_ENV_VAR"]
}
