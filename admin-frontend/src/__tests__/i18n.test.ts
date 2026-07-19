import { describe, it, expect } from 'vitest';
import ko from '../locales/ko.json';
import en from '../locales/en.json';

const messages: Record<string, Record<string, unknown>> = { ko, en };

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
  const reference = keyPaths(ko).sort();

  it.each(Object.keys(messages))('%s has the same keys as ko', (locale) => {
    expect(keyPaths(messages[locale]).sort()).toEqual(reference);
  });
});
