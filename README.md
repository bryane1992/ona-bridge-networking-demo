# Docker Bridge Networking on Ona — Demo

Demonstrates that Docker bridge networking is broken by default on Ona environments, and how to fix it.

## The problem

Two services (`frontend` and `backend`) run via Docker Compose with bridge networking. The frontend resolves the backend by container name (`backend:8081`) using Docker's built-in DNS.

On Ona environments, this fails because:

1. The VM firewall (supervisor) uses **nftables**
2. The Docker daemon is configured to use **iptables-legacy** (forced by the devcontainers docker-in-docker feature)
3. These are separate rule stores — Docker's bridge NAT rules are invisible to nft
4. Traffic from bridge-networked containers can't route through the Ona proxy

Hitting `http://localhost:8080` returns a timeout or connection error because the frontend can't reach the backend.

## Branches

### `main` — shows the problem

Bridge networking with no workaround. Frontend fails to reach backend by container name.

### `workaround` — applies the fix

Adds a script that switches Docker to `iptables-nft` and restarts the daemon. Bridge networking and Docker DNS resolution work correctly.

## How to test

1. Open this repo in an Ona environment
2. Wait for `docker compose up` to finish (runs on start)
3. Open port 8080 — on `main` you'll see a timeout; on `workaround` you'll see the backend response

## The workaround

```bash
sudo update-alternatives --set iptables /usr/sbin/iptables-nft
sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-nft

for table in nat filter mangle raw; do
  sudo iptables-legacy -t "$table" -F 2>/dev/null || true
  sudo iptables-legacy -t "$table" -X 2>/dev/null || true
done

sudo pkill -x dockerd
while pgrep -x dockerd >/dev/null; do sleep 1; done
sudo service docker start
```

## Related

- [PDE-754](https://linear.app/ona-team/issue/PDE-754) — Linear ticket for proper fix
- [devcontainers/features#1633](https://github.com/devcontainers/features/issues/1633) — upstream issue
