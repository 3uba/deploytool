# deploytool

A small single-binary CLI for deploying Git projects to a Linux server. One
command, `deploytool deploy <project>`, pulls the latest code from your repo,
rotates the previous release into a numbered backup, wires in a shared `.env`,
and (optionally) rebuilds and restarts the project with Docker Compose or a
Dockerfile.

It is a lightweight, personal-scale alternative to a full CI/CD pipeline: no
agents, no YAML pipelines, just a config file and a binary on the server.

## How it works

Each project you register lives under `root_dir/<project>/` with this layout:

```
<root_dir>/<project>/
├── current/      the live checkout (freshly cloned on every deploy)
├── backup/       previous releases, rotated into numbered folders (1, 2, 3, ...)
└── shared/       files that persist across deploys, e.g. your real .env
```

Running a deploy does, in order:

1. Reads the project's entry from `.config`.
2. Verifies Git access (`git ls-remote`).
3. Moves the current release into `backup/<n+1>/`.
4. Clones a fresh copy of the repo into `current/`.
5. Copies `shared/.env` into `current/.env`.
6. If `current/` contains a `docker-compose.yaml`, prompts to run
   `docker compose up -d --build prod-<project>`; otherwise, if a `Dockerfile`
   is present, prompts to build an image and run it on the project's port.

## Requirements

- A Linux server (the build targets `linux/amd64`)
- [Go 1.19+](https://go.dev/) to build the binary
- `git`
- `docker` with the Compose plugin (only needed for the Docker deploy step)

## Install

Build the `linux/amd64` binary and put it on your `PATH`. From a clone on the
server (or build locally and copy the binary over):

```bash
git clone https://github.com/3uba/deploytool
cd deploytool
./build                 # produces ./app/deploytool (linux/amd64)

# install it somewhere on PATH and point DT_PATH at the tool directory
sudo mkdir -p /opt/deploytool
sudo cp app/deploytool /opt/deploytool/
echo 'export PATH=$PATH:/opt/deploytool'   >> ~/.bashrc
echo 'export DT_PATH=/opt/deploytool'      >> ~/.bashrc
source ~/.bashrc
```

`DT_PATH` tells deploytool where it is installed (used by `update` and
`uninstall`).

## Configuration

Projects are stored in a `.config` file next to the tool. Copy the example and
fill it in, or use the interactive `deploytool create` command.

```ini
# .config
root_dir=/home/ubuntu

# one block per project
myapp_user=git-username
myapp_token=ghp_your_access_token
myapp_git_url=github.com/you/myapp.git
myapp_port=8080
```

`root_dir` is the base directory under which every project's
`current/backup/shared` tree is created. The per-project fields are the Git
credentials used to clone, plus the port used for the Dockerfile deploy path.

> **Security note:** deploytool authenticates to Git by embedding the token in
> the clone URL (`https://user:token@host`). That token can be visible in the
> server's process list while a clone runs, and it is stored in plaintext in
> `.config`. Use a scoped, least-privilege deploy token, keep `.config`
> readable only by the deploy user (`chmod 600 .config`), and never commit it
> (it is in `.gitignore`).

## Usage

```bash
deploytool create             # interactively register a new project in .config
deploytool deploy <project>   # deploy that project (alias: d)
deploytool update             # update deploytool itself (git pull in DT_PATH)
deploytool uninstall          # remove deploytool
deploytool nginx              # nginx setup (not implemented yet)
```

The `<project>` passed to `deploy` is the name you gave it in `create`.

## Status

This is a personal deployment helper, built for my own servers. It works for
the flow above but is intentionally small. The `nginx` command is a stub and
does nothing yet. Contributions and issues are welcome.

## License

MIT, see [LICENSE](LICENSE).
