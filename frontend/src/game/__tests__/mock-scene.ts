import { vi } from 'vitest';
import type { SceneContext, Palette } from '../draw/palettes';
import { DAY, NIGHT } from '../draw/palettes';

function createMockGraphics() {
  const gfx: Record<string, any> = {
    fillStyle: vi.fn().mockReturnThis(),
    fillRect: vi.fn().mockReturnThis(),
    fillCircle: vi.fn().mockReturnThis(),
    fillEllipse: vi.fn().mockReturnThis(),
    fillPath: vi.fn().mockReturnThis(),
    lineStyle: vi.fn().mockReturnThis(),
    beginPath: vi.fn().mockReturnThis(),
    moveTo: vi.fn().mockReturnThis(),
    lineTo: vi.fn().mockReturnThis(),
    closePath: vi.fn().mockReturnThis(),
    strokePath: vi.fn().mockReturnThis(),
    clear: vi.fn().mockReturnThis(),
    setPosition: vi.fn().mockReturnThis(),
    setDepth: vi.fn().mockReturnThis(),
    setAlpha: vi.fn().mockReturnThis(),
    destroy: vi.fn(),
  };
  return gfx;
}

function createMockSprite() {
  return {
    setScale: vi.fn().mockReturnThis(),
    setDepth: vi.fn().mockReturnThis(),
    setFrame: vi.fn().mockReturnThis(),
    setTint: vi.fn().mockReturnThis(),
    setFlipX: vi.fn().mockReturnThis(),
    setVisible: vi.fn().mockReturnThis(),
    play: vi.fn().mockReturnThis(),
    stop: vi.fn().mockReturnThis(),
    destroy: vi.fn(),
    x: 200,
    preFX: { addGlow: vi.fn() },
  };
}

export function createMockScene(dark = false): SceneContext {
  const graphics = createMockGraphics();
  const sprite = createMockSprite();

  return {
    p: dark ? { ...NIGHT } : { ...DAY },
    _dark: dark,
    _treeGfx: [],
    _animals: [],
    _transitioning: false,
    _reducedMotion: false,
    _onDocumentClick: null,
    flowers: [],
    smileTimer: null,
    add: {
      graphics: vi.fn(() => graphics),
      sprite: vi.fn(() => sprite),
      image: vi.fn(() => sprite),
    },
    addTween: vi.fn(() => null),
    addTimedEvent: vi.fn(() => null),
    addDelayedCall: vi.fn(() => null),
  } as unknown as SceneContext;
}
