<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-4">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">
          {{ t('achievementsView.title') }}
        </h1>
        <p class="text-wood-dark text-sm">
          {{
            t('achievementsView.progress', {
              earned: earnedCount,
              total: ACHIEVEMENTS.length,
            })
          }}
        </p>
        <p class="text-wood-dark/70 text-[10px] break-keep mt-1">
          {{ t('achievementsView.medalNote') }}
        </p>
      </div>
    </div>

    <div class="wooden-panel p-2">
      <div class="wooden-panel-inner p-4 space-y-2">
        <div
          v-for="a in ACHIEVEMENTS"
          :key="a.id"
          class="achievement-row flex items-center gap-3 px-3 py-2.5 parchment-bg"
          :class="[
            earned.has(a.id) ? '' : 'opacity-40',
            rarityClass(a.id, earned.has(a.id)),
          ]"
        >
          <TwEmoji
            :emoji="earned.has(a.id) ? a.icon : '🔒'"
            size="1.75rem"
            class="shrink-0"
          />
          <div class="flex-1 min-w-0">
            <p class="text-secondary text-sm font-bold">
              {{
                earned.has(a.id) || !a.secret
                  ? t(`achievements.${a.id}.name`)
                  : '???'
              }}
            </p>
            <p class="text-wood-dark text-xs">
              {{
                earned.has(a.id) || !a.secret
                  ? t(`achievements.${a.id}.description`)
                  : '???'
              }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <p class="text-center text-parchment/50 text-xs mt-3 break-keep">
      {{ t('achievementsView.localOnly') }}
    </p>

    <div class="flex justify-center gap-2 mt-3">
      <button class="pixel-btn text-sm px-6 py-2" @click="router.back()">
        {{ t('common.back') }}
      </button>
      <router-link
        to="/hall-of-fame"
        class="pixel-btn text-sm px-6 py-2 no-underline"
      >
        {{ t('hallOfFame.title') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useAchievements } from '../composables/useAchievements';
import TwEmoji from '../components/TwEmoji.vue';

const { t } = useI18n();
const router = useRouter();

const { earned, ACHIEVEMENTS } = useAchievements();
const earnedCount = computed(() => earned.value.size);

function rarityClass(id: string, isEarned: boolean): string {
  if (!isEarned) return '';
  if (id === 'gold') return 'rarity-gold';
  if (id === 'silver') return 'rarity-silver';
  if (id === 'bronze') return 'rarity-bronze';
  return '';
}
</script>

<style scoped>
.achievement-row {
  position: relative;
  overflow: hidden;
  isolation: isolate;
}

.achievement-row > * {
  position: relative;
  z-index: 2;
}

/* ─── Bronze: warm metallic glow ─── */
.parchment-bg.rarity-bronze {
  background: linear-gradient(
    135deg,
    rgba(205, 127, 50, 0.18),
    rgba(140, 75, 25, 0.06)
  );
  border-color: #cd7f32;
  box-shadow:
    0 0 6px rgba(205, 127, 50, 0.35),
    inset 0 0 10px rgba(205, 127, 50, 0.1);
}

/* ─── Silver: restrained iridescence (cyan/purple) ─── */
.parchment-bg.rarity-silver {
  background: linear-gradient(
    135deg,
    rgba(220, 220, 230, 0.18),
    rgba(160, 160, 180, 0.05)
  );
  border-color: #c0c0c8;
  box-shadow:
    0 0 8px rgba(220, 220, 230, 0.4),
    inset 0 0 10px rgba(220, 220, 230, 0.12);
}

.parchment-bg.rarity-silver::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    110deg,
    hsla(280, 50%, 70%, 0.18) 0%,
    hsla(200, 60%, 70%, 0.18) 50%,
    hsla(280, 50%, 70%, 0.18) 100%
  );
  background-size: 250% 100%;
  mix-blend-mode: overlay;
  animation: rarity-iridescent-shift 8s linear infinite;
  pointer-events: none;
  z-index: 0;
}

/* ─── Gold: rainbow holographic ─── */
.parchment-bg.rarity-gold {
  background: linear-gradient(
    135deg,
    rgba(255, 215, 0, 0.22),
    rgba(255, 165, 0, 0.08) 50%,
    rgba(255, 215, 0, 0.18)
  );
  border-color: #ffd700;
  box-shadow:
    0 0 10px rgba(255, 215, 0, 0.5),
    inset 0 0 12px rgba(255, 215, 0, 0.15);
  animation: rarity-gold-pulse 4s ease-in-out infinite;
  text-shadow:
    0 1px 2px rgba(0, 0, 0, 0.45),
    0 0 1px rgba(0, 0, 0, 0.35);
}

.parchment-bg.rarity-gold::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    110deg,
    hsl(330, 70%, 65%) 0%,
    hsl(30, 75%, 65%) 14%,
    hsl(60, 75%, 62%) 28%,
    hsl(130, 55%, 58%) 42%,
    hsl(190, 60%, 62%) 56%,
    hsl(250, 60%, 68%) 70%,
    hsl(310, 65%, 65%) 84%,
    hsl(330, 70%, 65%) 100%
  );
  background-size: 280% 100%;
  mix-blend-mode: overlay;
  filter: saturate(0.7) brightness(0.9);
  opacity: 0.45;
  animation: rarity-holo-shift 7s linear infinite;
  pointer-events: none;
  z-index: 0;
}

.parchment-bg.rarity-gold::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    115deg,
    transparent 35%,
    rgba(255, 255, 255, 0.25) 50%,
    transparent 65%
  );
  transform: translateX(-100%);
  animation: rarity-shimmer 4s ease-in-out infinite;
  pointer-events: none;
  z-index: 1;
}

@keyframes rarity-shimmer {
  0% {
    transform: translateX(-100%);
  }
  60%,
  100% {
    transform: translateX(100%);
  }
}

@keyframes rarity-iridescent-shift {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 250% 50%;
  }
}

@keyframes rarity-holo-shift {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 280% 50%;
  }
}

@keyframes rarity-gold-pulse {
  0%,
  100% {
    box-shadow:
      0 0 10px rgba(255, 215, 0, 0.5),
      inset 0 0 12px rgba(255, 215, 0, 0.15);
  }
  50% {
    box-shadow:
      0 0 14px rgba(255, 215, 0, 0.7),
      inset 0 0 14px rgba(255, 215, 0, 0.22);
  }
}

@media (prefers-reduced-motion: reduce) {
  .parchment-bg.rarity-silver::before,
  .parchment-bg.rarity-gold,
  .parchment-bg.rarity-gold::before,
  .parchment-bg.rarity-gold::after {
    animation: none;
  }
}
</style>
