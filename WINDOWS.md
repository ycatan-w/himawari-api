# 🌻 Himawari - Installation Guide for Windows

## Installation

1. **Download** the correct `.zip` file for your computer:
   - Choose `386` for 32-bit Windows
   - Choose `amd64` for 64-bit Windows
2. **Extract** the downloaded `.zip` file to a folder (e.g. on your Desktop).
3. **Start the installer**:
   - Double-click `install.bat`, **or**
   - Right-click `install.ps1` and select **Run with PowerShell**
4. During installation:
   - Windows will automatically prompt you for administrator rights (UAC).
     The script restarts itself with elevated permissions, so you don’t need to run it as admin manually.
   - After approval, the installer will:
     - Copy the program files to `C:\Program Files\Himawari`
     - Create and start a Windows service named **HimawariServer**
     - Initialize the database (ready to use)
5. Once installation is complete, open your browser and go to:
   👉 [http://localhost:9740/web](http://localhost:9740/web)
   You can create an account or log in.

💡 **Tip:** To confirm the service is running, open the **Services** panel in Windows and look for `HimawariServer`.

---

## Uninstallation

1. Go back to the folder where you extracted the installer, or go to:
   `C:\Program Files\Himawari`
2. Run the uninstaller:
   - Double-click `uninstall.bat`, **or**
   - Right-click `uninstall.ps1` and select **Run with PowerShell**
3. During uninstallation:
   - Windows will prompt you for administrator rights (UAC).
     The script will restart itself with elevated permissions, so you don’t need to run it as admin manually.
   - After approval, the uninstaller will:
     - Stop and remove the Windows service
     - Remove the program files and uninstall script
     - Delete the Himawari database

💡 **Tip:** After uninstall, check the **Services** panel to confirm `HimawariServer` is no longer listed.
