<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-4">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">{{ t('hallOfFame.title') }}</h1>
        <p class="text-wood-dark text-xs break-keep">
          {{ t('hallOfFame.subtitle') }}
        </p>
        <p class="text-wood-dark/70 text-[10px] break-keep mt-1">
          {{ t('achievementsView.medalNote') }}
        </p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-6">
      <div class="pixel-spinner"></div>
    </div>

    <div
      v-else-if="entries.length === 0"
      class="parchment-bg p-6 text-center text-wood-dark/80"
    >
      {{ t('hallOfFame.empty') }}
    </div>

    <div v-else class="wooden-panel p-2">
      <div class="wooden-panel-inner p-4 space-y-2">
        <div
          v-for="(entry, index) in entries"
          :key="entry.id"
          class="flex items-center gap-3 px-3 py-2.5 parchment-bg"
        >
          <span class="text-primary text-sm font-bold w-6"
            >{{ index + 1 }}.</span
          >
          <BoringAvatar
            :name="entry.nickname"
            :size="28"
            class="rounded-full border-2 border-wood-dark/20"
          />
          <span class="text-secondary text-sm font-semibold flex-1">{{
            entry.nickname
          }}</span>
          <span class="text-xs text-wood-dark/50">{{
            formatDate(entry.created_at)
          }}</span>
        </div>
      </div>
    </div>

    <div class="flex justify-center mt-4">
      <button class="pixel-btn text-sm px-6 py-2" @click="router.back()">
        {{ t('common.back') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { getHallOfFame } from '../services/api';
import BoringAvatar from '../components/BoringAvatar.vue';
import type { HallOfFameEntry } from '../types';

const { t, locale } = useI18n();
const router = useRouter();
const entries = ref<HallOfFameEntry[]>([]);
const loading = ref(true);

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

onMounted(async () => {
  try {
    entries.value = await getHallOfFame();
  } catch (e) {
    console.error('Failed to load hall of fame:', e);
  } finally {
    loading.value = false;
  }
});
</script>
