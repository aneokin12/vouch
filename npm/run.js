#!/usr/bin/env node

const { spawnSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const binName = os.platform() === 'win32' ? 'vouch.exe' : 'vouch';
const binPath = path.join(__dirname, '../bin', binName);

if (!fs.existsSync(binPath)) {
    console.error(`Vouch binary not found at ${binPath}.`);
    console.error(`Please make sure it was downloaded properly via 'npm install'.`);
    process.exit(1);
}

// Pass all arguments correctly
const args = process.argv.slice(2);
const result = spawnSync(binPath, args, { stdio: 'inherit' });

if (result.error) {
    console.error(result.error);
    process.exit(1);
}
process.exit(result.status || 0);
