# Satisfy
A command line tool for installing a linux Satisfactory dedicated server.

## What does it do?
- Installs pre-requisites
  - Adds `multiverse` or `nonfree` repository (ubuntu/debian)
  - Adds 32-bit package support w/ dpkg
  - Installs `lib32gcc1`
  - Installs `steamcmd`
  - Symlinks `steamcmd` to `/usr/games/steamcmd`
- Creates a `satisfactory` user
- Creates a `satisfactory.service` systemd service

## Installation
Download the latest release for your operating system (currently only Ubuntu and Debian are supported). Run the following command as root:
```
satisfy install
```