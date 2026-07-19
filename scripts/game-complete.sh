#!/usr/bin/env node

// Sets the game completion flag in localStorage
// so the photo gallery is unlocked without playing the game.
//
// Usage: Run this script, then paste the copied command
// into the browser console on the wedding site.

const encoded = btoa(encodeURIComponent(JSON.stringify(Date.now())));
const cmd = `localStorage.setItem('_wcg', '${encoded}')`;

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
  console.log('Paste in browser console to mark game as completed.');
} catch {
  console.log(cmd);
}
