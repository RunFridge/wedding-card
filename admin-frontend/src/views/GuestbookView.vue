<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useGuestbook } from '@/composables/useGuestbook';
import { useSort } from '@/composables/useSort';
import { usePagination } from '@/composables/usePagination';
import { exportCSV } from '@/lib/csv';
import PaginationBar from '@/components/PaginationBar.vue';
import SortableHeader from '@/components/SortableHeader.vue';
import { useBulkSelection } from '@/composables/useBulkSelection';
import { useBulkAction } from '@/composables/useBulkAction';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import ConfirmDialog from '@/components/ConfirmDialog.vue';
import BulkActionBar from '@/components/BulkActionBar.vue';
import {
  MoreHorizontal,
  Eye,
  EyeOff,
  Trash2,
  RefreshCw,
  Download,
} from 'lucide-vue-next';

const { t } = useI18n();
const {
  entries,
  loading,
  load,
  toggleVisibility,
  toggleVisibilityRaw,
  deleteEntry,
  deleteEntryRaw,
} = useGuestbook();
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
  sortKey,
  sortDir,
  toggleSort,
  sortedItems: sortedEntries,
} = useSort(entries);
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

function exportGuestbook() {
  exportCSV(
    'guestbook.csv',
    [
      'ID',
      'Nickname',
      'Message',
      'Secret',
      'Hidden',
      'Evaluated',
      'IP',
      'Date',
    ],
    entries.value.map((e) => [
      String(e.id),
      e.nickname,
      e.message,
      String(e.secret),
      String(e.hidden),
      String(e.evaluated),
      e.ip,
      e.created_at,
    ]),
  );
}

async function bulkHide() {
  const ids = [
    ...entries.value.filter((e) => isSelected(e.id)).map((e) => e.id),
  ];
  await bulk.execute(ids, (id) => toggleVisibilityRaw(id, true), reload);
}

async function bulkShow() {
  const ids = [
    ...entries.value.filter((e) => isSelected(e.id)).map((e) => e.id),
  ];
  await bulk.execute(ids, (id) => toggleVisibilityRaw(id, false), reload);
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

function truncate(text: string, max = 50) {
  return text.length > max ? text.slice(0, max) + '...' : text;
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('nav.guestbook') }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ t('guestbook.subtitle') }}
        </p>
      </div>
      <div class="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="!entries.length"
          @click="exportGuestbook"
        >
          <Download class="mr-2 h-4 w-4" />
          {{ t('common.csv') }}
        </Button>
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
      <Button
        size="sm"
        variant="outline"
        :disabled="bulk.running.value"
        @click="bulkHide"
      >
        <EyeOff class="mr-2 h-4 w-4" />
        {{ t('common.hide') }}
      </Button>
      <Button
        size="sm"
        variant="outline"
        :disabled="bulk.running.value"
        @click="bulkShow"
      >
        <Eye class="mr-2 h-4 w-4" />
        {{ t('common.show') }}
      </Button>
      <ConfirmDialog
        :title="t('common.bulkDelete')"
        :description="t('guestbook.bulkDeleteDesc', { count: selectedCount })"
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
      {{ t('guestbook.empty') }}
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
          <TableHead class="w-16">{{ t('common.id') }}</TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.nickname')"
              field="nickname"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead class="max-w-xs">
            <SortableHeader
              :label="t('guestbook.message')"
              field="message"
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
              :label="t('guestbook.secret')"
              field="secret"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.hidden')"
              field="hidden"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.evaluated')"
              field="evaluated"
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
        <TableRow v-for="entry in pagedEntries" :key="entry.id">
          <TableCell>
            <Checkbox
              :model-value="isSelected(entry.id)"
              :disabled="bulk.running.value"
              @update:model-value="toggleOne(entry.id, $event as boolean)"
            />
          </TableCell>
          <TableCell class="font-mono text-xs">{{ entry.id }}</TableCell>
          <TableCell class="font-medium">{{ entry.nickname }}</TableCell>
          <TableCell class="max-w-xs">
            <TooltipProvider v-if="entry.message.length > 50">
              <Tooltip>
                <TooltipTrigger class="cursor-default text-left">
                  {{ truncate(entry.message) }}
                </TooltipTrigger>
                <TooltipContent side="bottom" class="max-w-sm">
                  {{ entry.message }}
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
            <span v-else>{{ entry.message }}</span>
          </TableCell>
          <TableCell class="font-mono text-xs">{{ entry.ip }}</TableCell>
          <TableCell>
            <Badge
              v-if="entry.secret"
              class="bg-amber-500 text-white hover:bg-amber-600"
            >
              {{ t('common.yes') }}
            </Badge>
            <span v-else class="text-muted-foreground text-xs">{{
              t('common.no')
            }}</span>
          </TableCell>
          <TableCell>
            <Badge :variant="entry.hidden ? 'destructive' : 'secondary'">
              {{ entry.hidden ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell>
            <Badge :variant="entry.evaluated ? 'default' : 'outline'">
              {{ entry.evaluated ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell class="text-xs">{{
            formatDate(entry.created_at)
          }}</TableCell>
          <TableCell>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button
                  variant="ghost"
                  size="icon"
                  class="h-8 w-8"
                  :disabled="bulk.running.value"
                >
                  <MoreHorizontal class="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  @click="toggleVisibility(entry.id, !entry.hidden)"
                >
                  <EyeOff v-if="!entry.hidden" class="mr-2 h-4 w-4" />
                  <Eye v-else class="mr-2 h-4 w-4" />
                  {{ entry.hidden ? t('common.show') : t('common.hide') }}
                </DropdownMenuItem>
                <ConfirmDialog
                  :title="t('guestbook.deleteTitle')"
                  :description="
                    t('guestbook.deleteDesc', { name: entry.nickname })
                  "
                  @confirm="deleteEntry(entry.id)"
                >
                  <DropdownMenuItem
                    class="text-destructive focus:text-destructive"
                    @select.prevent
                  >
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ t('common.delete') }}
                  </DropdownMenuItem>
                </ConfirmDialog>
              </DropdownMenuContent>
            </DropdownMenu>
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
