<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { isDemo, loadDemoFlag } from '@/lib/demo';

const { t } = useI18n();

onMounted(loadDemoFlag);
</script>

<template>
  <div
    v-if="isDemo"
    class="fixed bottom-0 inset-x-0 z-50 bg-amber-500 py-1 text-center text-xs font-medium text-amber-950"
  >
    {{ t('demo.banner') }}
  </div>
</template>
