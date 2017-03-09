# hadoop-ottom8r
Automations for Hadoop operations

### Configuration
Default file: `hadoop-ottom8r.toml`

```
[Connection]
  nifi_host = "https://pi-datamart-proxy.corp.chartercom.com"
  nifi_user = ""
  nifi_pass = ""
  nifi_cert = ""

[Backup]
  backup_path = "/tmp/nifi-backup"
  config_file = "hadoop-ottom8r.toml"
  debug_mode = false
  log_level = "info"
  log_file = "nifi-backup.log"
  mock = false
```
Add username and password into TOML file. Once plaintext password saved, run hadoop-ottom8r with the "encrypt" switch: `hadoop-ottom8r --encrypt` - The application will apply AES encryption to your plaintext password and update the TOML file to reflect. *This needs to be performed, as plaintext passwords **will not** work.*