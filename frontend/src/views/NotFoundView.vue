<template>
  <div
    class="max-w-md mx-auto px-4 py-6 flex flex-col items-center justify-center min-h-[70vh]"
  >
    <div class="wooden-panel p-2 w-full">
      <div class="wooden-panel-inner px-6 py-8 text-center">
        <h1 class="text-5xl text-secondary mb-3">404</h1>
        <p class="text-wood-dark text-base">{{ t('notFound.message') }}</p>
        <div class="stardew-divider my-4"></div>
        <p class="blink text-parchment/70 text-sm tracking-widest mb-6">
          PAGE NOT FOUND
        </p>
        <router-link to="/" class="pixel-btn inline-block">
          {{ t('notFound.goHome') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
.blink {
  animation: blink 1.2s step-end infinite;
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
</style>
