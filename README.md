# Kommodity Attestation Extension
Talos Extension to perform attestation reporting for establishing machine trust as part of zero-trust architecture. 

## Attestation Types

The following table summarizes the different types of attestation performed by this extension, what each attestation checks, and where the information is fetched from:

| Attestation Type   | What It Checks (Summary)                          | Source of Information                      |
|--------------------|---------------------------------------------------|--------------------------------------------|
| AppArmor           | Whether AppArmor is enabled                       | `/sys/module/apparmor/parameters/enabled`  |
| SELinux            | SELinux enforcement mode                          | `/sys/fs/selinux/enforce`                  |
| Secure Boot        | If Secure Boot is enabled                         | `/sys/firmware/efi/efivars/SecureBoot-*`   |
| Kernel Lockdown    | Current kernel lockdown mode                      | `/sys/kernel/security/lockdown`            |
| SquashFS           | If root filesystem is read-only SquashFS          | `/proc/mounts`                             |
| Talos Extensions   | Installed Talos extensions and their hashes       | `/usr/local/etc/containers`                |
| Image Layers       | Metadata of image layers (name, version, author)  | `/etc/extensions.yaml`                     |
| Talos Version      | Running Talos OS version                          | `/etc/os-release`                          |

Each attestation type collects measurements, evidence, and metadata to generate a comprehensive attestation report for the machine.
