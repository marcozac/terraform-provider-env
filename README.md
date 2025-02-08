# Env Terraform Provider

This [Terraform](https://www.terraform.io/) provider offers a simple function, `getenv`, designed to retrieve environment variable values within your Terraform configurations.

It enables seamless integration of environment-specific data into your infrastructure as code workflows without the need to prefix your variables with `TF_VAR_` or other workarounds.

## 📦 **Installation**

To use this provider, add the following block to your Terraform configuration:

```hcl
terraform {
  required_providers {
    env = {
      source  = "marcozac/env"
      version = "~> 0.1"
    }
  }
}
```
