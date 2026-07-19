import { ref, computed } from 'vue';
import {
  getAchievements,
  setAchievements,
  clearAchievements as clearStore,
  getHallOfFameSubmitted,
} from '../lib/store';

export interface Achievement {
  id: string;
  icon: string;
  secret?: boolean;
}

export const ACHIEVEMENTS: Achievement[] = [
  { id: 'gallery', icon: '🏆' },
  { id: 'gold', icon: '🥇' },
  { id: 'silver', icon: '🥈' },
  { id: 'bronze', icon: '🥉' },
  { id: 'guestbook', icon: '📖' },
  { id: 'photo', icon: '📸' },
  { id: 'heart', icon: '💕' },
  { id: 'gallery_all', icon: '🖼️' },
  { id: 'theme', icon: '🎨' },
  { id: 'speed', icon: '⚡' },
  { id: 'wedding_day', icon: '💒', secret: true },
];

const PODIUM_IDS = ['gold', 'silver', 'bronze'] as const;
const HOF_REQUIRED_IDS = [
  'gallery',
  'guestbook',
  'photo',
  'heart',
  'gallery_all',
  'theme',
  'speed',
  'wedding_day',
] as const;

const earned = ref<Set<string>>(new Set(getAchievements()));
const showHofDialog = ref(false);
const allEarned = computed(() => {
  const e = earned.value;
  const hasAllRequired = HOF_REQUIRED_IDS.every((id) => e.has(id));
  const hasPodium = PODIUM_IDS.some((id) => e.has(id));
  return hasAllRequired && hasPodium;
});
const toastQueue: Achievement[] = [];
const currentToast = ref<Achievement | null>(null);
let toastTimer: ReturnType<typeof setTimeout> | null = null;

function save() {
  setAchievements([...earned.value]);
}

function showNextToast() {
  const next = toastQueue.shift();
  if (!next) {
    currentToast.value = null;
    return;
  }
  currentToast.value = next;
  navigator.vibrate?.([30, 50, 30]);
  toastTimer = setTimeout(() => {
    currentToast.value = null;
    setTimeout(showNextToast, 300);
  }, 3000);
}

function award(id: string) {
  if (earned.value.has(id)) return;
  earned.value = new Set([...earned.value, id]);
  save();
  const a = ACHIEVEMENTS.find((a) => a.id === id);
  if (a) {
    toastQueue.push(a);
    if (!currentToast.value) showNextToast();
  }
  if (allEarned.value && !getHallOfFameSubmitted()) {
    setTimeout(() => {
      showHofDialog.value = true;
    }, 3500);
  }
}

function clearAll() {
  clearStore();
  earned.value = new Set();
  toastQueue.length = 0;
  currentToast.value = null;
  if (toastTimer) {
    clearTimeout(toastTimer);
    toastTimer = null;
  }
}

export function useAchievements() {
  return {
    earned,
    currentToast,
    award,
    clearAll,
    allEarned,
    showHofDialog,
    ACHIEVEMENTS,
  };
}
