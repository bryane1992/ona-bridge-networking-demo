#!/usr/bin/env bash
# Workaround: switch Docker from iptables-legacy to iptables-nft.
#
# The devcontainers docker-in-docker feature forces iptables-legacy, but the
# Ona VM firewall uses nftables. These are separate rule stores, so Docker's
# bridge NAT rules are invisible to nft and bridge networking breaks.
#
# This script switches to iptables-nft and restarts dockerd so bridge
# networking and Docker DNS resolution work correctly.

set -e

echo "Switching iptables to nft backend..."
sudo update-alternatives --set iptables /usr/sbin/iptables-nft
sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-nft

echo "Flushing legacy iptables tables..."
for table in nat filter mangle raw; do
  sudo iptables-legacy -t "$table" -F 2>/dev/null || true
  sudo iptables-legacy -t "$table" -X 2>/dev/null || true
done

echo "Restarting Docker daemon..."
if pgrep -x dockerd >/dev/null; then
  sudo pkill -x dockerd
  while pgrep -x dockerd >/dev/null; do sleep 1; done
fi
sudo service docker start

echo "Done. Docker is now using iptables-nft. Bridge networking should work."
