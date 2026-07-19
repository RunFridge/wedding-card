<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import type { BulkProgress } from '@/composables/useBulkAction';
import { Button } from '@/components/ui/button';
import { X } from 'lucide-vue-next';

const { t } = useI18n();

defineProps<{
  selectedCount: number;
  progress: BulkProgress | null;
  running: boolean;
}>();

defineEmits<{
  clear: [];
}>();
</script>

<template>
  <div
    class="mb-4 flex items-center gap-3 rounded-lg border bg-muted/50 px-4 py-2"
  >
    <span class="text-sm font-medium">
      {{ t('common.selected', { count: selectedCount }) }}
    </span>

    <div class="flex items-center gap-2">
      <slot />
    </div>

    <Button
      variant="ghost"
      size="icon"
      class="ml-auto h-7 w-7"
      @click="$emit('clear')"
    >
      <X class="h-4 w-4" />
    </Button>

    <div v-if="running && progress" class="absolute bottom-0 left-0 right-0">
      <div class="h-1 w-full overflow-hidden bg-muted">
        <div
          class="h-full bg-primary transition-all duration-300"
          :style="{ width: `${(progress.completed / progress.total) * 100}%` }"
        />
      </div>
    </div>
  </div>
</template>
