package config

// ExampleConfig contains the example configuration YAML
const ExampleConfig = `# StackUp Configuration File
# Complete example with all supported features

profile: full-stack-dev

settings:
  auto_update_path: true
  verify_installations: true

tools:
  # Git - Simple package manager install
  - name: git
    display_name: "Git"
    version: latest
    description: "Version control system"
    linux:
      package_names:
        apt: git
        dnf: git
        pacman: git
    macos:
      brew: git
    windows:
      package_names:
        winget: Git.Git
        choco: git
    verify_command: "git --version"

  # Docker - depends on WSL on Windows
  - name: docker
    display_name: "Docker Desktop"
    version: latest
    description: "Container platform"
    windows:
      installer: "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
      type: exe
    macos:
      brew: docker
    linux:
      custom_commands:
        - command: curl
          args: ["-fsSL", "https://get.docker.com", "-o", "get-docker.sh"]
          description: "Download Docker install script"
        - command: sh
          args: ["get-docker.sh"]
          sudo: true
    verify_command: "docker --version"

presets:
  web-dev:
    description: "Web development stack"
    tools: ["git", "docker"]
`