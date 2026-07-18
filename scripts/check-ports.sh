#!/bin/bash
# check-ports.sh — Cek apakah port yang dibutuhkan tersedia

PORTS=(8200 8201)
ALL_OK=true

for PORT in "${PORTS[@]}"; do
    if ss -tlnp "sport = :$PORT" 2>/dev/null | grep -q ":$PORT"; then
        echo "[FAIL] Port $PORT — already in use"
        ALL_OK=false
    else
        echo "[OK]   Port $PORT — available"
    fi
done

echo ""
if $ALL_OK; then
    echo "✓ All ports available. Ready to start."
    exit 0
else
    echo "✗ Some ports are in use. Free them first."
    exit 1
fi
