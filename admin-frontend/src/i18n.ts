import { createI18n } from 'vue-i18n';
import ko from './locales/ko.json';
import en from './locales/en.json';

const STORAGE_KEY = 'wedding-locale';

export const messages = { ko, en };
export type Locale = keyof typeof messages;

export function detectLocale(): Locale {
  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved && saved in messages) return saved as Locale;
  return navigator.language.startsWith('en') ? 'en' : 'ko';
}

export const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'ko',
  messages,
});

export function setLocale(locale: Locale) {
  i18n.global.locale.value = locale;
  localStorage.setItem(STORAGE_KEY, locale);
  document.documentElement.lang = locale;
}

document.documentElement.lang = i18n.global.locale.value;
