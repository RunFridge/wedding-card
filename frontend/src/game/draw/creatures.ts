import { SceneContext } from './palettes';

export function drawFireflies(ctx: SceneContext, w: number, h: number): void {
  for (let i = 0; i < 10; i++) {
    const fx = w * 0.1 + Math.random() * w * 0.8;
    const fy = h * 0.35 + Math.random() * h * 0.35;

    const glow = ctx.add.graphics();
    glow.fillStyle(0xffee44, 0.12);
    glow.fillCircle(fx, fy, 6);
    glow.fillStyle(0xffee44, 0.7);
    glow.fillCircle(fx, fy, 2);

    ctx.addTween({
      targets: glow,
      x: '+=' + (-20 + Math.random() * 40),
      y: '+=' + (-15 + Math.random() * 30),
      alpha: { from: 0.8, to: 0.1 },
      duration: 3000 + Math.random() * 3000,
      yoyo: true,
      repeat: -1,
      ease: 'Sine.easeInOut',
      delay: i * 400,
    });
  }
}

export function drawButterflies(
  ctx: SceneContext,
  w: number,
  h: number,
): void {
  const colors = [0xff9ff3, 0x74b9ff, 0xffdd57, 0xffffff];

  for (let i = 0; i < 4; i++) {
    const bx = w * 0.2 + Math.random() * w * 0.6;
    const by = h * 0.3 + Math.random() * h * 0.3;
    const color = colors[i % colors.length];

    const gfx = ctx.add.graphics();
    gfx.fillStyle(color, 0.8);
    gfx.fillEllipse(bx - 4, by - 2, 6, 4);
    gfx.fillEllipse(bx + 4, by - 2, 6, 4);
    gfx.fillEllipse(bx - 3, by + 2, 4, 3);
    gfx.fillEllipse(bx + 3, by + 2, 4, 3);
    gfx.fillStyle(0x333333, 1);
    gfx.fillRect(bx - 0.5, by - 3, 1, 6);

    ctx.addTween({
      targets: gfx,
      x: '+=' + (-30 + Math.random() * 60),
      y: '+=' + (-20 + Math.random() * 40),
      alpha: { from: 0.9, to: 0.4 },
      scaleX: { from: 1, to: 0.7 },
      duration: 4000 + Math.random() * 3000,
      yoyo: true,
      repeat: -1,
      ease: 'Sine.easeInOut',
      delay: i * 800,
    });
  }
}

export function drawAnimals(ctx: SceneContext, w: number, h: number): void {
  ctx._animals.forEach((a) => {
    if (a.wanderTween) a.wanderTween.stop();
    if (!ctx._transitioning) a.sprite.destroy();
  });
  ctx._animals = [];

  const groundY = h * 0.65;
  const scale = Math.max(1.5, Math.min(3, w / 200));
  const nightTint = 0x223322;

  const spawnAnimal = (
    key: string,
    animIdle: string,
    animWalk: string,
    x: number,
    y: number,
  ): void => {
    const sprite = ctx.add.sprite(x, y, key, 0);
    sprite.setScale(scale);
    sprite.setDepth(0);
    sprite.setFrame(0);
    if (ctx._dark) sprite.setTint(nightTint);

    const entry: {
      sprite: Phaser.GameObjects.Sprite;
      wanderTween: Phaser.Tweens.Tween | null;
    } = { sprite, wanderTween: null };

    const wander = (): void => {
      const destX = w * 0.1 + Math.random() * w * 0.8;
      const dist = Math.abs(destX - sprite.x);
      const duration = 3000 + dist * 10 + Math.random() * 4000;

      sprite.setFlipX(destX < sprite.x);
      sprite.play(animWalk, true);

      const tween = ctx.addTween({
        targets: sprite,
        x: destX,
        duration,
        ease: 'Linear',
        onComplete: () => {
          sprite.stop();
          sprite.setFrame(0);
          const pause = 2000 + Math.random() * 5000;
          ctx.addDelayedCall(pause, wander);
        },
      });
      entry.wanderTween = tween;
    };

    ctx._animals.push(entry);

    const startDelay = Math.random() * 3000;
    ctx.addDelayedCall(startDelay, wander);
  };

  spawnAnimal(
    'chicken',
    'chicken_idle',
    'chicken_walk',
    w * 0.3 + Math.random() * w * 0.4,
    groundY + 10 + Math.random() * 30,
  );

  spawnAnimal(
    'chick',
    'chick_idle',
    'chick_walk',
    w * 0.3 + Math.random() * w * 0.15,
    groundY + 20 + Math.random() * 25,
  );
  spawnAnimal(
    'chick',
    'chick_idle',
    'chick_walk',
    w * 0.5 + Math.random() * w * 0.15,
    groundY + 25 + Math.random() * 25,
  );
}

export function drawCharacters(ctx: SceneContext, w: number, h: number): void {
  const cx = w * 0.5;
  const cy = h * 0.72;
  const charScale = Math.min(w / 400, 1) * 0.35;

  drawFloralArch(ctx, w, h, cx, cy, charScale);

  const outlineColor = ctx._dark ? 0x0a1a06 : 0x2d5a1b;
  const addOutline = (img: Phaser.GameObjects.Image): Phaser.GameObjects.Image => {
    if (img.preFX) {
      img.preFX.addGlow(outlineColor, 2, 0, false, 0.8, 12);
    }
    return img;
  };

  const charOpen = addOutline(
    ctx.add.image(cx, cy, 'char_open').setScale(charScale),
  );
  const charBlink = addOutline(
    ctx.add.image(cx, cy, 'char_blink').setScale(charScale).setVisible(false),
  );
  const charSmile = addOutline(
    ctx.add.image(cx, cy, 'char_smile').setScale(charScale).setVisible(false),
  );

  let smiling = false;

  const showDefault = (): void => {
    charOpen.setVisible(true);
    charBlink.setVisible(false);
    charSmile.setVisible(false);
  };

  const scheduleBlink = (): void => {
    ctx.addDelayedCall(2000 + Math.random() * 3000, () => {
      if (smiling) {
        scheduleBlink();
        return;
      }
      charOpen.setVisible(false);
      charBlink.setVisible(true);

      ctx.addDelayedCall(150, () => {
        if (!smiling) showDefault();
        scheduleBlink();
      });
    });
  };

  scheduleBlink();

  ctx._onDocumentClick = () => {
    smiling = true;
    charOpen.setVisible(false);
    charBlink.setVisible(false);
    charSmile.setVisible(true);

    if (ctx.smileTimer) ctx.smileTimer.remove();
    ctx.smileTimer = ctx.addDelayedCall(1500, () => {
      smiling = false;
      showDefault();
    });
  };
  document.addEventListener('pointerdown', ctx._onDocumentClick);
}

function drawFloralArch(
  ctx: SceneContext,
  w: number,
  h: number,
  cx: number,
  cy: number,
  charScale: number,
): void {
  const gfx = ctx.add.graphics();
  const imgH = 250 * charScale;
  const archW = imgH * 1.6;
  const archTopY = cy - imgH * 1.3;
  const archBaseY = cy + imgH * 0.5;

  gfx.lineStyle(6, ctx.p.trunkBase, 0.9);
  const archPoints: { x: number; y: number }[] = [];
  const steps = 32;

  const legSteps = 8;
  for (let i = 0; i <= legSteps; i++) {
    const t = i / legSteps;
    archPoints.push({
      x: cx - archW,
      y: archBaseY - t * (archBaseY - archTopY),
    });
  }

  const curveSteps = steps - legSteps * 2;
  for (let i = 0; i <= curveSteps; i++) {
    const t = i / curveSteps;
    const angle = t * Math.PI;
    const ax = cx - archW * Math.cos(angle);
    const ay = archTopY - Math.sin(angle) * archW * 0.5;
    archPoints.push({ x: ax, y: ay });
  }

  for (let i = 0; i <= legSteps; i++) {
    const t = i / legSteps;
    archPoints.push({
      x: cx + archW,
      y: archTopY + t * (archBaseY - archTopY),
    });
  }

  gfx.beginPath();
  gfx.moveTo(archPoints[0].x, archPoints[0].y);
  for (let i = 1; i < archPoints.length; i++) {
    gfx.lineTo(archPoints[i].x, archPoints[i].y);
  }
  gfx.strokePath();

  const leafColors = [ctx.p.canopy, ctx.p.canopyHL, ctx.p.hills[2]];
  for (let i = 0; i < archPoints.length; i += 2) {
    const pt = archPoints[i];
    const color = leafColors[i % leafColors.length];
    const side = i % 4 < 2 ? -1 : 1;

    gfx.fillStyle(color, 0.9);
    gfx.fillEllipse(pt.x + side * 6, pt.y - 3, 8, 5);
    gfx.fillEllipse(pt.x - side * 4, pt.y + 4, 7, 4);
  }

  if (ctx._dark) {
    const lightColors = [0xff4444, 0x44ff44, 0x4488ff, 0xffdd00, 0xff66cc, 0x44ffff];
    for (let i = 1; i < archPoints.length - 1; i += 2) {
      const pt = archPoints[i];
      const color = lightColors[i % lightColors.length];
      if (i > 1) {
        const prev = archPoints[i - 2];
        gfx.lineStyle(1, 0x333333, 0.6);
        gfx.beginPath();
        gfx.moveTo(prev.x, prev.y);
        gfx.lineTo(pt.x, pt.y);
        gfx.strokePath();
      }
      const bulb = ctx.add.graphics();
      bulb.fillStyle(color, 0.3);
      bulb.fillCircle(pt.x, pt.y, 6);
      bulb.fillStyle(color, 0.9);
      bulb.fillCircle(pt.x, pt.y, 3);
      bulb.fillStyle(0xffffff, 0.6);
      bulb.fillCircle(pt.x - 1, pt.y - 1, 1);
      ctx.addTween({
        targets: bulb,
        alpha: { from: 0.4 + Math.random() * 0.3, to: 1 },
        duration: 300 + Math.random() * 700,
        yoyo: true,
        repeat: -1,
        delay: Math.random() * 1000,
      });
    }
  } else {
    const flowerColors = [0xff6b6b, 0xff9ff3, 0xffffff, 0xffdd57, 0xff8a65, 0x74b9ff];
    for (let i = 1; i < archPoints.length - 1; i += 3) {
      const pt = archPoints[i];
      const color = flowerColors[i % flowerColors.length];
      const fs = 3.5 + Math.random() * 2;

      for (let petal = 0; petal < 5; petal++) {
        const angle = (petal / 5) * Math.PI * 2;
        const px = pt.x + Math.cos(angle) * fs * 0.5;
        const py = pt.y + Math.sin(angle) * fs * 0.5;
        gfx.fillStyle(color, 0.95);
        gfx.fillCircle(px, py, fs * 0.45);
      }
      gfx.fillStyle(0xffee88, 1);
      gfx.fillCircle(pt.x, pt.y, fs * 0.3);
    }
  }

  const hangPoints = [
    archPoints[2],
    archPoints[3],
    archPoints[steps - 2],
    archPoints[steps - 3],
  ];
  hangPoints.forEach((pt) => {
    const vineLen = 10 + Math.random() * 12;
    gfx.lineStyle(1.5, ctx.p.canopyHL, 0.7);
    gfx.beginPath();
    gfx.moveTo(pt.x, pt.y);
    gfx.lineTo(pt.x + (Math.random() - 0.5) * 6, pt.y + vineLen);
    gfx.strokePath();

    gfx.fillStyle(ctx.p.vineLeaf, 0.8);
    gfx.fillEllipse(pt.x + (Math.random() - 0.5) * 6, pt.y + vineLen, 4, 3);
  });
}
