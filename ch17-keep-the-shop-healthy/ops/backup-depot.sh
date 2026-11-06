#!/bin/sh
# Back up the production database, keep two weeks of copies here, and send
# them off this server.
#
# This runs on the SERVER, from cron, not in development. Install it as
# /usr/local/bin/backup-depot.sh, `chmod +x` it, and add to `crontab -e`:
#
#   15 3 * * * /usr/local/bin/backup-depot.sh
#
# rclone needs an `offsite` destination (`rclone config`). To be alerted when
# a night's backup does not happen, add a last line that calls your uptime
# monitor's heartbeat URL:
#
#   curl -fsS https://your-monitor.example/ping/abc123
set -eu
dir=/var/backups/depot
mkdir -p "$dir"
docker exec depot-db pg_dump -U postgres --format=custom depot_production > "$dir/depot-$(date +%F).dump"
find "$dir" -name 'depot-*.dump' -mtime +14 -delete
rclone copy "$dir" offsite:depot-backups
