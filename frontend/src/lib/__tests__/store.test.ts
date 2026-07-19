import { describe, it, expect, beforeEach } from 'vitest';
import {
  getGameCompletion,
  setGameCompletion,
  clearGameCompletion,
  getAchievements,
  setAchievements,
  clearAchievements,
  getHallOfFameSubmitted,
  setHallOfFameSubmitted,
  clearHallOfFameSubmitted,
  clearAll,
} from '../store';

beforeEach(() => {
  localStorage.clear();
});

describe('store — game completion', () => {
  it('returns null when no completion stored', () => {
    expect(getGameCompletion()).toBeNull();
  });

  it('stores and retrieves a timestamp', () => {
    setGameCompletion();
    const val = getGameCompletion();
    expect(val).toBeTypeOf('number');
    expect(val!).toBeGreaterThan(0);
  });

  it('clearGameCompletion removes the value', () => {
    setGameCompletion();
    clearGameCompletion();
    expect(getGameCompletion()).toBeNull();
  });

  it('handles corrupted data gracefully', () => {
    localStorage.setItem('_wcg', 'not-valid-base64!!!');
    expect(getGameCompletion()).toBeNull();
  });
});

describe('store — achievements', () => {
  it('returns empty array when nothing stored', () => {
    expect(getAchievements()).toEqual([]);
  });

  it('stores and retrieves achievement IDs', () => {
    setAchievements(['gold', 'silver', 'bronze']);
    expect(getAchievements()).toEqual(['gold', 'silver', 'bronze']);
  });

  it('handles empty array', () => {
    setAchievements([]);
    expect(getAchievements()).toEqual([]);
  });

  it('clearAchievements removes the value', () => {
    setAchievements(['gold']);
    clearAchievements();
    expect(getAchievements()).toEqual([]);
  });

  it('handles corrupted data gracefully', () => {
    localStorage.setItem('_wca', 'garbage');
    expect(getAchievements()).toEqual([]);
  });
});

describe('store — hall of fame', () => {
  it('getHallOfFameSubmitted returns false when not set', () => {
    expect(getHallOfFameSubmitted()).toBe(false);
  });

  it('setHallOfFameSubmitted stores encoded timestamp', () => {
    const before = Date.now();
    setHallOfFameSubmitted();
    const raw = localStorage.getItem('_wch');
    expect(raw).not.toBeNull();
    expect(raw).toMatch(/^[A-Za-z0-9+/=]+$/);
    const decoded = JSON.parse(decodeURIComponent(atob(raw!)));
    expect(decoded).toBeTypeOf('number');
    expect(decoded).toBeGreaterThanOrEqual(before);
    expect(decoded).toBeLessThanOrEqual(Date.now());
  });

  it('getHallOfFameSubmitted returns true after set', () => {
    setHallOfFameSubmitted();
    expect(getHallOfFameSubmitted()).toBe(true);
  });

  it('clearHallOfFameSubmitted removes the item', () => {
    setHallOfFameSubmitted();
    clearHallOfFameSubmitted();
    expect(getHallOfFameSubmitted()).toBe(false);
    expect(localStorage.getItem('_wch')).toBeNull();
  });
});

describe('store — clearAll', () => {
  it('clears game, achievements, and hall of fame', () => {
    setGameCompletion();
    setAchievements(['gold', 'theme']);
    setHallOfFameSubmitted();
    clearAll();
    expect(getGameCompletion()).toBeNull();
    expect(getAchievements()).toEqual([]);
    expect(getHallOfFameSubmitted()).toBe(false);
  });

  it('clearAll clears hall of fame too', () => {
    setHallOfFameSubmitted();
    expect(getHallOfFameSubmitted()).toBe(true);
    clearAll();
    expect(getHallOfFameSubmitted()).toBe(false);
    expect(localStorage.getItem('_wch')).toBeNull();
  });
});

describe('store — encoding', () => {
  it('data is base64 encoded in localStorage', () => {
    setAchievements(['test']);
    const raw = localStorage.getItem('_wca');
    expect(raw).not.toBeNull();
    expect(raw).not.toContain('test');
    expect(raw).toMatch(/^[A-Za-z0-9+/=]+$/);
  });

  it('round-trips unicode characters', () => {
    setAchievements(['금메달', '은메달', '동메달']);
    expect(getAchievements()).toEqual(['금메달', '은메달', '동메달']);
  });
});
