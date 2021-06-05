#!/usr/bin/env bash

set -exu

apt-get update
apt-get install -y golang git
apt-get clean

set +x
cat << EOF > /etc/systemd/system/dsm.service
[Unit]
Description=DSM
After=network.service

[Service]
ExecStart=/src/run-lietu.sh
User=vagrant
Group=vagrant
Restart=always

[Install]
WantedBy=default.target
EOF
set -x

systemctl daemon-reload
systemctl enable dsm
systemctl start dsm
