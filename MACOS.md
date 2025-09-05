# 🌻 Himawari - Installation Guide for macOS

## Installation

1. **Download** the version that matches your Mac:
   - `.pkg` for the easiest installation
   - `.tar.gz` if you prefer a manual install
   - Choose `amd64`, `arm64`, or the universal version
2. If you downloaded the `.pkg`:
   - Double-click it and follow the installation wizard
   - macOS will prompt for administrator credentials (via system dialog)
3. If you downloaded the `.tar.gz`:
   - Extract it (double-click in Finder or run `tar -xzf file.tar.gz` in Terminal)
   - Open Terminal, go into the extracted folder, and run:
     ```bash
     ./install.sh
     ```
   - If you forget to run with `sudo`, the script will stop and tell you the exact command to retry (e.g. `sudo ./install.sh`).
4. The installer will:
   - Copy the program to `/usr/local/bin/`
   - Create and start a **LaunchDaemon** named `com.himawari.server` (visible in `/Library/LaunchDaemons/`)
   - Initialize the database (ready to use)
5. Once installed, open your browser and visit:
   👉 [http://localhost:9740/web](http://localhost:9740/web)

💡 **Tip:** To confirm the daemon is running:

```bash
launchctl list | grep himawari
```

---

## Uninstallation

1. Open Terminal and run:

   ```bash
   himawari-server-uninstall
   ```

   - If you forget `sudo`, the script will tell you to retry with it.

2. The uninstaller will:

   - Stop and remove the LaunchDaemon
   - Remove the program files
   - Delete the Himawari database

3. ⚠️ After uninstall, you may need to manually delete the `himawari-server-uninstall` script from `/usr/local/bin/`.

💡 **Tip:** You can check that the daemon is gone with:

```bash
launchctl list | grep himawari
```
