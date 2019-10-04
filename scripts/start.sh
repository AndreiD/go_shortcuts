#!/usr/bin/env bash
pm2 stop "xxxxx"
git pull
pm2 start "xxxxx" -- --env = production
pm2 logs "xxxxx"
