import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createMockScene } from './mock-scene';
import { drawFireflies, drawButterflies, drawAnimals, drawCharacters } from '../draw/creatures';

beforeEach(() => {
  if (typeof globalThis.document === 'undefined') {
    (globalThis as any).document = {
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    };
  } else {
    vi.spyOn(document, 'addEventListener').mockImplementation(() => {});
  }
});

const W = 400;
const H = 800;

describe('creatures draw functions', () => {
  describe('drawFireflies', () => {
    it('creates 10 firefly glow graphics with tweens', () => {
      const ctx = createMockScene(true);
      drawFireflies(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalledTimes(10);
      expect(ctx.addTween).toHaveBeenCalledTimes(10);
    });
  });

  describe('drawButterflies', () => {
    it('creates 4 butterflies with tweens', () => {
      const ctx = createMockScene();
      drawButterflies(ctx, W, H);
      expect(ctx.add.graphics).toHaveBeenCalledTimes(4);
      expect(ctx.addTween).toHaveBeenCalledTimes(4);
    });
  });

  describe('drawAnimals', () => {
    it('spawns 3 animals (1 chicken + 2 chicks)', () => {
      const ctx = createMockScene();
      drawAnimals(ctx, W, H);
      expect(ctx._animals).toHaveLength(3);
      expect(ctx.add.sprite).toHaveBeenCalledTimes(3);
    });

    it('applies night tint in dark mode', () => {
      const ctx = createMockScene(true);
      drawAnimals(ctx, W, H);
      const sprite = (ctx.add.sprite as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(sprite.setTint).toHaveBeenCalledWith(0x223322);
    });

    it('does not apply tint in day mode', () => {
      const ctx = createMockScene(false);
      drawAnimals(ctx, W, H);
      const sprite = (ctx.add.sprite as ReturnType<typeof vi.fn>).mock.results[0].value;
      expect(sprite.setTint).not.toHaveBeenCalled();
    });

    it('cleans up previous animals on redraw', () => {
      const ctx = createMockScene();
      const oldSprite = { destroy: vi.fn() };
      ctx._animals = [{ sprite: oldSprite as any, wanderTween: null }];
      drawAnimals(ctx, W, H);
      expect(oldSprite.destroy).toHaveBeenCalled();
    });
  });

  describe('drawCharacters', () => {
    it('creates character images and registers click listener', () => {
      const ctx = createMockScene();
      drawCharacters(ctx, W, H);
      expect(ctx.add.image).toHaveBeenCalledTimes(3);
      expect(ctx._onDocumentClick).not.toBeNull();
    });

    it('schedules blink animation', () => {
      const ctx = createMockScene();
      drawCharacters(ctx, W, H);
      expect(ctx.addDelayedCall).toHaveBeenCalled();
    });
  });
});
