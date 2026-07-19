import { describe, it, expect, afterEach, vi } from 'vitest';
import { messages, detectLocale } from '../i18n';

function keyPaths(obj: Record<string, unknown>, prefix = ''): string[] {
  return Object.entries(obj).flatMap(([key, value]) => {
    const path = prefix ? `${prefix}.${key}` : key;
    if (value && typeof value === 'object') {
      return keyPaths(value as Record<string, unknown>, path);
    }
    return [path];
  });
}

describe('locale catalogs', () => {
  const reference = keyPaths(messages.ko).sort();

  it.each(Object.keys(messages))('%s has the same keys as ko', (locale) => {
    const paths = keyPaths(
      messages[locale as keyof typeof messages] as Record<string, unknown>,
    ).sort();
    expect(paths).toEqual(reference);
  });
});

describe('detectLocale', () => {
  afterEach(() => {
    localStorage.clear();
    vi.unstubAllGlobals();
  });

  it('prefers a saved valid locale', () => {
    localStorage.setItem('wedding-locale', 'en');
    vi.stubGlobal('navigator', { language: 'ko-KR' });
    expect(detectLocale()).toBe('en');
  });

  it('ignores a saved unknown locale', () => {
    localStorage.setItem('wedding-locale', 'xx');
    vi.stubGlobal('navigator', { language: 'ko-KR' });
    expect(detectLocale()).toBe('ko');
  });

  it('detects English browsers', () => {
    vi.stubGlobal('navigator', { language: 'en-US' });
    expect(detectLocale()).toBe('en');
  });

  it('falls back to Korean for other languages', () => {
    vi.stubGlobal('navigator', { language: 'ja' });
    expect(detectLocale()).toBe('ko');
  });
});
