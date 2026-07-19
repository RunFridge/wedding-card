<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useHallOfFame } from '@/composables/useHallOfFame';
import { useSort } from '@/composables/useSort';
import { usePagination } from '@/composables/usePagination';
import PaginationBar from '@/components/PaginationBar.vue';
import SortableHeader from '@/components/SortableHeader.vue';
import { useBulkSelection } from '@/composables/useBulkSelection';
import { useBulkAction } from '@/composables/useBulkAction';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import ConfirmDialog from '@/components/ConfirmDialog.vue';
import BulkActionBar from '@/components/BulkActionBar.vue';
import { Trash2, RefreshCw } from 'lucide-vue-next';

const { t } = useI18n();
const { entries, loading, load, deleteEntry, deleteEntryRaw } = useHallOfFame();
const {
  sortKey,
  sortDir,
  toggleSort,
  sortedItems: sortedEntries,
} = useSort(entries);
const {
  selectedCount,
  headerChecked,
  isSelected,
  toggleOne,
  toggleAll,
  clearSelection,
} = useBulkSelection(entries);
const bulk = useBulkAction();
const {
  page,
  pageSize,
  totalPages,
  paginatedItems: pagedEntries,
  setPageSize,
  resetPage,
} = usePagination(sortedEntries);

async function reload() {
  clearSelection();
  resetPage();
  await load();
}

async function bulkDelete() {
  const ids = [
    ...entries.value.filter((e) => isSelected(e.id)).map((e) => e.id),
  ];
  await bulk.execute(ids, (id) => deleteEntryRaw(id), reload);
}

onMounted(load);

function formatDate(str: string) {
  return new Date(str).toLocaleString();
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('nav.hallOfFame') }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ t('hallOfFame.subtitle') }}
        </p>
      </div>
      <div class="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="loading || bulk.running.value"
          @click="reload"
        >
          <RefreshCw
            class="mr-2 h-4 w-4"
            :class="{ 'animate-spin': loading }"
          />
          {{ t('common.reload') }}
        </Button>
      </div>
    </div>

    <BulkActionBar
      v-if="selectedCount > 0"
      :selected-count="selectedCount"
      :progress="bulk.progress.value"
      :running="bulk.running.value"
      @clear="clearSelection"
    >
      <ConfirmDialog
        :title="t('common.bulkDelete')"
        :description="t('hallOfFame.bulkDeleteDesc', { count: selectedCount })"
        @confirm="bulkDelete"
      >
        <Button size="sm" variant="destructive" :disabled="bulk.running.value">
          <Trash2 class="mr-2 h-4 w-4" />
          {{ t('common.delete') }}
        </Button>
      </ConfirmDialog>
    </BulkActionBar>

    <div v-if="loading && entries.length === 0" class="space-y-2">
      <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
    </div>

    <div
      v-else-if="entries.length === 0"
      class="py-8 text-center text-muted-foreground"
    >
      {{ t('hallOfFame.empty') }}
    </div>

    <Table v-else>
      <TableHeader>
        <TableRow>
          <TableHead class="w-10">
            <Checkbox
              :model-value="headerChecked"
              :disabled="bulk.running.value"
              @update:model-value="toggleAll($event as boolean)"
            />
          </TableHead>
          <TableHead class="w-16">#</TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.nickname')"
              field="nickname"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.ip')"
              field="ip"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.date')"
              field="created_at"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead class="w-16">{{ t('common.actions') }}</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="(entry, index) in pagedEntries" :key="entry.id">
          <TableCell>
            <Checkbox
              :model-value="isSelected(entry.id)"
              :disabled="bulk.running.value"
              @update:model-value="toggleOne(entry.id, $event as boolean)"
            />
          </TableCell>
          <TableCell class="font-mono">{{
            (page - 1) * pageSize + index + 1
          }}</TableCell>
          <TableCell class="font-medium">{{ entry.nickname }}</TableCell>
          <TableCell class="font-mono text-xs">{{ entry.ip }}</TableCell>
          <TableCell class="text-xs">{{
            formatDate(entry.created_at)
          }}</TableCell>
          <TableCell>
            <ConfirmDialog
              :title="t('hallOfFame.deleteTitle')"
              :description="
                t('hallOfFame.deleteDesc', { name: entry.nickname })
              "
              @confirm="deleteEntry(entry.id)"
            >
              <Button
                variant="ghost"
                size="icon"
                class="h-8 w-8 text-destructive hover:text-destructive"
                :disabled="bulk.running.value"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </ConfirmDialog>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
    <PaginationBar
      :page="page"
      :total-pages="totalPages"
      :page-size="pageSize"
      :total-items="entries.length"
      @update:page="page = $event"
      @update:page-size="setPageSize($event)"
    />
  </div>
</template>
