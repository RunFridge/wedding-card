<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAssetPhotos } from '@/composables/useAssetPhotos';
import { useSort } from '@/composables/useSort';
import { usePagination } from '@/composables/usePagination';
import PaginationBar from '@/components/PaginationBar.vue';
import SortableHeader from '@/components/SortableHeader.vue';
import type { BatchProgress } from '@/composables/useAssetPhotos';
import { useBulkSelection } from '@/composables/useBulkSelection';
import { useBulkAction } from '@/composables/useBulkAction';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
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
  Trash2,
  ExternalLink,
  RefreshCw,
  Upload,
  Gamepad2,
  Star,
  StarOff,
  X,
} from 'lucide-vue-next';

const { t } = useI18n();
const {
  photos,
  loading,
  load,
  uploadBatch,
  toggleGame,
  toggleGameRaw,
  setMain,
  deletePhoto,
  deletePhotoRaw,
} = useAssetPhotos();
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

const fileInput = ref<HTMLInputElement | null>(null);
const uploadLabel = ref('');
const uploading = ref(false);
const selectedFiles = ref<File[]>([]);
const uploadProgress = ref<BatchProgress | null>(null);
const uploadErrors = ref<{ fileName: string; message: string }[]>([]);

const gamePhotoCount = computed(
  () => photos.value.filter((p) => p.use_for_game).length,
);
const hasMainPhoto = computed(() => photos.value.some((p) => p.is_main_photo));

const uploadButtonText = computed(() => {
  if (uploading.value && uploadProgress.value) {
    return t('assetPhotos.uploadingProgress', {
      current: uploadProgress.value.completed + 1,
      total: uploadProgress.value.total,
    });
  }
  const count = selectedFiles.value.length;
  if (count <= 1) return t('assetPhotos.upload');
  return t('assetPhotos.uploadN', { count });
});

async function reload() {
  clearSelection();
  resetPage();
  await load();
}

async function handleUpload() {
  if (selectedFiles.value.length === 0) return;
  uploading.value = true;
  uploadErrors.value = [];
  try {
    const errors = await uploadBatch(
      selectedFiles.value,
      uploadLabel.value,
      (progress) => {
        uploadProgress.value = { ...progress };
      },
    );
    if (errors.length === 0) {
      uploadLabel.value = '';
      selectedFiles.value = [];
      if (fileInput.value) fileInput.value.value = '';
    }
    uploadErrors.value = errors;
  } finally {
    uploading.value = false;
    uploadProgress.value = null;
  }
}

async function bulkAddToGame() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => toggleGameRaw(id, true), reload);
}

async function bulkRemoveFromGame() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => toggleGameRaw(id, false), reload);
}

async function bulkDelete() {
  const ids = [
    ...photos.value.filter((p) => isSelected(p.id)).map((p) => p.id),
  ];
  await bulk.execute(ids, (id) => deletePhotoRaw(id), reload);
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement;
  selectedFiles.value = target.files ? Array.from(target.files) : [];
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
        <h1 class="text-2xl font-bold">{{ t('assetPhotos.title') }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ t('assetPhotos.subtitle') }}
        </p>
      </div>
      <Button
        variant="outline"
        size="sm"
        :disabled="loading || bulk.running.value"
        @click="reload"
      >
        <RefreshCw class="mr-2 h-4 w-4" :class="{ 'animate-spin': loading }" />
        {{ t('common.reload') }}
      </Button>
    </div>

    <div class="mb-4 flex items-center gap-3">
      <Badge variant="outline">{{
        t('assetPhotos.gamePhotoCount', { count: gamePhotoCount })
      }}</Badge>
      <Badge :variant="hasMainPhoto ? 'default' : 'outline'">
        {{
          hasMainPhoto
            ? t('assetPhotos.mainPhotoSet')
            : t('assetPhotos.mainPhotoNotSet')
        }}
      </Badge>
    </div>

    <div class="mb-6 rounded-lg border p-4">
      <div class="flex items-end gap-3">
        <div class="flex-1">
          <label class="mb-1 block text-sm font-medium">{{
            t('assetPhotos.image')
          }}</label>
          <input
            ref="fileInput"
            type="file"
            accept="image/*"
            multiple
            class="block w-full text-sm file:mr-4 file:rounded file:border-0 file:bg-primary file:px-4 file:py-2 file:text-sm file:font-semibold file:text-primary-foreground hover:file:bg-primary/90"
            @change="onFileChange"
          />
          <p
            v-if="selectedFiles.length > 1"
            class="mt-1 text-xs text-muted-foreground"
          >
            {{
              t('assetPhotos.filesSelected', { count: selectedFiles.length })
            }}
          </p>
        </div>
        <div class="w-48">
          <label class="mb-1 block text-sm font-medium">{{
            t('assetPhotos.labelOptional')
          }}</label>
          <Input
            v-model="uploadLabel"
            :placeholder="t('assetPhotos.labelPlaceholder')"
          />
        </div>
        <Button
          :disabled="selectedFiles.length === 0 || uploading"
          @click="handleUpload"
        >
          <Upload class="mr-2 h-4 w-4" />
          {{ uploadButtonText }}
        </Button>
      </div>

      <div v-if="uploading && uploadProgress" class="mt-3">
        <div
          class="mb-1 flex items-center justify-between text-xs text-muted-foreground"
        >
          <span>{{
            t('assetPhotos.uploadingFile', {
              name: uploadProgress.currentFileName,
            })
          }}</span>
          <span>{{ uploadProgress.completed }}/{{ uploadProgress.total }}</span>
        </div>
        <div class="h-2 w-full overflow-hidden rounded-full bg-muted">
          <div
            class="h-full bg-primary transition-all duration-300"
            :style="{
              width: `${(uploadProgress.completed / uploadProgress.total) * 100}%`,
            }"
          />
        </div>
      </div>

      <div
        v-if="uploadErrors.length > 0"
        class="mt-3 rounded-md border border-destructive/50 bg-destructive/10 p-3"
      >
        <div class="mb-1 flex items-center justify-between">
          <span class="text-sm font-medium text-destructive">
            {{ t('assetPhotos.uploadsFailed', uploadErrors.length) }}
          </span>
          <Button
            variant="ghost"
            size="icon"
            class="h-6 w-6"
            @click="uploadErrors = []"
          >
            <X class="h-3 w-3" />
          </Button>
        </div>
        <ul class="space-y-1 text-xs text-destructive">
          <li v-for="(err, i) in uploadErrors" :key="i">
            <span class="font-medium">{{ err.fileName }}:</span>
            {{ err.message }}
          </li>
        </ul>
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
        @click="bulkAddToGame"
      >
        <Gamepad2 class="mr-2 h-4 w-4" />
        {{ t('assetPhotos.addToGame') }}
      </Button>
      <Button
        size="sm"
        variant="outline"
        :disabled="bulk.running.value"
        @click="bulkRemoveFromGame"
      >
        <Gamepad2 class="mr-2 h-4 w-4" />
        {{ t('assetPhotos.removeFromGame') }}
      </Button>
      <ConfirmDialog
        :title="t('common.bulkDelete')"
        :description="t('assetPhotos.bulkDeleteDesc', { count: selectedCount })"
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
      {{ t('assetPhotos.empty') }}
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
          <TableHead class="w-16">{{ t('assetPhotos.thumb') }}</TableHead>
          <TableHead class="w-12">{{ t('common.id') }}</TableHead>
          <TableHead>
            <SortableHeader
              :label="t('assetPhotos.label')"
              field="label"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('common.originalFile')"
              field="original_filename"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('assetPhotos.game')"
              field="use_for_game"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead>
            <SortableHeader
              :label="t('assetPhotos.main')"
              field="is_main_photo"
              :active-field="sortKey"
              :direction="sortDir"
              @sort="toggleSort"
            />
          </TableHead>
          <TableHead class="w-12">
            <SortableHeader
              :label="t('assetPhotos.sort')"
              field="sort_order"
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
        <TableRow v-for="photo in pagedPhotos" :key="photo.id">
          <TableCell>
            <Checkbox
              :model-value="isSelected(photo.id)"
              :disabled="bulk.running.value"
              @update:model-value="toggleOne(photo.id, $event as boolean)"
            />
          </TableCell>
          <TableCell>
            <a v-if="photo.thumbnail_url" :href="photo.url" target="_blank">
              <img
                :src="photo.thumbnail_url"
                :alt="photo.label || photo.original_filename"
                class="h-10 w-10 rounded object-cover"
              />
            </a>
          </TableCell>
          <TableCell class="font-mono text-xs">{{ photo.id }}</TableCell>
          <TableCell class="font-medium">{{ photo.label || '-' }}</TableCell>
          <TableCell class="max-w-[200px] truncate text-xs">{{
            photo.original_filename
          }}</TableCell>
          <TableCell>
            <Badge
              :variant="photo.use_for_game ? 'default' : 'outline'"
              class="cursor-pointer"
              @click="toggleGame(photo.id, !photo.use_for_game)"
            >
              {{ photo.use_for_game ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell>
            <Badge
              :variant="photo.is_main_photo ? 'default' : 'outline'"
              class="cursor-pointer"
              @click="setMain(photo.id, !photo.is_main_photo)"
            >
              {{ photo.is_main_photo ? t('common.yes') : t('common.no') }}
            </Badge>
          </TableCell>
          <TableCell class="font-mono text-xs">{{
            photo.sort_order
          }}</TableCell>
          <TableCell class="text-xs">{{
            formatDate(photo.created_at)
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
                  @click="toggleGame(photo.id, !photo.use_for_game)"
                >
                  <Gamepad2 class="mr-2 h-4 w-4" />
                  {{
                    photo.use_for_game
                      ? t('assetPhotos.removeFromGame')
                      : t('assetPhotos.addToGame')
                  }}
                </DropdownMenuItem>
                <DropdownMenuItem
                  @click="setMain(photo.id, !photo.is_main_photo)"
                >
                  <Star v-if="!photo.is_main_photo" class="mr-2 h-4 w-4" />
                  <StarOff v-else class="mr-2 h-4 w-4" />
                  {{
                    photo.is_main_photo
                      ? t('assetPhotos.unsetMain')
                      : t('assetPhotos.setAsMain')
                  }}
                </DropdownMenuItem>
                <DropdownMenuItem
                  v-if="photo.url"
                  as="a"
                  :href="photo.url"
                  target="_blank"
                >
                  <ExternalLink class="mr-2 h-4 w-4" />
                  {{ t('assetPhotos.viewFull') }}
                </DropdownMenuItem>
                <ConfirmDialog
                  :title="t('assetPhotos.deleteTitle')"
                  :description="
                    t('assetPhotos.deleteDesc', {
                      name: photo.label || photo.original_filename,
                    })
                  "
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
