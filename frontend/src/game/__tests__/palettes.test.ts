import { describe, it, expect } from 'vitest';
import { DAY, NIGHT, SUNSET_SKY, DAWN_SKY } from '../draw/palettes';
import type { Palette } from '../draw/palettes';

function paletteKeys(): (keyof Palette)[] {
  return [
    'skyTop', 'skyBot', 'cloud', 'birdColor',
    'distTrunk', 'distCanopy', 'hills',
    'groundBase', 'groundPatch', 'groundTuft',
    'pathBase', 'pathEdge', 'pathLine',
    'fencePost', 'fenceRail', 'fenceDark',
    'trunkBase', 'trunkLight', 'trunkBark', 'trunkRoot',
    'canopy', 'canopyShadow', 'canopyHL',
    'bushBody', 'bushShadow', 'bushHL',
    'flowerStem', 'flowerLeaf', 'grassColors',
    'vineStem', 'vineLeaf',
    'mushroomStem', 'mushroomCap',
  ];
}

describe('palettes', () => {
  it('DAY and NIGHT have identical keys', () => {
    expect(Object.keys(DAY).sort()).toEqual(Object.keys(NIGHT).sort());
  });

  it('DAY has all required palette keys', () => {
    for (const key of paletteKeys()) {
      expect(DAY).toHaveProperty(key);
    }
  });

  it('NIGHT has all required palette keys', () => {
    for (const key of paletteKeys()) {
      expect(NIGHT).toHaveProperty(key);
    }
  });

  it('sky colors are RGB tuples of length 3', () => {
    for (const palette of [DAY, NIGHT]) {
      expect(palette.skyTop).toHaveLength(3);
      expect(palette.skyBot).toHaveLength(3);
      palette.skyTop.forEach((v) => expect(v).toBeGreaterThanOrEqual(0));
      palette.skyBot.forEach((v) => expect(v).toBeLessThanOrEqual(255));
    }
  });

  it('hills are tuples of length 3', () => {
    expect(DAY.hills).toHaveLength(3);
    expect(NIGHT.hills).toHaveLength(3);
  });

  it('grassColors has at least 1 color', () => {
    expect(DAY.grassColors.length).toBeGreaterThan(0);
    expect(NIGHT.grassColors.length).toBeGreaterThan(0);
  });

  it('DAY has non-null cloud and birdColor', () => {
    expect(DAY.cloud).not.toBeNull();
    expect(DAY.birdColor).not.toBeNull();
  });

  it('NIGHT has null cloud and birdColor', () => {
    expect(NIGHT.cloud).toBeNull();
    expect(NIGHT.birdColor).toBeNull();
  });

  it('SUNSET_SKY and DAWN_SKY have top and bot tuples', () => {
    expect(SUNSET_SKY.top).toHaveLength(3);
    expect(SUNSET_SKY.bot).toHaveLength(3);
    expect(DAWN_SKY.top).toHaveLength(3);
    expect(DAWN_SKY.bot).toHaveLength(3);
  });
});
