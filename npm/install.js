#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');

const pkg = require('../package.json');
const repo = 'aneokin12/vouch';
// Use version from package.json
const version = pkg.version;

const platform = os.platform();
const arch = os.arch();

const goos = platform === 'win32' ? 'windows' : platform;
const goarch = arch === 'x64' ? 'amd64' : (arch === 'arm64' ? 'arm64' : '');

if (!goarch) {
    console.error(`Unsupported architecture: ${arch}`);
    process.exit(1);
}

// Ensure the binary name matches the GoReleaser output
// GoReleaser formats it as ProjectName_Os_Arch
// e.g. vouch_linux_amd64, vouch_windows_amd64.exe, vouch_darwin_arm64
const ext = goos === 'windows' ? '.exe' : '';
const filename = `vouch_${goos}_${goarch}${ext}`;
const downloadUrl = `https://github.com/${repo}/releases/download/v${version}/${filename}`;

const binDir = path.join(__dirname, '../bin');
if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true });

const binName = goos === 'windows' ? 'vouch.exe' : 'vouch';
const binPath = path.join(binDir, binName);

// Skip if already downloaded
if (fs.existsSync(binPath)) {
    console.log('Binary already exists, skipping download.');
    process.exit(0);
}

console.log(`Downloading Vouch v${version} for ${goos}-${goarch}...`);

function download(url, dest) {
    return new Promise((resolve, reject) => {
        https.get(url, (res) => {
            if (res.statusCode === 301 || res.statusCode === 302) {
                return download(res.headers.location, dest).then(resolve).catch(reject);
            }
            if (res.statusCode !== 200) {
                return reject(new Error(`Failed to download: HTTP ${res.statusCode}`));
            }
            const file = fs.createWriteStream(dest);
            res.pipe(file);
            file.on('finish', () => {
                file.close();
                resolve();
            });
            file.on('error', (err) => {
                fs.unlink(dest, () => reject(err));
            });
        }).on('error', reject);
    });
}

download(downloadUrl, binPath)
    .then(() => {
        if (goos !== 'windows') {
            fs.chmodSync(binPath, 0o755);
        }
        console.log('Successfully installed Vouch!');
    })
    .catch((err) => {
        console.error(`Error downloading binary: ${err.message}`);
        console.error(`Attempted URL: ${downloadUrl}`);
        console.error('This may be because the release is not published yet, or your platform/architecture is unsupported.');
        process.exit(1);
    });
