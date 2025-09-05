# 🌻 Himawari - Installation Guide for Linux

> Tested on Ubuntu (steps may vary slightly on other distributions)

## Installation

1. **Download** the correct `.tar.gz` archive for your system:
   - `386` for 32-bit
   - `amd64` for 64-bit
   - `arm64` for ARM processors
2. **Extract** the archive:

   ```bash
   tar -xzf file.tar.gz
   ```

3. **Run the installer** from the extracted folder:

   ```bash
   ./install.sh
   ```

   - If you forget to use `sudo`, the script will stop and show you the exact command to run (e.g. `sudo ./install.sh`).

4. The installer will:

   - Copy the program to `/usr/local/bin/`
   - Create and start a systemd service named **himawari**
   - Initialize the database (ready to use)

5. Once installed, open your browser and visit:
   👉 [http://localhost:9740/web](http://localhost:9740/web)

💡 **Tip:** To confirm the service is running:

```bash
systemctl status himawari
```

---

## Uninstallation

1. Open Terminal and run:

   ```bash
   himawari-server-uninstall
   ```

   - If you forget `sudo`, the script will show you the correct command.

2. The uninstaller will:

   - Stop and remove the systemd service
   - Remove the program files
   - Delete the Himawari database

3. ⚠️ After uninstall, you may need to manually delete the `himawari-server-uninstall` script from `/usr/local/bin/`.

💡 **Tip:** To confirm the service is fully removed:

```bash
systemctl status himawari
```
