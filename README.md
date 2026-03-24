# Terraform Provider Centron

A [Terraform](https://www.terraform.io) provider for managing resources in the [Centron Cloud](https://www.centron.de/).

> [!WARNING]
> **This provider is currently under active development (MVP phase).** Features are limited and API mappings are still being iterated on. Feel free to contribute, open issues, or submit Pull Requests!

## Features

### What it can do
- **Authentication**: Automatically generates and manages OAuth2 Bearer tokens (`client_credentials`) against the Centron API.
- **Servers**: Provisions Virtual Machines using the `centron_server` resource block. Includes support for specifying parameters like CPU cores, Memory, Image ID, Pool ID, Disks, and Subnets.

### What it cannot do (Yet)
- **Data Sources**: Querying existing Subnets, Disks, Pools, or Servers.
- **Advanced Networking**: Assigning additional IPv4/IPv6 addresses, attaching multiple network adapters, or creating VPCs/Subnets natively.
- **Storage**: Adding dynamic Disks or DVD ISO ISOs after creation.
- **Snapshots**: Creating and scheduling automated snapshots.
- **Subscriptions / Pool Management**: Fully managing Centron pool limits via Terraform.

---

## Configuration

To use the provider locally while it is under development, configure your credentials:

```hcl
terraform {
  required_providers {
    centron = {
      source = "joreichhardt/centron"
    }
  }
}

provider "centron" {
  # You can also use CENTRON_CLIENT_ID and CENTRON_CLIENT_SECRET environment variables
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Contributing and Local Testing

1. Clone the repository and navigate into the root directory.
2. Run `go mod tidy` and `go build`.
3. Configure a `.terraformrc` file overriding the provider installation:
   ```hcl
   provider_installation {
     dev_overrides {
       "registry.terraform.io/joreichhardt/centron" = "/path/to/your/compiled/directory"
     }
     direct {}
   }
   ```
4. Run `TF_CLI_CONFIG_FILE=/path/to/.terraformrc terraform plan` inside your test configurations folder to test changes.

Contributions are warmly welcomed! 🚀

## License

**AGPL-3.0**

This project is licensed under AGPL-3.0.

Commercial licenses are available for companies that:
- want to use this provider in closed-source environments
- want to avoid AGPL obligations

Contact: johannes.reichhardt@gmail.com
