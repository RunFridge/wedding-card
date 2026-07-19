import { describe, it, expect, vi } from 'vitest';
import { createMockScene } from './mock-scene';
import { drawTrees, drawBushes, drawMushrooms, drawTallGrass, drawFlowers, drawFenceVines } from '../draw/vegetation';

const W = 400;
const H = 800;

describe('vegetation draw functions', () => {
  describe('drawTrees', () => {
    it('creates 4 tree graphics and stores refs in _treeGfx', () => {
      const ctx = createMockScene();
      drawTrees(ctx, W, H);
      expect(ctx._treeGfx).toHaveLength(4);
    });

    it('draws fruits in day mode only', () => {
      const dayCtx = createMockScene(false);
      const nightCtx = createMockScene(true);
      drawTrees(dayCtx, W, H);
      drawTrees(nightCtx, W, H);
      const dayGfx = (dayCtx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      const nightGfx = (nightCtx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(dayGfx.fillCircle.mock.calls.length).toBeGreaterThan(nightGfx.fillCircle.mock.calls.length);
    });
  });

  describe('drawBushes', () => {
    it('draws bush bodies and berries', () => {
      const ctx = createMockScene();
      drawBushes(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillCircle).toHaveBeenCalled();
    });
  });

  describe('drawMushrooms', () => {
    it('draws stems and caps for all spots', () => {
      const ctx = createMockScene();
      drawMushrooms(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillEllipse).toHaveBeenCalled();
    });

    it('draws glow in night mode', () => {
      const ctx = createMockScene(true);
      drawMushrooms(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      const glowCalls = gfx.fillStyle.mock.calls.filter(
        (c: number[]) => c[0] === 0xffaa44,
      );
      expect(glowCalls.length).toBeGreaterThan(0);
    });
  });

  describe('drawTallGrass', () => {
    it('draws blade shapes using paths', () => {
      const ctx = createMockScene();
      drawTallGrass(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.beginPath).toHaveBeenCalled();
      expect(gfx.fillPath).toHaveBeenCalled();
    });
  });

  describe('drawFlowers', () => {
    it('creates flower graphics and stores in ctx.flowers', () => {
      const ctx = createMockScene();
      drawFlowers(ctx, W, H);
      expect(ctx.flowers.length).toBe(35);
    });

    it('adds sway tween', () => {
      const ctx = createMockScene();
      drawFlowers(ctx, W, H);
      expect(ctx.addTween).toHaveBeenCalled();
    });
  });

  describe('drawFenceVines', () => {
    it('draws vine stems and leaves', () => {
      const ctx = createMockScene();
      drawFenceVines(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.strokePath).toHaveBeenCalled();
    });
  });
});
