import { describe, it, expect, vi } from 'vitest';
import { createMockScene } from './mock-scene';
import { drawSky, drawStars, drawMoon, drawClouds, drawSun, drawBirds } from '../draw/sky';

const W = 400;
const H = 800;

describe('sky draw functions', () => {
  describe('drawSky', () => {
    it('creates a graphics object and draws sky gradient', () => {
      const ctx = createMockScene();
      drawSky(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillStyle).toHaveBeenCalled();
      expect(gfx.fillRect).toHaveBeenCalled();
    });

    it('works with night palette', () => {
      const ctx = createMockScene(true);
      drawSky(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
    });
  });

  describe('drawStars', () => {
    it('creates static stars and twinkling stars with tweens', () => {
      const ctx = createMockScene(true);
      drawStars(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      expect(ctx.addTween).toHaveBeenCalled();
    });
  });

  describe('drawMoon', () => {
    it('draws moon with glow, body, and crescent', () => {
      const ctx = createMockScene(true);
      drawMoon(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillCircle).toHaveBeenCalled();
    });
  });

  describe('drawClouds', () => {
    it('creates cloud graphics with drift tweens', () => {
      const ctx = createMockScene();
      drawClouds(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      expect(ctx.addTween).toHaveBeenCalled();
    });
  });

  describe('drawSun', () => {
    it('creates sun glow and ray animation', () => {
      const ctx = createMockScene();
      drawSun(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      expect(ctx.addTween).toHaveBeenCalled();
    });
  });

  describe('drawBirds', () => {
    it('creates at least one bird with wing flap timer', () => {
      const ctx = createMockScene();
      drawBirds(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalled();
      expect(ctx.addTimedEvent).toHaveBeenCalled();
      expect(ctx.addDelayedCall).toHaveBeenCalled();
    });
  });
});
