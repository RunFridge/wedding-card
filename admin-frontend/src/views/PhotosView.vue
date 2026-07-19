<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { usePhotos } from '@/composables/usePhotos';
import { useSort } from '@/composables/useSort';
import { usePagination } from '@/composables/usePagination';
import { exportCSV } from '@/lib/csv';
import PaginationBar from '@/components/PaginationBar.vue';
import SortableHeader from '@/components/SortableHeader.vue';
import { useConfig } from '@/composables/useConfig';
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
  ExternalLink,
  RefreshCw,
  Download,
  Upload,
  UploadCloud,
  HeartOff,
} from 'lucide-vue-next';

const { t } = useI18n();
const {
  photos,
  loading,
  load,
  toggleVisibility,
  toggleVisibilityRaw,
  deletePhoto,
  deletePhotoRaw,
  resetHearts,
} = usePhotos();
const { config, load: loadConfig, save: saveConfig } = useConfig();
const {
  sortKey,
  sortDir,
  toggleSort,
  sortedItems: sortedPhotos,
} = useSort(photos);
const {
  selectedCount,
  headerChecked,
  isSelected,
  toggleOne,
  toggleAll,
  clearSelection,
} = useBulkSelection(photos);
const bulk = useBulkAction();
const {
  page,
  pageSize,
  totalPages,
  paginatedItems: pagedPhotos,
  setPageSize,
  resetPage,
} = usePagination(sortedPhotos);
const toggling = ref(false);

async function toggleUpload() {
  if (!config.value) return;
  toggling.value = true;
  try {
    await saveConfig({
      ...config.value,
      photo_upload_enabled: !config.value.photo_upload_enabled,
    });
  } finally {
    toggling.value = false;
  }
}

async function reload() {
  clearSelection();
  resetPage();
  await load();
}

function exportPhotos() {
  exportCSV(
    'photos.csv',
    [
      'ID',
      'Name',
      'Hidden',
      'Evaluated',
      'Hearts',
      'URL',
      'Original URL',
      'IP',
      'Date',
    ],
    photos.value.map((p) => [
      String(p.id),
      p.name,
      String(p.hidden),
      String(p.evaluated),
      String(p.hearts),
      p.url || '',
      p.original_url || '',
      p.ip_address,
      p.upload_date,
    ]),
  );
}

async function bulkHide() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => toggleVisibilityRaw(id, true), reload);
}

async function bulkShow() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => toggleVisibilityRaw(id, false), reload);
}

async function bulkDelete() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => deletePhotoRaw(id), reload);
}

onMounted(() => {
  load();
  loadConfig();
});

function formatDate(str: string) {
  return new Date(str).toLocaleString();
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('nav.photos') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('photos.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <Button
          v-if="config"
          :variant="config.photo_upload_enabled ? 'default' : 'outline'"
          size="sm"
          :disabled="toggling"
          @click="toggleUpload"
        >
          <UploadCloud
            v-if="config.photo_upload_enabled"
            class="mr-2 h-4 w-4"
          />
          <Upload v-else class="mr-2 h-4 w-4" />
          {{
            config.photo_upload_enabled
              ? t('photos.uploadsEnabled')
              : t('photos.enableUploads')
          }}
        </Button>
        <Button
          variant="outline"
          size="sm"
          :disabled="!photos.length"
          @click="exportPhotos"
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
        :description="t('photos.bulkDeleteDesc', { count: selectedCount })"
        @confirm="bulkDelete"
      >
        <Button size="sm" variant="destructive" :disabled="bulk.running.value">
          <Trash2 class="mr-2 h-4 w-4" />
          {{ t('common.delete') }}
        </Button>
      </ConfirmDialog>
    </BulkActionBar>

    <div v-if="loading && photos.length === 0" class="space-y-2">
      <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
    </div>

    <div
      v-else-if="photos.length === 0"
      class="py-8 text-center text-muted-foreground"
    >
      {{ t('photos.empty') }}
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
              :label="t('photos.name')"
              field="name"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>{{ t('common.originalFile') }}</TableHead>
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
              :label="t('photos.hearts')"
              field="hearts"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.date')"
              field="upload_date"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>{{ t('photos.preview') }}</TableHead>
          <TableHead class="w-16">{{ t('common.actions') }}</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="photo in pagedPhotos" :key="photo.id">
          <TableCell>
            <Checkbox
              :model-value="isSelected(photo.id)"
              :disabled="bulk.running.value"
              @update:model-value="toggleOne(photo.id, $event as boolean)"
            />
          </TableCell>
          <TableCell class="font-mono text-xs">{{ photo.id }}</TableCell>
          <TableCell class="font-medium">{{ photo.name }}</TableCell>
          <TableCell class="max-w-[200px] truncate text-xs">{{
            photo.original_filename
          }}</TableCell>
          <TableCell>
            <Badge :variant="photo.hidden ? 'destructive' : 'secondary'">
              {{ photo.hidden ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell>
            <Badge :variant="photo.evaluated ? 'default' : 'outline'">
              {{ photo.evaluated ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell class="font-mono text-xs">{{ photo.hearts }}</TableCell>
          <TableCell class="text-xs">{{
            formatDate(photo.upload_date)
          }}</TableCell>
          <TableCell>
            <div class="flex gap-2">
              <a
                v-if="photo.url"
                :href="photo.url"
                target="_blank"
                class="inline-flex items-center gap-1 text-xs text-primary hover:underline"
              >
                {{ t('photos.view') }} <ExternalLink class="h-3 w-3" />
              </a>
              <a
                v-if="photo.original_url"
                :href="photo.original_url"
                target="_blank"
                class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:underline"
              >
                {{ t('photos.original') }} <ExternalLink class="h-3 w-3" />
              </a>
            </div>
          </TableCell>
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
                  @click="toggleVisibility(photo.id, !photo.hidden)"
                >
                  <EyeOff v-if="!photo.hidden" class="mr-2 h-4 w-4" />
                  <Eye v-else class="mr-2 h-4 w-4" />
                  {{ photo.hidden ? t('common.show') : t('common.hide') }}
                </DropdownMenuItem>
                <ConfirmDialog
                  :title="t('photos.resetHeartsTitle')"
                  :description="
                    t('photos.resetHeartsDesc', { name: photo.name })
                  "
                  @confirm="resetHearts(photo.id)"
                >
                  <DropdownMenuItem @select.prevent>
                    <HeartOff class="mr-2 h-4 w-4" />
                    {{ t('photos.resetHeartsTitle') }}
                  </DropdownMenuItem>
                </ConfirmDialog>
                <ConfirmDialog
                  :title="t('photos.deleteTitle')"
                  :description="t('photos.deleteDesc', { name: photo.name })"
                  @confirm="deletePhoto(photo.id)"
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
      :total-items="photos.length"
      @update:page="page = $event"
      @update:page-size="setPageSize($event)"
    />
  </div>
</template>
