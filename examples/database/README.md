# Database Connection Management

This example demonstrates how `defer` simplifies database cleanup.

## Setup
```bash
npm install better-sqlite3
node setup.js  # Creates sample database
```

## Run
```bash
djs users.djs
```

## What it demonstrates

- Guaranteed connection cleanup even on errors
- Simpler than try/finally blocks
- Common DevOps pattern
