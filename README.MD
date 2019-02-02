# unicorn

unicorn generates email random email alias address for self hosted https://mailinabox.email/ servers.
The motivation behind this is to have privacy for all the email addresses that we give out.

The server settings are in `settings.yaml`.
```yaml

---
domain: box.email.xyz
username: user@email.xyz
emaildomain : email.xyz
```

## password
The password is retrieved from the  OSX `keychain` and it is not stored in the config file as it posses
a big vulnerability. The assumption is that the `keychain` has password for the `username` in the `settings.yaml`.

This will generate a random email alias based on the settings `87AA8A3D4385348FCDB63FE9617A22E4@email.xyz` 
mapping a forward to the original email address which is also in the `settings.yaml`.

## Download
https://github.com/naveensrinivasan/unicorn/releases
