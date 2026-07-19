import Phaser from 'phaser';

export interface Palette {
  skyTop: [number, number, number];
  skyBot: [number, number, number];
  cloud: number | null;
  birdColor: number | null;
  distTrunk: number;
  distCanopy: number;
  hills: [number, number, number];
  groundBase: number;
  groundPatch: number;
  groundTuft: number;
  pathBase: number;
  pathEdge: number;
  pathLine: number;
  fencePost: number;
  fenceRail: number;
  fenceDark: number;
  trunkBase: number;
  trunkLight: number;
  trunkBark: number;
  trunkRoot: number;
  canopy: number;
  canopyShadow: number;
  canopyHL: number;
  bushBody: number;
  bushShadow: number;
  bushHL: number;
  flowerStem: number;
  flowerLeaf: number;
  grassColors: number[];
  vineStem: number;
  vineLeaf: number;
  mushroomStem: number;
  mushroomCap: number;
}

export interface SceneContext extends Phaser.Scene {
  p: Palette;
  _dark: boolean;
  _treeGfx: Phaser.GameObjects.Graphics[];
  _animals: Array<{
    sprite: Phaser.GameObjects.Sprite;
    wanderTween: Phaser.Tweens.Tween | null;
  }>;
  _transitioning: boolean;
  _reducedMotion: boolean;
  _onDocumentClick: (() => void) | null;
  flowers: Phaser.GameObjects.Graphics[];
  smileTimer: Phaser.Time.TimerEvent | null;
  addTween(
    config: Phaser.Types.Tweens.TweenBuilderConfig,
  ): Phaser.Tweens.Tween | null;
  addTimedEvent(
    config: Phaser.Types.Time.TimerEventConfig,
  ): Phaser.Time.TimerEvent | null;
  addDelayedCall(
    delay: number,
    callback: () => void,
  ): Phaser.Time.TimerEvent | null;
}

export interface SkyGradient {
  top: [number, number, number];
  bot: [number, number, number];
}

export const DAY: Palette = {
  skyTop: [100, 170, 235],
  skyBot: [155, 220, 245],
  cloud: 0xffffff,
  birdColor: 0xffffff,
  distTrunk: 0x245a16,
  distCanopy: 0x1a5c10,
  hills: [0x2d6b1a, 0x3a8522, 0x4a9a2e],
  groundBase: 0x3a7d22,
  groundPatch: 0x4a8d2e,
  groundTuft: 0x326e1c,
  pathBase: 0xc4a45a,
  pathEdge: 0xb8973e,
  pathLine: 0xd4b46a,
  fencePost: 0x8b6914,
  fenceRail: 0xa0722a,
  fenceDark: 0x6b4513,
  trunkBase: 0x5c3a0e,
  trunkLight: 0x6b4513,
  trunkBark: 0x4a3008,
  trunkRoot: 0x5c3a0e,
  canopy: 0x2a6e1a,
  canopyShadow: 0x1f5a12,
  canopyHL: 0x3a8c24,
  bushBody: 0x2d6b1a,
  bushShadow: 0x1f5a12,
  bushHL: 0x3a9a28,
  flowerStem: 0x2d6b1a,
  flowerLeaf: 0x3a8522,
  grassColors: [0x4a9a2e, 0x3a8522, 0x56b030, 0x2d8b1a],
  vineStem: 0x2d8b1a,
  vineLeaf: 0x3a9a28,
  mushroomStem: 0xf5e6c8,
  mushroomCap: 0xcc3333,
};

export const NIGHT: Palette = {
  skyTop: [12, 16, 42],
  skyBot: [25, 35, 70],
  cloud: null,
  birdColor: null,
  distTrunk: 0x0c2808,
  distCanopy: 0x0a2006,
  hills: [0x0f2a08, 0x14300c, 0x183810],
  groundBase: 0x1a3510,
  groundPatch: 0x1f3e15,
  groundTuft: 0x142e0c,
  pathBase: 0x3a2a15,
  pathEdge: 0x2e2010,
  pathLine: 0x4a3520,
  fencePost: 0x5a4010,
  fenceRail: 0x6a4a18,
  fenceDark: 0x3a2808,
  trunkBase: 0x2a1a08,
  trunkLight: 0x3a2510,
  trunkBark: 0x1a0e04,
  trunkRoot: 0x2a1a08,
  canopy: 0x0f2208,
  canopyShadow: 0x0a1a06,
  canopyHL: 0x1a3510,
  bushBody: 0x0f2208,
  bushShadow: 0x0a1a06,
  bushHL: 0x1a3012,
  flowerStem: 0x142e0c,
  flowerLeaf: 0x1a3510,
  grassColors: [0x183810, 0x142e0c, 0x1f4015, 0x14300c],
  vineStem: 0x142e0c,
  vineLeaf: 0x1a3012,
  mushroomStem: 0x8a7a60,
  mushroomCap: 0x882222,
};

export const SUNSET_SKY: SkyGradient = {
  top: [220, 100, 40],
  bot: [255, 170, 80],
};
export const DAWN_SKY: SkyGradient = {
  top: [70, 90, 150],
  bot: [170, 140, 180],
};
