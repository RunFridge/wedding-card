#!/usr/bin/env node

// Sets the hall of fame submitted flag in localStorage
// so the congratulations dialog won't appear again.
//
// Usage: Run this script, then paste the copied command
// into the browser console on the wedding site.

const encoded = btoa(encodeURIComponent(JSON.stringify(Date.now())));
const cmd = `localStorage.setItem('_wch', '${encoded}')`;

try {
  require('child_process').execSync(
    process.platform === 'darwin'
      ? `printf '%s' ${JSON.stringify(cmd)} | pbcopy`
      : process.platform === 'win32'
        ? `echo|set /p=${JSON.stringify(cmd)} | clip`
        : `printf '%s' ${JSON.stringify(cmd)} | xclip -selection clipboard 2>/dev/null || printf '%s' ${JSON.stringify(cmd)} | xsel --clipboard 2>/dev/null || printf '%s' ${JSON.stringify(cmd)} | wl-copy 2>/dev/null`,
    { stdio: 'pipe' },
  );
  console.log('Copied to clipboard!');
  console.log('Paste in browser console to mark hall of fame as submitted.');
} catch {
  console.log(cmd);
}
