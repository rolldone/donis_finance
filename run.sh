#!/bin/bash
# donis_finance server runner
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR" || exit 1

PIDFILE="$DIR/pid.txt"

# Read the PID from the file
if [ -f "$PIDFILE" ]; then
    pid=$(cat "$PIDFILE")
    # Kill the process if it's running
    if kill -0 "$pid" 2> /dev/null; then
        echo "Current PID: $pid"
        kill "$pid"
        echo "Killed previous process"
        sleep 1
    fi
fi

# Build the binary
echo "Building..."
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server" 2>&1
echo "✅ Build OK"

# Run the server in background via docker exec
nohup docker exec donis-finance-app-1 sh -c "cd /app && ./server" > /dev/null 2>&1 &
echo $! > "$PIDFILE"
echo "Server started PID: $(cat $PIDFILE)"
echo "Access at http://localhost:8200"
