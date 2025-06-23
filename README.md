## Setting Up ADB Proxy as a Systemd Service

### 1. Create the Systemd Service File

Create a new file at `/etc/systemd/system/adb-proxy.service` with the following content:

```ini
[Unit]
Description=ADB Proxy from WSL to Windows
After=network.target

[Service]
ExecStart=/usr/local/bin/adb-proxy
Restart=always

[Install]
WantedBy=default.target
```

---

### 2. Install the `adb-proxy` Binary

Copy your compiled `adb-proxy` binary to `/usr/local/bin` and make it executable:

```sh
sudo cp path/to/adb-proxy /usr/local/bin/adb-proxy
sudo chmod +x /usr/local/bin/adb-proxy
```

---

### 3. Enable and Start the Service

Reload systemd, then enable and start the service:

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now adb-proxy.service
```

---

### 4. Verify the Service Status

Check that the service is running:

```sh
systemctl status adb-proxy.service
```

---

### 5. Start the ADB Server on Windows

On your Windows machine, start the ADB server so it listens on all interfaces:

```sh
adb -a start-server
```

This ensures that the ADB server is accessible from WSL.