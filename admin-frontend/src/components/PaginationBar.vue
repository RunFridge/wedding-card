<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { Button } from '@/components/ui/button';
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';

const { t } = useI18n();

defineProps<{
  page: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
}>();

const emit = defineEmits<{
  'update:page': [value: number];
  'update:pageSize': [value: number];
}>();

const sizes = [10, 25, 50];
</script>

<template>
  <div
    class="flex flex-col sm:flex-row items-center justify-between gap-2 pt-4 text-sm text-muted-foreground"
  >
    <div class="flex items-center gap-2">
      <span>{{ t('pagination.rows') }}</span>
      <select
        :value="pageSize"
        class="h-8 rounded-md border bg-background px-2 text-sm"
        @change="
          emit(
            'update:pageSize',
            Number(($event.target as HTMLSelectElement).value),
          )
        "
      >
        <option v-for="s in sizes" :key="s" :value="s">{{ s }}</option>
      </select>
      <span class="hidden sm:inline">{{
        t('pagination.ofItems', { count: totalItems })
      }}</span>
    </div>
    <div class="flex items-center gap-1">
      <Button
        variant="outline"
        size="icon"
        class="h-8 w-8"
        :disabled="page <= 1"
        @click="emit('update:page', page - 1)"
      >
        <ChevronLeft class="h-4 w-4" />
      </Button>
      <span class="px-2">{{ page }} / {{ totalPages }}</span>
      <Button
        variant="outline"
        size="icon"
        class="h-8 w-8"
        :disabled="page >= totalPages"
        @click="emit('update:page', page + 1)"
      >
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>
