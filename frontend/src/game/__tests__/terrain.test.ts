import { describe, it, expect, vi } from 'vitest';
import { createMockScene } from './mock-scene';
import { drawDistantTrees, drawHills, drawGround, drawPath, drawFence } from '../draw/terrain';

const W = 400;
const H = 800;

describe('terrain draw functions', () => {
  describe('drawDistantTrees', () => {
    it('draws distant tree trunks and canopies', () => {
      const ctx = createMockScene();
      drawDistantTrees(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillRect).toHaveBeenCalled();
      expect(gfx.fillCircle).toHaveBeenCalled();
    });
  });

  describe('drawHills', () => {
    it('draws 3 hill layers', () => {
      const ctx = createMockScene();
      drawHills(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalledTimes(3);
    });
  });

  describe('drawGround', () => {
    it('fills ground base, patches, and tufts', () => {
      const ctx = createMockScene();
      drawGround(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillStyle).toHaveBeenCalled();
      expect(gfx.fillRect).toHaveBeenCalled();
    });
  });

  describe('drawPath', () => {
    it('draws center path with edges and lines', () => {
      const ctx = createMockScene();
      drawPath(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillRect).toHaveBeenCalled();
    });
  });

  describe('drawFence', () => {
    it('draws posts, rails, and dark edges', () => {
      const ctx = createMockScene();
      drawFence(ctx, W, H);
      const gfx = (ctx.add.graphics as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(gfx.fillRect).toHaveBeenCalled();
    });
  });
});
