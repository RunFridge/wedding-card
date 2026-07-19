import { SceneContext } from './palettes';

export function drawTrees(ctx: SceneContext, w: number, h: number): void {
  const groundY = h * 0.62;
  ctx._treeGfx = [];
  const treeData = [
    { x: w * 0.08, canopyY: h * 0.34, scale: 1.2 },
    { x: w * -0.04, canopyY: h * 0.38, scale: 1.0 },
    { x: w * 0.92, canopyY: h * 0.34, scale: 1.2 },
    { x: w * 1.04, canopyY: h * 0.38, scale: 1.0 },
  ];

  treeData.forEach(({ x, canopyY, scale }) => {
    const gfx = ctx.add.graphics();
    ctx._treeGfx.push(gfx);
    const trunkW = 16 * scale;
    const canopyR = 40 * scale;
    const trunkTop = canopyY + canopyR * 0.3;
    const trunkH = groundY - trunkTop;

    gfx.fillStyle(ctx.p.trunkBase, 1);
    gfx.fillRect(x - trunkW / 2, trunkTop, trunkW, trunkH);

    gfx.fillStyle(ctx.p.trunkLight, 1);
    gfx.fillRect(x - trunkW / 2 + trunkW * 0.65, trunkTop, trunkW * 0.35, trunkH);

    gfx.fillStyle(ctx.p.trunkBark, 0.3);
    for (let ty = trunkTop; ty < groundY; ty += 6) {
      gfx.fillRect(x - trunkW / 2, ty, trunkW, 1);
    }

    gfx.fillStyle(ctx.p.trunkRoot, 0.7);
    gfx.fillEllipse(x, groundY, trunkW * 2, 8 * scale);

    gfx.fillStyle(ctx.p.canopy, 1);
    gfx.fillCircle(x, canopyY, canopyR);
    gfx.fillCircle(x - canopyR * 0.55, canopyY + canopyR * 0.35, canopyR * 0.75);
    gfx.fillCircle(x + canopyR * 0.55, canopyY + canopyR * 0.35, canopyR * 0.75);
    gfx.fillCircle(x - canopyR * 0.3, canopyY - canopyR * 0.4, canopyR * 0.6);
    gfx.fillCircle(x + canopyR * 0.3, canopyY - canopyR * 0.4, canopyR * 0.6);

    gfx.fillStyle(ctx.p.canopyShadow, 1);
    gfx.fillCircle(x + canopyR * 0.3, canopyY + canopyR * 0.1, canopyR * 0.5);
    gfx.fillCircle(x - canopyR * 0.2, canopyY - canopyR * 0.3, canopyR * 0.35);

    gfx.fillStyle(ctx.p.canopyHL, 0.5);
    gfx.fillCircle(x - canopyR * 0.4, canopyY - canopyR * 0.15, canopyR * 0.3);

    if (!ctx._dark) {
      gfx.fillStyle(0xff4444, 1);
      gfx.fillCircle(x - canopyR * 0.3, canopyY + canopyR * 0.1, 3 * scale);
      gfx.fillCircle(x + canopyR * 0.4, canopyY - canopyR * 0.1, 3 * scale);
      gfx.fillStyle(0xffaa00, 1);
      gfx.fillCircle(x + canopyR * 0.1, canopyY + canopyR * 0.3, 2.5 * scale);
    }
  });
}

export function drawBushes(ctx: SceneContext, w: number, h: number): void {
  const gfx = ctx.add.graphics();
  const groundY = h * 0.62;

  const bushData = [
    { x: w * 0.18, y: groundY + 5 },
    { x: w * 0.32, y: groundY + 14 },
    { x: w * 0.68, y: groundY + 8 },
    { x: w * 0.82, y: groundY + 16 },
    { x: w * 0.1, y: groundY + 20 },
    { x: w * 0.9, y: groundY + 22 },
  ];

  const berryColors = ctx._dark
    ? [0x992222, 0x442288, 0x995500, 0x993366]
    : [0xff4444, 0x6644ff, 0xff8800, 0xff66aa];

  bushData.forEach(({ x, y }, i) => {
    gfx.fillStyle(ctx.p.bushBody, 1);
    gfx.fillCircle(x, y, 13);
    gfx.fillCircle(x - 9, y + 3, 10);
    gfx.fillCircle(x + 10, y + 2, 11);

    gfx.fillStyle(ctx.p.bushShadow, 1);
    gfx.fillCircle(x + 5, y - 2, 7);

    gfx.fillStyle(ctx.p.bushHL, 0.5);
    gfx.fillCircle(x - 5, y - 4, 5);

    const bColor = berryColors[i % berryColors.length];
    gfx.fillStyle(bColor, 1);
    gfx.fillCircle(x - 5, y - 6, 2.5);
    gfx.fillCircle(x + 8, y - 2, 2.5);
    gfx.fillCircle(x + 1, y + 4, 2);
    gfx.fillCircle(x - 8, y + 1, 2);

    if (!ctx._dark) {
      gfx.fillStyle(0xffffff, 0.4);
      gfx.fillCircle(x - 5.5, y - 7, 1);
      gfx.fillCircle(x + 7.5, y - 3, 1);
    }
  });
}

export function drawMushrooms(ctx: SceneContext, w: number, h: number): void {
  const groundY = h * 0.62;
  const gfx = ctx.add.graphics();
  const spots = [
    { x: w * 0.15, y: groundY + 4 },
    { x: w * 0.42, y: groundY + 18 },
    { x: w * 0.58, y: groundY + 10 },
    { x: w * 0.85, y: groundY + 6 },
    { x: w * 0.7, y: groundY + 25 },
  ];

  spots.forEach(({ x, y }) => {
    const s = 0.8 + Math.random() * 0.5;
    gfx.fillStyle(ctx.p.mushroomStem, 1);
    gfx.fillRect(x - 2 * s, y - 6 * s, 4 * s, 6 * s);
    gfx.fillStyle(ctx.p.mushroomCap, 1);
    gfx.fillEllipse(x, y - 7 * s, 10 * s, 7 * s);
    gfx.fillStyle(0xffffff, ctx._dark ? 0.4 : 0.9);
    gfx.fillCircle(x - 2 * s, y - 9 * s, 1.5 * s);
    gfx.fillCircle(x + 3 * s, y - 8 * s, 1 * s);

    if (ctx._dark) {
      gfx.fillStyle(0xffaa44, 0.08);
      gfx.fillCircle(x, y - 7 * s, 12 * s);
    }
  });
}

export function drawTallGrass(ctx: SceneContext, w: number, h: number): void {
  const groundY = h * 0.62;
  const gfx = ctx.add.graphics();

  for (let i = 0; i < 40; i++) {
    const x = Math.random() * w;
    const y = groundY + Math.random() * 6;
    const bladeCount = 3 + Math.floor(Math.random() * 3);
    const color =
      ctx.p.grassColors[Math.floor(Math.random() * ctx.p.grassColors.length)];

    gfx.fillStyle(color, 1);
    for (let b = 0; b < bladeCount; b++) {
      const bx = x + (b - bladeCount / 2) * 2;
      const bladeH = 8 + Math.random() * 10;
      const lean = (Math.random() - 0.5) * 4;
      gfx.beginPath();
      gfx.moveTo(bx - 1, y);
      gfx.lineTo(bx + lean, y - bladeH);
      gfx.lineTo(bx + 1, y);
      gfx.closePath();
      gfx.fillPath();
    }
  }
}

export function drawFlowers(ctx: SceneContext, w: number, h: number): void {
  const groundY = h * 0.62;
  const flowerColors = ctx._dark
    ? [0x993838, 0x998830, 0x995588, 0x888888, 0x406688, 0x994838]
    : [0xff6b6b, 0xffdd57, 0xff9ff3, 0xffffff, 0x74b9ff, 0xff8a65];
  ctx.flowers = [];

  for (let i = 0; i < 35; i++) {
    const x = Math.random() * w;
    const y = groundY + 10 + Math.random() * (h - groundY - 50);
    const color = flowerColors[Math.floor(Math.random() * flowerColors.length)];
    const size = 3 + Math.random() * 3;

    const gfx = ctx.add.graphics();

    gfx.fillStyle(ctx.p.flowerStem, 1);
    gfx.fillRect(x - 1, y - size * 2, 2, size * 2);

    gfx.fillStyle(ctx.p.flowerLeaf, 0.8);
    gfx.fillEllipse(x - 4, y - size, 5, 3);
    gfx.fillEllipse(x + 4, y - size, 5, 3);

    const petalCount = 4 + Math.floor(Math.random() * 3);
    const flowerY = y - size * 2;
    for (let p = 0; p < petalCount; p++) {
      const angle = (p / petalCount) * Math.PI * 2;
      const px = x + Math.cos(angle) * size * 0.6;
      const py = flowerY + Math.sin(angle) * size * 0.6;
      gfx.fillStyle(color, 1);
      gfx.fillCircle(px, py, size * 0.5);
    }

    gfx.fillStyle(ctx._dark ? 0x998850 : 0xffee88, 1);
    gfx.fillCircle(x, flowerY, size * 0.35);

    ctx.flowers.push(gfx);
  }

  ctx.addTween({
    targets: ctx.flowers,
    y: '-=2',
    duration: 1800,
    yoyo: true,
    repeat: -1,
    ease: 'Sine.easeInOut',
  });
}

export function drawFenceVines(ctx: SceneContext, w: number, h: number): void {
  const fenceY = h * 0.6;
  const postHeight = 36;
  const gfx = ctx.add.graphics();

  for (let x = 30; x < w; x += 48) {
    if (Math.random() > 0.5) continue;
    const vineTop = fenceY - postHeight * 0.8;
    const vineLen = 12 + Math.random() * 10;

    gfx.lineStyle(2, ctx.p.vineStem, 0.8);
    gfx.beginPath();
    gfx.moveTo(x, vineTop);
    gfx.lineTo(x + 4, vineTop + vineLen * 0.4);
    gfx.lineTo(x - 2, vineTop + vineLen * 0.7);
    gfx.lineTo(x + 3, vineTop + vineLen);
    gfx.strokePath();

    gfx.fillStyle(ctx.p.vineLeaf, 0.9);
    gfx.fillEllipse(x + 4, vineTop + vineLen * 0.4, 5, 3);
    gfx.fillEllipse(x - 2, vineTop + vineLen * 0.7, 4, 3);

    if (Math.random() > 0.4) {
      gfx.fillStyle(ctx._dark ? 0x885588 : 0xff9ff3, 0.9);
      gfx.fillCircle(x + 3, vineTop + vineLen, 2.5);
      gfx.fillStyle(ctx._dark ? 0x998850 : 0xffee88, 1);
      gfx.fillCircle(x + 3, vineTop + vineLen, 1);
    }
  }
}
