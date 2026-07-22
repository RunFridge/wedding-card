<template>
  <div class="min-h-screen font-pixel">
    <div ref="phaserContainer" class="fixed inset-0 z-0"></div>

    <DemoRibbon v-if="DEMO_MODE" />

    <!-- Top-right buttons -->
    <div class="fixed top-3 right-3 z-50 flex gap-1.5">
      <button
        class="w-9 h-9 rounded-full bg-wood-dark/50 backdrop-blur flex items-center justify-center text-parchment/80 hover:text-parchment cursor-pointer transition-colors"
        :aria-label="t('app.share')"
        @pointerdown.stop
        @click="shareInvitation"
      >
        <Share2 class="w-4 h-4" />
      </button>
      <button
        class="w-9 h-9 rounded-full bg-wood-dark/50 backdrop-blur flex items-center justify-center text-parchment/80 hover:text-parchment cursor-pointer transition-colors"
        :aria-label="hideUI ? t('app.showUi') : t('app.hideUi')"
        @pointerdown.stop
        @click="hideUI = !hideUI"
      >
        <EyeOff v-if="hideUI" class="w-4 h-4" />
        <Eye v-else class="w-4 h-4" />
      </button>
      <button
        class="w-9 h-9 rounded-full bg-wood-dark/50 backdrop-blur flex items-center justify-center text-parchment/80 hover:text-parchment cursor-pointer transition-colors"
        :aria-label="isDark ? t('app.dayMode') : t('app.nightMode')"
        :disabled="transitioning"
        @pointerdown.stop
        @click="toggle"
      >
        <Sun v-if="isDark" class="w-4 h-4" />
        <Moon v-else class="w-4 h-4" />
      </button>
      <button
        class="w-9 h-9 rounded-full bg-wood-dark/50 backdrop-blur flex items-center justify-center text-parchment/80 hover:text-parchment cursor-pointer transition-colors text-xs font-bold"
        :aria-label="t('app.switchLanguage')"
        @pointerdown.stop
        @click="cycleLocale"
      >
        {{ nextLocaleLabel }}
      </button>
    </div>

    <main
      v-show="!hideUI"
      class="relative z-10 pb-20"
      style="padding-bottom: calc(5rem + env(safe-area-inset-bottom, 0px))"
    >
      <router-view v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" :key="route.path" />
        </Transition>
      </router-view>
    </main>

    <nav
      v-show="!hideUI"
      :aria-label="t('nav.menu')"
      class="fixed bottom-0 w-full z-50"
      style="padding-bottom: env(safe-area-inset-bottom, 0px)"
    >
      <div class="max-w-md mx-auto">
        <div class="wooden-panel flex items-stretch py-0 px-0" role="menubar">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            :aria-current="isActive(item.path) ? 'page' : undefined"
            role="menuitem"
            class="flex-1 flex flex-col items-center justify-center py-2 transition-transform cursor-pointer"
            :class="
              isActive(item.path)
                ? 'text-parchment scale-110'
                : 'text-parchment/60 hover:text-parchment/80'
            "
          >
            <TwEmoji :emoji="item.icon" size="1.25rem" class="mb-0.5" />
            <span class="text-[10px] tracking-wider uppercase cursor-pointer">{{
              t(item.key)
            }}</span>
            <div
              v-if="isActive(item.path)"
              class="w-1.5 h-1.5 bg-parchment mt-0.5 cursor-pointer"
              style="clip-path: polygon(50% 0%, 100% 100%, 0% 100%)"
            ></div>
          </router-link>
        </div>
      </div>
    </nav>

    <PrivacyNotice v-show="!hideUI" />

    <!-- Achievement toast (global) -->
    <Teleport to="body">
      <Transition name="achievement">
        <div
          v-if="achievementToast"
          :key="achievementToast.id"
          class="fixed bottom-20 left-1/2 -translate-x-1/2 z-[200] w-[90vw] max-w-sm"
        >
          <div class="achievement-toast flex items-center gap-3 px-4 py-3">
            <span class="text-3xl">{{ achievementToast.icon }}</span>
            <div class="flex-1 min-w-0">
              <p
                class="achievement-title text-xs uppercase tracking-widest mb-0.5"
              >
                Achievement Unlocked
              </p>
              <p class="text-parchment text-sm font-bold">
                {{ t(`achievements.${achievementToast.id}.name`) }}
              </p>
              <p class="text-parchment/60 text-xs">
                {{ t(`achievements.${achievementToast.id}.description`) }}
              </p>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Hall of Fame dialog (global) -->
    <Teleport to="body">
      <Transition name="dialog">
        <div
          v-if="showHofDialog"
          role="dialog"
          aria-modal="true"
          aria-labelledby="hof-dialog-title"
          class="fixed inset-0 z-[100] flex items-center justify-center px-4"
        >
          <div class="fixed inset-0 bg-black/50"></div>
          <div class="relative wooden-panel p-3 w-full max-w-xs">
            <div class="wooden-panel-inner px-5 py-5 text-center">
              <TwEmoji emoji="🎉" size="2.5rem" class="mb-2" />
              <p
                id="hof-dialog-title"
                class="text-primary text-lg font-bold mb-1"
              >
                {{ t('app.hofCongrats') }}
              </p>
              <p class="text-secondary text-sm mb-1 break-keep">
                {{ t('app.hofAllAchieved') }}
              </p>
              <p class="text-wood-dark/80 text-xs mb-4 break-keep">
                {{ t('app.hofLeaveName') }}
              </p>
              <input
                ref="hofInputRef"
                v-model="hofName"
                type="text"
                maxlength="30"
                class="pixel-input w-full mb-3"
                :class="{ 'opacity-40 pointer-events-none': hofAnonymous }"
                :disabled="hofAnonymous"
                :placeholder="t('common.namePlaceholder')"
                @keyup.enter="onHofSubmit"
              />
              <label
                class="flex items-center gap-2 mb-4 cursor-pointer text-sm text-wood-dark/80 select-none justify-center"
                @click.prevent="hofAnonymous = !hofAnonymous"
              >
                <span
                  class="pixel-checkbox"
                  :class="{ 'is-checked': hofAnonymous }"
                >
                  <span v-if="hofAnonymous" class="pixel-checkmark"
                    >&#10003;</span
                  >
                </span>
                {{ t('app.hofAnonymous') }}
              </label>
              <button
                class="pixel-btn text-sm px-6 w-full"
                @click="onHofSubmit"
              >
                {{ t('common.confirm') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Hall of Fame confetti -->
    <Teleport to="body">
      <div
        v-if="showHofDialog"
        class="fixed inset-0 pointer-events-none z-[200] overflow-hidden"
      >
        <span
          v-for="i in 40"
          :key="i"
          class="confetti-piece"
          :style="{
            left: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 2}s`,
            animationDuration: `${2.5 + Math.random() * 2}s`,
            fontSize: `${14 + Math.random() * 14}px`,
          }"
          >{{ ['🎊', '🎉', '✨', '⭐', '💫', '🌟'][i % 6] }}</span
        >
      </div>
    </Teleport>

    <!-- Share copy toast -->
    <Teleport to="body">
      <Transition name="achievement">
        <div
          v-if="shareToast"
          role="status"
          aria-live="polite"
          class="fixed top-14 left-1/2 -translate-x-1/2 z-[200] px-4 py-2 rounded bg-wood-dark/90 text-parchment text-sm whitespace-nowrap"
        >
          {{ t('app.linkCopied') }}
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Eye, EyeOff, Sun, Moon, Share2 } from 'lucide-vue-next';
import { useHearts } from './composables/useHearts';
import { useAchievements } from './composables/useAchievements';
import { setHallOfFameSubmitted, getHallOfFameSubmitted } from './lib/store';
import { createHallOfFameEntry } from './services/api';
import {
  WEDDING_DATETIME,
  GROOM_ENG_NAME,
  BRIDE_ENG_NAME,
  VENUE_NAME,
  DEMO_MODE,
} from './config/wedding';
import DemoRibbon from './components/DemoRibbon.vue';
import { useWeddingDate } from './composables/useWeddingDate';
import { useTheme } from './composables/useTheme';
import { useI18n } from 'vue-i18n';
import { messages, setLocale, type Locale } from './i18n';
import TwEmoji from './components/TwEmoji.vue';
import PrivacyNotice from './components/PrivacyNotice.vue';

const { connect, disconnect, onGameReset, offGameReset } = useHearts();
const {
  currentToast: achievementToast,
  award,
  clearAll: clearAchievements,
  allEarned,
  showHofDialog,
} = useAchievements();

const hideUI = ref(false);
const shareToast = ref(false);
let shareToastTimer: ReturnType<typeof setTimeout> | null = null;

async function shareInvitation() {
  const title = t('app.shareTitle', {
    groom: GROOM_ENG_NAME,
    bride: BRIDE_ENG_NAME,
  });
  const text = `${weddingDate.value} ${weddingDay.value} | ${VENUE_NAME}`;
  const url = location.origin;

  if (navigator.share) {
    try {
      await navigator.share({ title, text, url });
      return;
    } catch {
      // User cancelled or share failed — fall through to clipboard
    }
  }

  try {
    await navigator.clipboard.writeText(url);
  } catch {
    // Clipboard API not available
  }
  shareToast.value = true;
  if (shareToastTimer) clearTimeout(shareToastTimer);
  shareToastTimer = setTimeout(() => {
    shareToast.value = false;
  }, 2000);
}
const hofName = ref('');
const hofAnonymous = ref(false);
const hofInputRef = ref<HTMLInputElement | null>(null);

function onHofSubmit() {
  const name = hofAnonymous.value
    ? t('app.anonymousPrefix') + Math.random().toString(36).slice(2, 8)
    : hofName.value;
  onHofConfirm(name);
}

const appRouter = useRouter();
appRouter.afterEach((to) => {
  if (
    to.name === 'achievements' &&
    allEarned.value &&
    !getHallOfFameSubmitted()
  ) {
    showHofDialog.value = true;
    hofName.value = '';
    nextTick(() => hofInputRef.value?.focus());
  }
});

async function onHofConfirm(name: string) {
  if (!hofAnonymous.value && !name.trim()) return;
  try {
    await createHallOfFameEntry(name.trim());
    setHallOfFameSubmitted();
  } catch (e) {
    console.error('Failed to submit hall of fame entry:', e);
  }
  showHofDialog.value = false;
}
const { isDark, toggle: rawToggle, transitioning } = useTheme();

const { t, locale } = useI18n();
const { date: weddingDate, day: weddingDay } = useWeddingDate();
const locales = Object.keys(messages) as Locale[];
const nextLocale = computed(
  () => locales[(locales.indexOf(locale.value as Locale) + 1) % locales.length],
);
const nextLocaleLabel = computed(
  () => messages[nextLocale.value].app.localeLabel,
);
function cycleLocale() {
  setLocale(nextLocale.value);
}

function toggle() {
  rawToggle();
  award('theme');
}
onMounted(() => {
  connect();
  onGameReset(clearAchievements);

  const now = new Date();
  if (
    now.getFullYear() === WEDDING_DATETIME.getFullYear() &&
    now.getMonth() === WEDDING_DATETIME.getMonth() &&
    now.getDate() === WEDDING_DATETIME.getDate()
  ) {
    award('wedding_day');
  }
});
onUnmounted(() => {
  disconnect();
  offGameReset(clearAchievements);
});

const route = useRoute();
const phaserContainer = ref<HTMLDivElement | null>(null);
let phaserGame: any = null;

onMounted(() => {
  if (!phaserContainer.value) return;
  const container = phaserContainer.value;
  const deferInit = async () => {
    const [Phaser, { default: BackgroundScene }] = await Promise.all([
      import('phaser'),
      import('./game/BackgroundScene'),
    ]);
    if (!container.isConnected) return;
    phaserGame = new Phaser.default.Game({
      type: Phaser.default.CANVAS,
      parent: container,
      width: 448,
      height: 960,
      scene: [BackgroundScene],
      transparent: false,
      backgroundColor: isDark.value ? '#0f1f0a' : '#3a7d22',
      banner: false,
      scale: {
        mode: Phaser.default.Scale.RESIZE,
        autoCenter: Phaser.default.Scale.CENTER_HORIZONTALLY,
      },
      audio: { noAudio: true },
      render: { pixelArt: true, antialias: false },
    });
  };
  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(deferInit);
  } else {
    setTimeout(deferInit, 1);
  }
});

onUnmounted(() => {
  phaserGame?.destroy(true);
  phaserGame = null;
});

const navItems = [
  { path: '/', icon: '🏠', key: 'nav.home' },
  { path: '/map', icon: '🗺️', key: 'nav.map' },
  { path: '/guestbook', icon: '📖', key: 'nav.guestbook' },
  { path: '/game', icon: '🎮', key: 'nav.game' },
  { path: '/photo', icon: '📸', key: 'nav.photo' },
];

function isActive(path: string): boolean {
  return route.path === path;
}
</script>

<style>
:deep(canvas) {
  image-rendering: pixelated;
  width: 100% !important;
  height: 100% !important;
}

.page-enter-active,
.page-leave-active {
  transition: opacity 0.2s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}

.achievement-toast {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border: 2px solid #c4943a;
  box-shadow:
    0 0 15px rgba(196, 148, 58, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.achievement-title {
  color: #c4943a;
}

.achievement-enter-active {
  transition:
    transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1),
    opacity 0.3s ease;
}

.achievement-leave-active {
  transition:
    transform 0.4s ease,
    opacity 0.4s ease;
}

.achievement-enter-from,
.achievement-leave-to {
  transform: translate(-50%, 40px);
  opacity: 0;
}

.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.confetti-piece {
  position: absolute;
  top: -30px;
  animation: confetti-fall linear forwards;
}

@keyframes confetti-fall {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    transform: translateY(100vh) rotate(720deg);
    opacity: 0;
  }
}
</style>
