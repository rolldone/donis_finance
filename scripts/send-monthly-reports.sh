#!/bin/bash
# ============================================================
# Send monthly financial reports to all members
# Intended to be triggered by cron on the 1st of each month.
#
# Usage:
#   ./scripts/send-monthly-reports.sh          # previous month
#   ./scripts/send-monthly-reports.sh --month 6 --year 2026
#   ./scripts/send-monthly-reports.sh --dry-run
# ============================================================
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR" || exit 1

go run ./cmd/console donisfinance:send-bulk-reports "$@"
