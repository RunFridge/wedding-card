import satori from 'satori';
import { html } from 'satori-html';
import { Resvg } from '@resvg/resvg-js';
import { createHash } from 'node:crypto';
import { readFileSync, writeFileSync } from 'node:fs';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

// Read config values directly to avoid importing Vue-dependent module
const configPath = resolve(__dirname, '../src/config/wedding.ts');
const configSource = readFileSync(configPath, 'utf-8');

function extractString(name: string): string {
  const match = configSource.match(
    new RegExp(`export let ${name}\\s*=\\s*'([^']*)'`),
  );
  if (!match) throw new Error(`Could not extract ${name} from wedding config`);
  return match[1];
}

function extractDatetime(): Date {
  const match = configSource.match(/new Date\('([^']+)'\)/);
  if (!match)
    throw new Error('Could not extract WEDDING_DATETIME from wedding config');
  return new Date(match[1]);
}

const groomName = extractString('GROOM_ENG_NAME');
const brideName = extractString('BRIDE_ENG_NAME');
const venueName = extractString('VENUE_NAME');
const dt = extractDatetime();
const dateStr = dt
  .toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    timeZone: 'Asia/Seoul',
  })
  .replace(/\. /g, '. ');
const dayStr = dt.toLocaleDateString('ko-KR', {
  weekday: 'long',
  timeZone: 'Asia/Seoul',
});
const timeStr = dt.toLocaleTimeString('ko-KR', {
  hour: 'numeric',
  minute: '2-digit',
  hour12: true,
  timeZone: 'Asia/Seoul',
});

const WIDTH = 1200;
const HEIGHT = 630;

const fontData = readFileSync(
  resolve(__dirname, '../src/assets/fonts/NeoDunggeunmo.woff'),
);

const markup = html(`
<div style="display: flex; flex-direction: column; width: ${WIDTH}px; height: ${HEIGHT}px; background: #A0722A; position: relative; padding: 16px;">

  <div style="display: flex; position: absolute; top: 8px; left: 8px; width: 12px; height: 12px; border-radius: 50%; background: #8B6914; border: 2px solid #6B4513;"></div>
  <div style="display: flex; position: absolute; top: 8px; right: 8px; width: 12px; height: 12px; border-radius: 50%; background: #8B6914; border: 2px solid #6B4513;"></div>
  <div style="display: flex; position: absolute; bottom: 8px; left: 8px; width: 12px; height: 12px; border-radius: 50%; background: #8B6914; border: 2px solid #6B4513;"></div>
  <div style="display: flex; position: absolute; bottom: 8px; right: 8px; width: 12px; height: 12px; border-radius: 50%; background: #8B6914; border: 2px solid #6B4513;"></div>

  <div style="display: flex; flex: 1; flex-direction: column; align-items: center; justify-content: center; background: #F5E6C8; border: 2px solid #C4943A; padding: 40px 32px;">

    <div style="display: flex; font-size: 28px; color: #8B7355; letter-spacing: 0.3em; margin-bottom: 20px;">Wedding Invitation</div>

    <div style="display: flex; align-items: center; width: 60%; margin-bottom: 20px;">
      <div style="display: flex; flex: 1; height: 2px; background: linear-gradient(to right, transparent, #C4943A, transparent);"></div>
      <div style="display: flex; width: 12px; height: 12px; background: #8B6914; margin-left: 12px; margin-right: 12px; transform: rotate(45deg);"></div>
      <div style="display: flex; flex: 1; height: 2px; background: linear-gradient(to right, transparent, #C4943A, transparent);"></div>
    </div>

    <div style="display: flex; align-items: center; margin-bottom: 20px;">
      <span style="font-size: 72px; color: #6B4226;">${groomName}</span>
      <span style="font-size: 36px; color: #8B7355; margin-left: 20px; margin-right: 20px;">${'&'}</span>
      <span style="font-size: 72px; color: #6B4226;">${brideName}</span>
    </div>

    <div style="display: flex; align-items: center; width: 60%; margin-bottom: 20px;">
      <div style="display: flex; flex: 1; height: 2px; background: linear-gradient(to right, transparent, #C4943A, transparent);"></div>
      <div style="display: flex; width: 12px; height: 12px; background: #8B6914; margin-left: 12px; margin-right: 12px; transform: rotate(45deg);"></div>
      <div style="display: flex; flex: 1; height: 2px; background: linear-gradient(to right, transparent, #C4943A, transparent);"></div>
    </div>

    <div style="display: flex; font-size: 36px; color: #5C3A0E; margin-top: 4px;">${dateStr}</div>
    <div style="display: flex; font-size: 28px; color: #8B7355; margin-top: 8px;">${dayStr} ${timeStr}</div>
    <div style="display: flex; font-size: 28px; color: #8B7355; margin-top: 6px;">${venueName}</div>

  </div>
</div>
`);

const svg = await satori(markup, {
  width: WIDTH,
  height: HEIGHT,
  fonts: [
    {
      name: 'NeoDunggeunmo',
      data: fontData,
      weight: 400,
      style: 'normal',
    },
  ],
});

const resvg = new Resvg(svg, {
  fitTo: { mode: 'width', value: WIDTH },
});
const png = resvg.render().asPng();

const hash = createHash('md5').update(png).digest('hex').slice(0, 8);

const outPath = resolve(__dirname, '../public/title_card.png');
writeFileSync(outPath, png);

// Update index.html with cache-busting query string
const indexPath = resolve(__dirname, '../index.html');
let indexHtml = readFileSync(indexPath, 'utf-8');
indexHtml = indexHtml.replace(/title_card\.png(\?v=[a-f0-9]*)?/g, `title_card.png?v=${hash}`);
writeFileSync(indexPath, indexHtml);

console.log(`OG image generated: ${outPath} (${png.length} bytes)`);
