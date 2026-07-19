const K_GAME = '\x5f\x77\x63\x67'; // _wcg
const K_ACH = '\x5f\x77\x63\x61'; // _wca
const K_HOF = '\x5f\x77\x63\x68'; // _wch

function encode(data: unknown): string {
  return btoa(encodeURIComponent(JSON.stringify(data)));
}

function decode<T>(raw: string): T | null {
  try {
    return JSON.parse(decodeURIComponent(atob(raw)));
  } catch {
    return null;
  }
}

export function getGameCompletion(): number | null {
  const raw = localStorage.getItem(K_GAME);
  if (!raw) return null;
  const val = decode<number>(raw);
  return typeof val === 'number' && !isNaN(val) ? val : null;
}

export function setGameCompletion() {
  localStorage.setItem(K_GAME, encode(Date.now()));
}

export function clearGameCompletion() {
  localStorage.removeItem(K_GAME);
}

export function getAchievements(): string[] {
  const raw = localStorage.getItem(K_ACH);
  if (!raw) return [];
  return decode<string[]>(raw) ?? [];
}

export function setAchievements(ids: string[]) {
  localStorage.setItem(K_ACH, encode(ids));
}

export function clearAchievements() {
  localStorage.removeItem(K_ACH);
}

export function getHallOfFameSubmitted(): boolean {
  const raw = localStorage.getItem(K_HOF);
  if (!raw) return false;
  const val = decode<number>(raw);
  return typeof val === 'number' && !isNaN(val);
}

export function setHallOfFameSubmitted() {
  localStorage.setItem(K_HOF, encode(Date.now()));
}

export function clearHallOfFameSubmitted() {
  localStorage.removeItem(K_HOF);
}

export function clearAll() {
  clearGameCompletion();
  clearAchievements();
  clearHallOfFameSubmitted();
}
