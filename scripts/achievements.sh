#!/usr/bin/env node

const readline = require('readline');

const ACHIEVEMENTS = [
  { id: 'gallery', icon: '🏆', name: '사진 갤러리 해금' },
  { id: 'gold', icon: '🥇', name: '금메달' },
  { id: 'silver', icon: '🥈', name: '은메달' },
  { id: 'bronze', icon: '🥉', name: '동메달' },
  { id: 'guestbook', icon: '📖', name: '축하 메시지' },
  { id: 'photo', icon: '📸', name: '포토그래퍼' },
  { id: 'heart', icon: '💕', name: '첫 좋아요' },
  { id: 'gallery_all', icon: '🖼️', name: '갤러리 마스터' },
  { id: 'theme', icon: '🎨', name: '인테리어 디자이너' },
  { id: 'speed', icon: '⚡', name: '스피드 러너' },
  { id: 'wedding_day', icon: '💒', name: '결혼식 당일 (secret)' },
];

const selected = new Set();
let cursor = 0;

function render() {
  process.stdout.write('\x1B[2J\x1B[H');
  console.log('┌─────────────────────────────────────────┐');
  console.log('│  Achievement Selector                   │');
  console.log('│  ↑/↓ move  SPACE toggle  A all  ENTER ok│');
  console.log('└─────────────────────────────────────────┘');
  console.log('');

  ACHIEVEMENTS.forEach((a, i) => {
    const pointer = i === cursor ? '>' : ' ';
    const check = selected.has(a.id) ? '✓' : ' ';
    const line = `${pointer} [${check}] ${a.icon}  ${a.name}`;
    if (i === cursor) {
      process.stdout.write(`\x1B[7m${line}\x1B[0m\n`);
    } else {
      console.log(line);
    }
  });

  console.log('');
  console.log(`Selected: ${selected.size}/${ACHIEVEMENTS.length}`);
}

function output() {
  const ids = ACHIEVEMENTS.filter((a) => selected.has(a.id)).map((a) => a.id);
  const encoded = btoa(encodeURIComponent(JSON.stringify(ids)));
  const cmd = `localStorage.setItem('_wca', '${encoded}')`;

  try {
    require('child_process').execSync(
      process.platform === 'darwin'
        ? `printf '%s' ${JSON.stringify(cmd)} | pbcopy`
        : process.platform === 'win32'
          ? `echo|set /p=${JSON.stringify(cmd)} | clip`
          : `printf '%s' ${JSON.stringify(cmd)} | xclip -selection clipboard 2>/dev/null || printf '%s' ${JSON.stringify(cmd)} | xsel --clipboard 2>/dev/null || printf '%s' ${JSON.stringify(cmd)} | wl-copy 2>/dev/null`,
      { stdio: 'pipe' },
    );
    console.log('');
    console.log('Copied to clipboard!');
  } catch {
    console.log('');
    console.log(cmd);
  }

  console.log('');
  console.log('── Selected achievements ──');
  ids.forEach((id) => {
    const a = ACHIEVEMENTS.find((x) => x.id === id);
    console.log(`  ${a.icon}  ${a.name}`);
  });
  console.log('');
}

readline.emitKeypressEvents(process.stdin);
if (process.stdin.isTTY) process.stdin.setRawMode(true);

render();

process.stdin.on('keypress', (str, key) => {
  if (key.ctrl && key.name === 'c') {
    process.stdout.write('\x1B[?25h');
    process.exit(0);
  }

  if (key.name === 'up') {
    cursor = (cursor - 1 + ACHIEVEMENTS.length) % ACHIEVEMENTS.length;
  } else if (key.name === 'down') {
    cursor = (cursor + 1) % ACHIEVEMENTS.length;
  } else if (key.name === 'space') {
    const id = ACHIEVEMENTS[cursor].id;
    if (selected.has(id)) selected.delete(id);
    else selected.add(id);
  } else if (str === 'a' || str === 'A') {
    if (selected.size === ACHIEVEMENTS.length) {
      selected.clear();
    } else {
      ACHIEVEMENTS.forEach((a) => selected.add(a.id));
    }
  } else if (key.name === 'return') {
    process.stdout.write('\x1B[2J\x1B[H');
    process.stdout.write('\x1B[?25h');
    if (process.stdin.isTTY) process.stdin.setRawMode(false);
    output();
    process.exit(0);
  }

  render();
});

process.stdout.write('\x1B[?25l');
