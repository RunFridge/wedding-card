import { ref } from 'vue';

const STORAGE_KEY = 'wedding-theme';
const HALF_DURATION = 750;

function isNightByHour() {
  const h = new Date().getHours();
  return h >= 18 || h < 6;
}

function resolveInitialTheme(): boolean {
  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved) return saved === 'dark';
  return isNightByHour();
}

const isDark = ref(resolveInitialTheme());
const transitioning = ref(false);

// Apply initial class synchronously — no event, no transition
(() => {
  const html = document.documentElement;
  if (isDark.value) {
    html.classList.add('dark');
  } else {
    html.classList.remove('dark');
  }
})();

function toggle() {
  if (transitioning.value) return;
  transitioning.value = true;

  const newDark = !isDark.value;

  // Dispatch immediately so Phaser starts its fade-out
  window.dispatchEvent(
    new CustomEvent('theme-change', {
      detail: { dark: newDark, duration: HALF_DURATION * 2 },
    }),
  );

  // At midpoint: swap CSS class + persist
  setTimeout(() => {
    const html = document.documentElement;
    if (newDark) {
      html.classList.add('dark');
    } else {
      html.classList.remove('dark');
    }
    isDark.value = newDark;
    localStorage.setItem(STORAGE_KEY, newDark ? 'dark' : 'light');
  }, HALF_DURATION);

  // At end: unlock toggle
  setTimeout(() => {
    transitioning.value = false;
  }, HALF_DURATION * 2);
}

export function useTheme() {
  return { isDark, toggle, transitioning };
}
