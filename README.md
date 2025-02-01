# Directory backup client-server

This project contains a client that creates a backup archive 
and a server that runs on a different machine and receives backups.

The purpose is to backup the important "projects" directory to a second machine
on the same LAN network.

## Security

- This app is designed to run on LAN, so no encryption at the moment, but can be added relatively easily using HTTPS (ListenAndServeTLS)
- The server requires a password, and this password is stored as cleartext, as this is more of a key than a password
- The server has a rate-limit, so attackers cannot fill up the storage or force deletion of old backups.

**Summary:** there are things to improve, **do not use outside protected LAN**