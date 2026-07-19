<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRankings } from '@/composables/useRankings';
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
import { Input } from '@/components/ui/input';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import api from '@/lib/axios';

const { t } = useI18n();
const { scores, loading, load, deleteScore, deleteScoreRaw, purgeAll } =
  useRankings();
const {
  selectedCount,
  headerChecked,
  isSelected,
  toggleOne,
  toggleAll,
  clearSelection,
} = useBulkSelection(scores);
const bulk = useBulkAction();
const {
  sortKey,
  sortDir,
  toggleSort,
  sortedItems: sortedScores,
} = useSort(scores);
const {
  page,
  pageSize,
  totalPages,
  paginatedItems: pagedScores,
  setPageSize,
  resetPage,
} = usePagination(sortedScores);

const purging = ref(false);
const purgePassword = ref('');
const purgeError = ref('');

async function reload() {
  clearSelection();
  resetPage();
  await load();
}

async function handlePurge() {
  if (!purgePassword.value) {
    purgeError.value = t('rankings.passwordRequired');
    return;
  }
  purging.value = true;
  purgeError.value = '';
  try {
    const res = await api.post('/admin/verify', {
      password: purgePassword.value,
    });
    if (res.data.status !== 'ok') {
      purgeError.value = t('rankings.wrongPassword');
      return;
    }
    await purgeAll();
    clearSelection();
    purgePassword.value = '';
  } catch {
    purgeError.value = t('rankings.wrongPassword');
  } finally {
    purging.value = false;
  }
}

async function bulkDelete() {
  const ids = [
    ...scores.value.filter((s) => isSelected(s.id)).map((s) => s.id),
  ];
  await bulk.execute(ids, (id) => deleteScoreRaw(id), reload);
}

onMounted(load);

function formatDate(str: string) {
  return new Date(str).toLocaleString();
}

function formatTime(ms: number) {
  return (ms / 1000).toFixed(2) + 's';
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('nav.rankings') }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ t('rankings.subtitle') }}
        </p>
      </div>
      <div class="flex gap-2">
        <AlertDialog
          @update:open="
            purgePassword = '';
            purgeError = '';
          "
        >
          <AlertDialogTrigger as-child>
            <Button
              variant="destructive"
              size="sm"
              :disabled="loading || bulk.running.value || purging"
            >
              <Trash2 class="mr-2 h-4 w-4" />
              {{ purging ? t('rankings.purging') : t('rankings.purgeAll') }}
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{{
                t('rankings.purgeTitle')
              }}</AlertDialogTitle>
              <AlertDialogDescription>
                {{ t('rankings.purgeDesc') }}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <Input
              v-model="purgePassword"
              type="password"
              :placeholder="t('rankings.adminPasswordPlaceholder')"
              @keyup.enter="handlePurge"
            />
            <p v-if="purgeError" class="text-sm text-destructive">
              {{ purgeError }}
            </p>
            <AlertDialogFooter>
              <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
              <AlertDialogAction
                :disabled="!purgePassword || purging"
                @click.prevent="handlePurge"
              >
                {{
                  purging ? t('rankings.verifying') : t('rankings.confirmPurge')
                }}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
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
        :description="t('rankings.bulkDeleteDesc', { count: selectedCount })"
        @confirm="bulkDelete"
      >
        <Button size="sm" variant="destructive" :disabled="bulk.running.value">
          <Trash2 class="mr-2 h-4 w-4" />
          {{ t('common.delete') }}
        </Button>
      </ConfirmDialog>
    </BulkActionBar>

    <div v-if="loading && scores.length === 0" class="space-y-2">
      <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
    </div>

    <div
      v-else-if="scores.length === 0"
      class="py-8 text-center text-muted-foreground"
    >
      {{ t('rankings.empty') }}
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
          <TableHead class="w-16">{{ t('rankings.rank') }}</TableHead>
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
              :label="t('rankings.time')"
              field="time_ms"
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
        <TableRow v-for="score in pagedScores" :key="score.id">
          <TableCell>
            <Checkbox
              :model-value="isSelected(score.id)"
              :disabled="bulk.running.value"
              @update:model-value="toggleOne(score.id, $event as boolean)"
            />
          </TableCell>
          <TableCell class="font-mono">{{
            scores.indexOf(score) + 1
          }}</TableCell>
          <TableCell class="font-medium">{{ score.nickname }}</TableCell>
          <TableCell class="font-mono">{{
            formatTime(score.time_ms)
          }}</TableCell>
          <TableCell class="font-mono text-xs">{{ score.ip }}</TableCell>
          <TableCell class="text-xs">{{
            formatDate(score.created_at)
          }}</TableCell>
          <TableCell>
            <ConfirmDialog
              :title="t('rankings.deleteTitle')"
              :description="t('rankings.deleteDesc', { name: score.nickname })"
              @confirm="deleteScore(score.id)"
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
      :total-items="scores.length"
      @update:page="page = $event"
      @update:page-size="setPageSize($event)"
    />
  </div>
</template>
