<template>
  <div class="max-w-md mx-auto">
    <!-- Title Screen Section -->
    <section
      class="relative h-screen overflow-hidden flex flex-col items-center px-4 pt-[10vh] sm:pt-[15vh]"
    >
      <!-- Title Sign -->
      <div class="wooden-panel p-3 w-full max-w-sm">
        <div class="wooden-panel-inner px-6 py-8 text-center">
          <p class="text-wood-dark/80 text-sm tracking-[0.3em] mb-3">
            Wedding Invitation
          </p>
          <div class="stardew-divider my-3"></div>
          <h1 class="text-3xl text-secondary leading-snug">
            {{ GROOM_ENG_NAME }}
            <span class="text-wood-dark/60 text-lg mx-1">&amp;</span>
            {{ BRIDE_ENG_NAME }}
          </h1>
          <div class="stardew-divider my-3"></div>
          <p class="text-wood-dark text-base mt-2">{{ weddingDate }}</p>
          <p class="text-wood-dark/80 text-sm mt-1">
            {{ weddingDay }} {{ weddingTime }}
          </p>
          <p class="text-wood-dark/80 text-sm mt-1">{{ VENUE_NAME }}</p>
        </div>
      </div>

      <!-- Press Start Button -->
      <button
        class="press-start mt-12 text-parchment text-lg tracking-widest cursor-pointer bg-transparent border-none"
        @click="scrollToContent"
        @keydown.enter="scrollToContent"
      >
        &#9654; PRESS START
      </button>
    </section>

    <!-- Invitation Content Section -->
    <section ref="contentSection" class="relative px-4 py-8 text-center">
      <div class="panel-enter wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-6 py-6">
          <h2 class="text-2xl text-secondary mb-3">
            {{ GROOM_KOR_NAME }} & {{ BRIDE_KOR_NAME }}
          </h2>
          <p ref="shortGreetingRef" class="text-wood-dark">
            {{ shortGreetingDisplay }}
          </p>
        </div>
      </div>

      <div class="panel-enter wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-6 py-4">
          <p class="text-lg text-secondary mb-1">{{ weddingDate }}</p>
          <p class="text-base text-secondary">
            {{ weddingDay }} {{ weddingTime }}
          </p>
          <p v-if="!isKST" class="text-wood-dark/60 text-sm mt-1">
            {{ localTime }}
          </p>
          <div class="stardew-divider my-3"></div>
          <p class="text-secondary text-base font-semibold">
            {{ VENUE_NAME }}<br />{{ VENUE_FLOOR }} {{ VENUE_HALL }}
          </p>
          <p class="text-wood-dark/80 text-sm mt-1">{{ VENUE_ADDRESS }}</p>
        </div>
      </div>

      <div class="panel-enter parchment-bg p-4 mb-6">
        <div class="rounded overflow-hidden border-2 border-wood-dark">
          <img
            :src="mainPhoto"
            :alt="`${GROOM_ENG_NAME} & ${BRIDE_ENG_NAME}`"
            class="w-full"
            width="340"
            height="510"
          />
        </div>
      </div>

      <div class="panel-enter wooden-panel p-2 mb-6">
        <div
          ref="mainGreetRef"
          class="wooden-panel-inner px-6 py-5 text-wood-dark leading-relaxed text-base whitespace-pre-line break-keep"
        >
          {{ mainGreetDisplay }}
        </div>
      </div>

      <div class="panel-enter wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-6 py-4 text-center">
          <p class="text-wood-dark/80 text-sm mb-2">
            {{ t('home.untilWedding') }}
          </p>
          <template v-if="countdown">
            <p class="text-2xl text-primary font-semibold">
              D-{{ countdown.days }}
            </p>
            <p class="text-lg text-primary font-semibold mt-1">
              {{ countdown.time }}
            </p>
          </template>
          <p v-else class="text-lg text-secondary">{{ t('home.thankYou') }}</p>
        </div>
      </div>

      <div class="panel-enter wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-4 py-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-wood-dark/80 mb-1">
                <TwEmoji emoji="🤵" size="0.875rem" />
                {{ t('home.groomParents') }}
              </p>
              <p class="text-secondary text-base">
                {{ GROOM_FATHER_KOR_NAME }} · {{ GROOM_MOTHER_KOR_NAME }}
              </p>
              <p class="text-secondary font-semibold mt-1">
                {{ GROOM_BIRTH_ORDER }} {{ GROOM_KOR_NAME }}
              </p>
            </div>
            <div>
              <p class="text-sm text-wood-dark/80 mb-1">
                <TwEmoji emoji="👰" size="0.875rem" />
                {{ t('home.brideParents') }}
              </p>
              <p class="text-secondary text-base">
                {{ BRIDE_FATHER_KOR_NAME }} · {{ BRIDE_MOTHER_KOR_NAME }}
              </p>
              <p class="text-secondary font-semibold mt-1">
                {{ BRIDE_BIRTH_ORDER }} {{ BRIDE_KOR_NAME }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="panel-enter wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-4 py-4">
          <button
            @click="accountsOpen = !accountsOpen"
            class="w-full flex justify-between items-center bg-transparent border-none p-0"
          >
            <h2 class="text-secondary text-sm">{{ t('home.gift') }}</h2>
            <span
              class="text-secondary text-sm transition-transform"
              :class="accountsOpen ? 'rotate-180' : ''"
              >▼</span
            >
          </button>

          <div v-if="accountsOpen" class="mt-3 select-none">
            <div class="stardew-divider mb-4"></div>

            <div class="mb-4">
              <p class="text-sm text-wood-dark/80 mb-2">
                <TwEmoji emoji="🤵" size="0.875rem" /> {{ t('home.groomSide') }}
              </p>
              <div class="space-y-2 text-left">
                <button
                  v-for="a in groomAccounts"
                  :key="a.name"
                  :aria-label="t('home.copyAccountAria', { name: a.name })"
                  @click="copyAccount(a.account)"
                  class="flex flex-wrap justify-between items-center w-full text-left cursor-pointer hover:bg-wood-dark/5 active:bg-wood-dark/10 rounded px-1 py-1 -mx-1 transition-colors bg-transparent border-none"
                >
                  <span class="text-secondary">{{ a.name }}</span>
                  <span class="text-wood-dark text-sm flex items-center gap-1"
                    >{{ a.account }}
                    <span class="text-xs text-wood-dark/40">{{
                      t('home.copy')
                    }}</span></span
                  >
                </button>
              </div>
            </div>

            <div class="stardew-divider mb-4"></div>

            <div>
              <p class="text-sm text-wood-dark/80 mb-2">
                <TwEmoji emoji="👰" size="0.875rem" /> {{ t('home.brideSide') }}
              </p>
              <div class="space-y-2 text-left">
                <button
                  v-for="a in brideAccounts"
                  :key="a.name"
                  :aria-label="t('home.copyAccountAria', { name: a.name })"
                  @click="copyAccount(a.account)"
                  class="flex flex-wrap justify-between items-center w-full text-left cursor-pointer hover:bg-wood-dark/5 active:bg-wood-dark/10 rounded px-1 py-1 -mx-1 transition-colors bg-transparent border-none"
                >
                  <span class="text-secondary">{{ a.name }}</span>
                  <span class="text-wood-dark text-sm flex items-center gap-1"
                    >{{ a.account }}
                    <span class="text-xs text-wood-dark/40">{{
                      t('home.copy')
                    }}</span></span
                  >
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="panel-enter flex flex-col gap-3">
        <router-link to="/map" class="pixel-btn block text-center">
          {{ t('map.title') }}
        </router-link>
      </div>
    </section>

    <Teleport to="body">
      <Transition name="toast">
        <div
          v-if="copyToast"
          class="fixed bottom-8 left-1/2 -translate-x-1/2 z-[100] parchment-bg px-6 py-2 text-sm text-secondary border-2 border-wood-dark shadow-lg w-[80vw] max-w-sm text-center"
        >
          {{ t('home.accountCopied') }}
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import TwEmoji from '../components/TwEmoji.vue';
import fallbackPhoto from '../assets/fallback_photo.svg';
import {
  MAIN_PHOTO_URL,
  GROOM_ENG_NAME,
  GROOM_KOR_NAME,
  BRIDE_ENG_NAME,
  BRIDE_KOR_NAME,
  GROOM_FATHER_KOR_NAME,
  GROOM_MOTHER_KOR_NAME,
  BRIDE_FATHER_KOR_NAME,
  BRIDE_MOTHER_KOR_NAME,
  GROOM_BANK_ACCOUNT,
  BRIDE_BANK_ACCOUNT,
  GROOM_FATHER_BANK_ACCOUNT,
  GROOM_MOTHER_BANK_ACCOUNT,
  BRIDE_FATHER_BANK_ACCOUNT,
  BRIDE_MOTHER_BANK_ACCOUNT,
  GROOM_BIRTH_ORDER,
  BRIDE_BIRTH_ORDER,
  WEDDING_DATETIME,
  VENUE_NAME,
  VENUE_ADDRESS,
  VENUE_FLOOR,
  VENUE_HALL,
  SHORT_GREETING,
  MAIN_GREET_TEXT,
} from '../config/wedding';
import { useI18n } from 'vue-i18n';
import { useWeddingDate } from '../composables/useWeddingDate';

const { t, locale } = useI18n();
const {
  date: weddingDate,
  day: weddingDay,
  time: weddingTime,
} = useWeddingDate();

interface AccountInfo {
  name: string;
  account: string;
}

const mainPhoto = computed(() => MAIN_PHOTO_URL || fallbackPhoto);

const contentSection = ref<HTMLElement | null>(null);
const copyToast = ref(false);
const accountsOpen = ref(false);
const now = ref(new Date());

const shortGreetingRef = ref<HTMLElement | null>(null);
const mainGreetRef = ref<HTMLElement | null>(null);
const shortGreetingDisplay = ref('');
const mainGreetDisplay = ref('');
const typingTimers: ReturnType<typeof setInterval>[] = [];
let typingObserver: IntersectionObserver | null = null;

function typeText(text: string, display: { value: string }, speed: number) {
  let i = 0;
  const timer = setInterval(() => {
    if (i < text.length) {
      display.value = text.slice(0, ++i);
    } else {
      clearInterval(timer);
    }
  }, speed);
  typingTimers.push(timer);
}

let ddayTimer: ReturnType<typeof setInterval> | null = null;
onMounted(() => {
  ddayTimer = setInterval(() => {
    now.value = new Date();
  }, 1000);

  typingObserver = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        if (entry.target === shortGreetingRef.value) {
          typeText(SHORT_GREETING, shortGreetingDisplay, 60);
        } else if (entry.target === mainGreetRef.value) {
          typeText(MAIN_GREET_TEXT, mainGreetDisplay, 60);
        }
        typingObserver!.unobserve(entry.target);
      }
    },
    { threshold: 0.1 },
  );

  if (shortGreetingRef.value) typingObserver.observe(shortGreetingRef.value);
  if (mainGreetRef.value) typingObserver.observe(mainGreetRef.value);
});
onUnmounted(() => {
  if (ddayTimer) clearInterval(ddayTimer);
  typingTimers.forEach(clearInterval);
  typingObserver?.disconnect();
});

const countdown = computed(() => {
  const diff = WEDDING_DATETIME.getTime() - now.value.getTime();
  if (diff <= 0) return null;
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((diff % (1000 * 60)) / 1000);
  const hh = String(hours).padStart(2, '0');
  const mm = String(minutes).padStart(2, '0');
  const ss = String(seconds).padStart(2, '0');
  return { days, time: `${hh}:${mm}:${ss}` };
});

const isKST = computed<boolean>(() => {
  const kstOffset = WEDDING_DATETIME.toLocaleString('en-US', {
    timeZone: 'Asia/Seoul',
  });
  const localOffset = WEDDING_DATETIME.toLocaleString('en-US');
  return kstOffset === localOffset;
});

const localTime = computed<string>(() => {
  return WEDDING_DATETIME.toLocaleString(locale.value, {
    month: 'long',
    day: 'numeric',
    weekday: 'short',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
    timeZoneName: 'short',
  });
});

const groomAccounts: AccountInfo[] = [
  { name: GROOM_KOR_NAME, account: GROOM_BANK_ACCOUNT },
  { name: GROOM_FATHER_KOR_NAME, account: GROOM_FATHER_BANK_ACCOUNT },
  { name: GROOM_MOTHER_KOR_NAME, account: GROOM_MOTHER_BANK_ACCOUNT },
];

const brideAccounts: AccountInfo[] = [
  { name: BRIDE_KOR_NAME, account: BRIDE_BANK_ACCOUNT },
  { name: BRIDE_FATHER_KOR_NAME, account: BRIDE_FATHER_BANK_ACCOUNT },
  { name: BRIDE_MOTHER_KOR_NAME, account: BRIDE_MOTHER_BANK_ACCOUNT },
];

async function copyAccount(account: string) {
  try {
    await navigator.clipboard.writeText(account);
    copyToast.value = true;
    setTimeout(() => {
      copyToast.value = false;
    }, 1500);
  } catch {
    const textarea = document.createElement('textarea');
    textarea.value = account;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    copyToast.value = true;
    setTimeout(() => {
      copyToast.value = false;
    }, 1500);
  }
}

function scrollToContent() {
  contentSection.value?.scrollIntoView({ behavior: 'smooth' });
}

let panelObserver: IntersectionObserver | null = null;

onMounted(() => {
  panelObserver = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add('panel-visible');
          panelObserver?.unobserve(entry.target);
        }
      }
    },
    { threshold: 0.15 },
  );
  document.querySelectorAll('.panel-enter').forEach((el) => {
    panelObserver!.observe(el);
  });
});

onUnmounted(() => {
  panelObserver?.disconnect();
});
</script>

<style scoped>
.press-start {
  animation: blink 1.2s step-end infinite;
  text-shadow:
    2px 2px 0 rgba(0, 0, 0, 0.5),
    -1px -1px 0 rgba(0, 0, 0, 0.3);
}

@keyframes blink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 10px);
}

.panel-enter {
  opacity: 0;
  transform: translateY(24px);
  transition:
    opacity 0.6s ease,
    transform 0.6s ease;
}

.panel-enter.panel-visible {
  opacity: 1;
  transform: translateY(0);
}

.panel-enter:nth-child(2) {
  transition-delay: 0.08s;
}
.panel-enter:nth-child(3) {
  transition-delay: 0.16s;
}
.panel-enter:nth-child(4) {
  transition-delay: 0.24s;
}
</style>
