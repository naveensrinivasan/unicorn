# unicorn

unicorn generates email random email alias address for self hosted https://mailinabox.email/ servers.
The motivation behind this is to have privacy for all the email addresses that we give out.

The server settings are in `settings.yaml`.
```yaml

---
domain: box.email.xyz
username: user@email.xyz
emaildomain : email.xyz
password: password
```

This will generate a random email alias based on the settings `87AA8A3D-4385-348F-CDB6-3FE9617A22E4@email.xyz` 
mapping a forward to the original email address which is also in the `settings.yaml`.