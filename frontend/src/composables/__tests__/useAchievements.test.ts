import { describe, it, expect, beforeEach, vi } from 'vitest';
import ko from '../../locales/ko.json';
import en from '../../locales/en.json';

beforeEach(() => {
  localStorage.clear();
  vi.resetModules();
});

async function freshModule() {
  const mod = await import('../useAchievements');
  return mod;
}

describe('useAchievements', () => {
  it('exports ACHIEVEMENTS with 11 entries', async () => {
    const { ACHIEVEMENTS } = await freshModule();
    expect(ACHIEVEMENTS).toHaveLength(11);
  });

  it('each achievement has id, icon, and catalog entries in every locale', async () => {
    const { ACHIEVEMENTS } = await freshModule();
    const catalogs = { ko, en };
    for (const a of ACHIEVEMENTS) {
      expect(a.id).toBeTypeOf('string');
      expect(a.icon).toBeTypeOf('string');
      for (const [localeName, catalog] of Object.entries(catalogs)) {
        const entry = (
          catalog.achievements as Record<
            string,
            { name: string; description: string }
          >
        )[a.id];
        expect(entry, `${localeName} achievements.${a.id}`).toBeDefined();
        expect(entry.name).toBeTypeOf('string');
        expect(entry.description).toBeTypeOf('string');
      }
    }
  });

  it('all achievement IDs are unique', async () => {
    const { ACHIEVEMENTS } = await freshModule();
    const ids = ACHIEVEMENTS.map((a) => a.id);
    expect(new Set(ids).size).toBe(ids.length);
  });

  it('wedding_day is marked as secret', async () => {
    const { ACHIEVEMENTS } = await freshModule();
    const wd = ACHIEVEMENTS.find((a) => a.id === 'wedding_day');
    expect(wd?.secret).toBe(true);
  });

  it('starts with empty earned set', async () => {
    const { useAchievements } = await freshModule();
    const { earned } = useAchievements();
    expect(earned.value.size).toBe(0);
  });

  it('award adds to earned set', async () => {
    const { useAchievements } = await freshModule();
    const { earned, award } = useAchievements();
    award('gallery');
    expect(earned.value.has('gallery')).toBe(true);
  });

  it('awarding same id twice does not duplicate', async () => {
    const { useAchievements } = await freshModule();
    const { earned, award } = useAchievements();
    award('gold');
    award('gold');
    expect(earned.value.size).toBe(1);
  });

  it('award persists to localStorage', async () => {
    const { useAchievements } = await freshModule();
    const { award } = useAchievements();
    award('theme');
    const raw = localStorage.getItem('_wca');
    expect(raw).not.toBeNull();
  });

  it('clearAll removes all earned achievements', async () => {
    const { useAchievements } = await freshModule();
    const { earned, award, clearAll } = useAchievements();
    award('gold');
    award('silver');
    clearAll();
    expect(earned.value.size).toBe(0);
  });

  it('loads previously saved achievements', async () => {
    const { setAchievements } = await import('../../lib/store');
    setAchievements(['gallery', 'guestbook', 'photo']);
    const { useAchievements } = await freshModule();
    const { earned } = useAchievements();
    expect(earned.value.size).toBe(3);
    expect(earned.value.has('gallery')).toBe(true);
    expect(earned.value.has('guestbook')).toBe(true);
    expect(earned.value.has('photo')).toBe(true);
  });

  it('award sets currentToast', async () => {
    const { useAchievements } = await freshModule();
    const { currentToast, award } = useAchievements();
    award('heart');
    expect(currentToast.value).not.toBeNull();
    expect(currentToast.value?.id).toBe('heart');
  });

  it('allEarned is true when all 11 are earned', async () => {
    const { useAchievements, ACHIEVEMENTS } = await freshModule();
    const { allEarned, award } = useAchievements();
    expect(allEarned.value).toBe(false);
    for (const a of ACHIEVEMENTS) {
      award(a.id);
    }
    expect(allEarned.value).toBe(true);
  });

  it('allEarned is false when not all earned', async () => {
    const { useAchievements, ACHIEVEMENTS } = await freshModule();
    const { allEarned, award } = useAchievements();
    for (let i = 0; i < ACHIEVEMENTS.length - 1; i++) {
      award(ACHIEVEMENTS[i].id);
    }
    expect(allEarned.value).toBe(false);
  });

  const NON_PODIUM_REQUIRED = [
    'gallery',
    'guestbook',
    'photo',
    'heart',
    'gallery_all',
    'theme',
    'speed',
    'wedding_day',
  ];

  it('allEarned is true with bronze + wedding_day + all 7 non-medals', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.forEach(award);
    award('bronze');
    expect(allEarned.value).toBe(true);
  });

  it('allEarned is true with only silver podium (no gold, no bronze)', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.forEach(award);
    award('silver');
    expect(allEarned.value).toBe(true);
  });

  it('allEarned is true with only gold podium', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.forEach(award);
    award('gold');
    expect(allEarned.value).toBe(true);
  });

  it('allEarned is false without wedding_day even with all 3 medals + 7 non-medals', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.filter((id) => id !== 'wedding_day').forEach(award);
    award('gold');
    award('silver');
    award('bronze');
    expect(allEarned.value).toBe(false);
  });

  it('allEarned is false without any podium even with wedding_day + 7 non-medals', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.forEach(award);
    expect(allEarned.value).toBe(false);
  });

  it('allEarned is false if missing one of the 7 non-medals', async () => {
    const { useAchievements } = await freshModule();
    const { allEarned, award } = useAchievements();
    NON_PODIUM_REQUIRED.filter((id) => id !== 'photo').forEach(award);
    award('gold');
    expect(allEarned.value).toBe(false);
  });

  it('clearAll resets earned and allEarned', async () => {
    const { useAchievements, ACHIEVEMENTS } = await freshModule();
    const { earned, allEarned, award, clearAll } = useAchievements();
    for (const a of ACHIEVEMENTS) {
      award(a.id);
    }
    expect(allEarned.value).toBe(true);
    clearAll();
    expect(earned.value.size).toBe(0);
    expect(allEarned.value).toBe(false);
  });

  it('showHofDialog triggers when all earned and hof not submitted', async () => {
    vi.useFakeTimers();
    const { useAchievements, ACHIEVEMENTS } = await freshModule();
    const { showHofDialog, award } = useAchievements();
    expect(showHofDialog.value).toBe(false);
    for (const a of ACHIEVEMENTS) {
      award(a.id);
    }
    vi.advanceTimersByTime(4000);
    expect(showHofDialog.value).toBe(true);
    vi.useRealTimers();
  });
});
