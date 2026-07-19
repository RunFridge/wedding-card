<script setup lang="ts">
import type { SortDir } from '@/composables/useSort';

defineProps<{
  label: string;
  field: string;
  activeField: string | null;
  direction: SortDir;
}>();

defineEmits<{
  sort: [field: string];
}>();
</script>

<template>
  <button
    class="inline-flex items-center gap-1 bg-transparent border-none p-0 font-medium text-left cursor-pointer hover:text-foreground transition-colors"
    @click="$emit('sort', field)"
  >
    {{ label }}
    <span class="text-xs opacity-50" :class="{ 'opacity-100': activeField === field }">
      <template v-if="activeField === field">{{ direction === 'asc' ? '▲' : '▼' }}</template>
      <template v-else>⇅</template>
    </span>
  </button>
</template>
