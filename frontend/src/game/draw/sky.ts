import { SceneContext } from './palettes';

export function drawSky(ctx: SceneContext, w: number, h: number): void {
  const skyHeight = h * 0.55;
  const gfx = ctx.add.graphics();
  const [tr, tg, tb] = ctx.p.skyTop;
  const [br, bg, bb] = ctx.p.skyBot;

  const steps = 32;
  for (let i = 0; i < steps; i++) {
    const t = i / steps;
    const r = Math.round(tr + t * (br - tr));
    const g = Math.round(tg + t * (bg - tg));
    const b = Math.round(tb + t * (bb - tb));
    const color = (r << 16) | (g << 8) | b;
    const y = t * skyHeight;
    const sliceH = skyHeight / steps + 1;
    gfx.fillStyle(color, 1);
    gfx.fillRect(0, y, w, sliceH);
  }
}

export function drawStars(ctx: SceneContext, w: number, h: number): void {
  const skyHeight = h * 0.5;
  const gfx = ctx.add.graphics();

  for (let i = 0; i < 60; i++) {
    const x = Math.random() * w;
    const y = Math.random() * skyHeight;
    const size = 0.5 + Math.random() * 1.5;
    const brightness = 0.3 + Math.random() * 0.7;
    gfx.fillStyle(0xffffff, brightness);
    gfx.fillRect(
      Math.round(x),
      Math.round(y),
      Math.ceil(size),
      Math.ceil(size),
    );
  }

  for (let i = 0; i < 8; i++) {
    const star = ctx.add.graphics();
    const x = Math.random() * w;
    const y = Math.random() * skyHeight * 0.8;
    star.fillStyle(0xffffcc, 0.9);
    star.fillRect(Math.round(x), Math.round(y), 2, 2);

    ctx.addTween({
      targets: star,
      alpha: { from: 0.9, to: 0.2 },
      duration: 1500 + Math.random() * 2000,
      yoyo: true,
      repeat: -1,
      ease: 'Sine.easeInOut',
      delay: Math.random() * 2000,
    });
  }
}

export function drawMoon(ctx: SceneContext, w: number, h: number): void {
  const gfx = ctx.add.graphics();
  const mx = w * 0.8;
  const my = h * 0.08;

  gfx.fillStyle(0xeeeedd, 0.08);
  gfx.fillCircle(mx, my, 40);
  gfx.fillStyle(0xeeeedd, 0.15);
  gfx.fillCircle(mx, my, 28);
  gfx.fillStyle(0xeeeedd, 0.9);
  gfx.fillCircle(mx, my, 16);
  gfx.fillStyle(0x1a2450, 1);
  gfx.fillCircle(mx + 6, my - 3, 13);
  gfx.fillStyle(0xddddcc, 0.3);
  gfx.fillCircle(mx - 5, my + 2, 3);
  gfx.fillCircle(mx - 2, my - 4, 2);
}

export function drawClouds(ctx: SceneContext, w: number, h: number): void {
  const cloudData = [
    { x: w * 0.15, y: h * 0.08, s: 1.0 },
    { x: w * 0.55, y: h * 0.05, s: 1.3 },
    { x: w * 0.85, y: h * 0.12, s: 0.8 },
  ];

  cloudData.forEach(({ x, y, s }) => {
    const gfx = ctx.add.graphics();
    gfx.fillStyle(ctx.p.cloud!, 0.85);
    gfx.fillCircle(x, y, 18 * s);
    gfx.fillCircle(x - 16 * s, y + 4 * s, 14 * s);
    gfx.fillCircle(x + 16 * s, y + 4 * s, 14 * s);
    gfx.fillCircle(x - 8 * s, y - 6 * s, 12 * s);
    gfx.fillCircle(x + 10 * s, y - 4 * s, 13 * s);

    ctx.addTween({
      targets: gfx,
      x: '+=' + (20 + Math.random() * 15),
      duration: 12000 + Math.random() * 8000,
      yoyo: true,
      repeat: -1,
      ease: 'Sine.easeInOut',
    });
  });
}

export function drawSun(ctx: SceneContext, w: number, h: number): void {
  const gfx = ctx.add.graphics();
  const sx = w * 0.82;
  const sy = h * 0.06;

  gfx.fillStyle(0xfff176, 0.3);
  gfx.fillCircle(sx, sy, 32);
  gfx.fillStyle(0xffee58, 0.5);
  gfx.fillCircle(sx, sy, 22);
  gfx.fillStyle(0xffeb3b, 0.9);
  gfx.fillCircle(sx, sy, 14);

  const rays = ctx.add.graphics();
  rays.fillStyle(0xfff9c4, 0.2);
  for (let i = 0; i < 8; i++) {
    const angle = (i / 8) * Math.PI * 2;
    const rx = sx + Math.cos(angle) * 36;
    const ry = sy + Math.sin(angle) * 36;
    rays.fillCircle(rx, ry, 5);
  }

  ctx.addTween({
    targets: rays,
    alpha: 0.3,
    duration: 3000,
    yoyo: true,
    repeat: -1,
    ease: 'Sine.easeInOut',
  });
}

export function drawBirds(ctx: SceneContext, w: number, h: number): void {
  const skyH = h * 0.4;
  const birdCount = 1 + Math.floor(Math.random() * 2);

  for (let i = 0; i < birdCount; i++) {
    const gfx = ctx.add.graphics();
    const by = h * 0.05 + Math.random() * skyH * 0.5;
    const s = 0.7 + Math.random() * 0.4;
    const startX = -40 - Math.random() * w;

    gfx.setPosition(startX, by);

    const drawBird = (wingPhase: boolean): void => {
      gfx.clear();
      gfx.lineStyle(2 * s, ctx.p.birdColor!, 0.9);
      const wingY = wingPhase ? -5 * s : 3 * s;
      gfx.beginPath();
      gfx.moveTo(-9 * s, 0);
      gfx.lineTo(-4 * s, wingY);
      gfx.lineTo(0, 0);
      gfx.strokePath();
      gfx.beginPath();
      gfx.moveTo(0, 0);
      gfx.lineTo(4 * s, wingY);
      gfx.lineTo(9 * s, 0);
      gfx.strokePath();
      gfx.fillStyle(ctx.p.birdColor!, 0.9);
      gfx.fillCircle(0, 0, 1.5 * s);
    };

    drawBird(true);

    let wingUp = true;
    ctx.addTimedEvent({
      delay: 280 + Math.random() * 120,
      loop: true,
      callback: () => {
        wingUp = !wingUp;
        drawBird(wingUp);
      },
    });

    const flyAcross = (): void => {
      const waitTime = 5000 + Math.random() * 8000;
      ctx.addDelayedCall(waitTime, () => {
        const newY = h * 0.05 + Math.random() * skyH * 0.5;
        gfx.setPosition(-40, newY);
        const duration = 10000 + Math.random() * 8000;
        const drift = (Math.random() - 0.5) * skyH * 0.2;
        ctx.addTween({
          targets: gfx,
          x: w + 60,
          y: newY + drift,
          duration,
          ease: 'Linear',
          onComplete: () => flyAcross(),
        });
      });
    };

    ctx.addDelayedCall(i * 4000 + Math.random() * 3000, flyAcross);
  }
}
