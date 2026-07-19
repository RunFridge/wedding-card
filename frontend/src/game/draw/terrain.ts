import { SceneContext } from './palettes';

export function drawDistantTrees(
  ctx: SceneContext,
  w: number,
  h: number,
): void {
  const gfx = ctx.add.graphics();
  const hillBaseY = h * 0.48;

  for (let x = 0; x < w; x += 20 + Math.random() * 15) {
    const canopyH = 12 + Math.random() * 10;
    const canopyR = 8 + Math.random() * 5;

    gfx.fillStyle(ctx.p.distTrunk, 0.5);
    gfx.fillRect(x - 1.5, hillBaseY - canopyH, 3, canopyH);
    gfx.fillStyle(ctx.p.distCanopy, 0.6);
    gfx.fillCircle(x, hillBaseY - canopyH, canopyR);
  }
}

export function drawHills(ctx: SceneContext, w: number, h: number): void {
  const hillLayers = [
    {
      baseY: h * 0.45,
      amplitude: 30,
      frequency: 0.004,
      color: ctx.p.hills[0],
      offset: 0,
    },
    {
      baseY: h * 0.5,
      amplitude: 25,
      frequency: 0.006,
      color: ctx.p.hills[1],
      offset: 100,
    },
    {
      baseY: h * 0.55,
      amplitude: 20,
      frequency: 0.008,
      color: ctx.p.hills[2],
      offset: 200,
    },
  ];

  hillLayers.forEach(({ baseY, amplitude, frequency, color, offset }) => {
    const gfx = ctx.add.graphics();
    gfx.fillStyle(color, 1);

    const points: { x: number; y: number }[] = [{ x: 0, y: h }];
    for (let x = 0; x <= w; x += 4) {
      const y =
        baseY +
        Math.sin((x + offset) * frequency) * amplitude +
        Math.sin((x + offset) * frequency * 2.3) * (amplitude * 0.4);
      points.push({ x, y });
    }
    points.push({ x: w, y: h });

    gfx.beginPath();
    gfx.moveTo(points[0].x, points[0].y);
    points.forEach((p) => gfx.lineTo(p.x, p.y));
    gfx.closePath();
    gfx.fillPath();
  });
}

export function drawGround(ctx: SceneContext, w: number, h: number): void {
  const groundY = h * 0.62;
  const gfx = ctx.add.graphics();

  gfx.fillStyle(ctx.p.groundBase, 1);
  gfx.fillRect(0, groundY, w, h - groundY);

  gfx.fillStyle(ctx.p.groundPatch, 1);
  for (let x = 0; x < w; x += 24) {
    const patchW = 10 + Math.random() * 14;
    const patchY = groundY + 8 + Math.random() * (h - groundY - 30);
    gfx.fillRect(x, patchY, patchW, 4);
  }

  gfx.fillStyle(ctx.p.groundTuft, 1);
  for (let x = 0; x < w; x += 8) {
    const tuftH = 4 + Math.random() * 8;
    gfx.fillRect(x, groundY - tuftH, 3, tuftH);
  }
}

export function drawPath(ctx: SceneContext, w: number, h: number): void {
  const gfx = ctx.add.graphics();
  const pathCenterX = w * 0.5;
  const pathTopY = h * 0.62;
  const pathW = 100;

  gfx.fillStyle(ctx.p.pathBase, 1);
  gfx.fillRect(pathCenterX - pathW / 2, pathTopY, pathW, h - pathTopY);

  gfx.fillStyle(ctx.p.pathEdge, 1);
  gfx.fillRect(pathCenterX - pathW / 2, pathTopY, 3, h - pathTopY);
  gfx.fillRect(pathCenterX + pathW / 2 - 3, pathTopY, 3, h - pathTopY);

  gfx.fillStyle(ctx.p.pathLine, 0.4);
  for (let y = pathTopY; y < h; y += 16) {
    gfx.fillRect(pathCenterX - pathW / 2 + 5, y, pathW - 10, 2);
  }
}

export function drawFence(ctx: SceneContext, w: number, h: number): void {
  const fenceY = h * 0.6;
  const postSpacing = 48;
  const postWidth = 8;
  const postHeight = 36;
  const railHeight = 4;

  const gfx = ctx.add.graphics();

  gfx.fillStyle(ctx.p.fencePost, 1);
  for (let x = 12; x < w; x += postSpacing) {
    gfx.fillRect(x, fenceY - postHeight, postWidth, postHeight);

    gfx.fillStyle(ctx.p.fenceRail, 1);
    gfx.fillRect(x + 2, fenceY - postHeight + 4, postWidth - 4, postHeight - 4);

    gfx.fillStyle(ctx.p.fenceDark, 1);
    gfx.fillRect(x, fenceY - postHeight, postWidth, 3);
    gfx.fillRect(x - 2, fenceY - postHeight - 2, postWidth + 4, 3);

    gfx.fillStyle(ctx.p.fencePost, 1);
  }

  gfx.fillStyle(ctx.p.fenceRail, 1);
  gfx.fillRect(0, fenceY - postHeight * 0.7, w, railHeight);
  gfx.fillRect(0, fenceY - postHeight * 0.35, w, railHeight);

  gfx.fillStyle(ctx.p.fenceDark, 1);
  gfx.fillRect(0, fenceY - postHeight * 0.7 + railHeight, w, 2);
  gfx.fillRect(0, fenceY - postHeight * 0.35 + railHeight, w, 2);
}
