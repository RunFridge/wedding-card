<script setup lang="ts">
import { onMounted, onUnmounted, reactive, ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useSystemSettings } from '@/composables/useSystemSettings';
import { useConfig } from '@/composables/useConfig';
import { isDemo } from '@/lib/demo';
import type { SystemSettings } from '@/types/admin';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { Checkbox } from '@/components/ui/checkbox';
import ConfirmDialog from '@/components/ConfirmDialog.vue';
import {
  Save,
  Wifi,
  AlertTriangle,
  CheckCircle2,
  RotateCcw,
  Loader2,
} from 'lucide-vue-next';

const MODERATION_CATEGORIES = [
  'sexual',
  'hate',
  'harassment',
  'self-harm',
  'sexual/minors',
  'hate/threatening',
  'violence/graphic',
  'self-harm/intent',
  'self-harm/instructions',
  'harassment/threatening',
  'violence',
  'illicit',
  'illicit/violent',
] as const;

const { t } = useI18n();
const {
  settings,
  loading,
  saving,
  testing,
  error,
  load,
  save,
  testS3,
  restart,
} = useSystemSettings();
const { config, load: loadConfig, save: saveConfig } = useConfig();

const form = reactive({
  bcrypt_cost: 10,
  game_timer_ms: 30000,
  rate_limit_enabled: true,
  s3_bucket: '',
  s3_region: '',
  s3_endpoint: '',
  s3_access_key: '',
  s3_secret_key: '',
  use_moderation: false,
  openai_api_key: '',
  moderation_thresholds: {} as Record<string, number>,
});

const thresholdsDirty = ref(false);

const sensitiveSet = reactive({
  s3_access_key: false,
  s3_secret_key: false,
  openai_api_key: false,
});

const useLocalStorage = ref(true);

function switchToLocal() {
  useLocalStorage.value = true;
  form.s3_bucket = '';
  form.s3_region = '';
  form.s3_endpoint = '';
  form.s3_access_key = '';
  form.s3_secret_key = '';
  for (const f of [
    's3_bucket',
    's3_region',
    's3_endpoint',
    's3_access_key',
    's3_secret_key',
  ]) {
    markDirty(f);
  }
}

const success = ref(false);
const resultBanner = ref<{
  restart_required: boolean;
  s3_reinitialized: boolean;
} | null>(null);
const s3TestResult = ref<{ success: boolean; error?: string } | null>(null);
const dirtyFields = ref(new Set<string>());

const hasChanges = computed(
  () => dirtyFields.value.size > 0 || thresholdsDirty.value,
);

function markDirty(field: string) {
  dirtyFields.value.add(field);
}

function populateForm(s: SystemSettings) {
  form.bcrypt_cost = s.bcrypt_cost;
  form.game_timer_ms = s.game_timer_ms;
  form.rate_limit_enabled = s.rate_limit_enabled;
  form.s3_bucket = s.s3_bucket;
  form.s3_region = s.s3_region;
  form.s3_endpoint = s.s3_endpoint;
  form.s3_access_key = '';
  form.s3_secret_key = '';
  form.use_moderation = s.use_moderation;
  form.openai_api_key = '';
  sensitiveSet.s3_access_key = s.s3_access_key_set;
  sensitiveSet.s3_secret_key = s.s3_secret_key_set;
  sensitiveSet.openai_api_key = s.openai_api_key_set;
  form.moderation_thresholds = {
    ...(config.value?.moderation_thresholds ?? {}),
  };
  dirtyFields.value.clear();
  thresholdsDirty.value = false;
  useLocalStorage.value = !s.s3_bucket;
}

onMounted(async () => {
  await Promise.all([load(), loadConfig()]);
  if (settings.value) populateForm(settings.value);
});

async function handleSave() {
  success.value = false;
  resultBanner.value = null;

  let anyChange = false;

  if (thresholdsDirty.value) {
    // reload to minimize overwrites from concurrent edits in /-/admin/settings
    await loadConfig();
    if (config.value) {
      try {
        await saveConfig({
          ...config.value,
          moderation_thresholds: { ...form.moderation_thresholds },
        });
        thresholdsDirty.value = false;
        anyChange = true;
      } catch {
        return;
      }
    }
  }

  if (dirtyFields.value.size > 0) {
    const updates: Record<string, unknown> = {};
    for (const field of dirtyFields.value) {
      if (field === 'moderation_thresholds') continue;
      updates[field] = form[field as keyof typeof form];
    }
    const result = await save(updates);
    if (!result) return;
    resultBanner.value = result;
    if (settings.value) populateForm(settings.value);
    anyChange = true;
  }

  if (anyChange) {
    success.value = true;
    setTimeout(() => {
      success.value = false;
      resultBanner.value = null;
    }, 5000);
  }
}

async function handleTestS3() {
  s3TestResult.value = null;
  const result = await testS3({
    s3_endpoint: form.s3_endpoint,
    s3_region: form.s3_region,
    s3_bucket: form.s3_bucket,
    s3_access_key: form.s3_access_key,
    s3_secret_key: form.s3_secret_key,
  });
  if (result) {
    s3TestResult.value = result;
    setTimeout(() => (s3TestResult.value = null), 5000);
  }
}

const restarting = ref(false);
const restartCountdown = ref(10);
const restartFailed = ref(false);
let countdownTimer: ReturnType<typeof setInterval> | null = null;
let pollTimer: ReturnType<typeof setInterval> | null = null;

function cleanupTimers() {
  if (countdownTimer) {
    clearInterval(countdownTimer);
    countdownTimer = null;
  }
  if (pollTimer) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
}

onUnmounted(cleanupTimers);

function reloadPage() {
  window.location.reload();
}

async function handleRestart() {
  const ok = await restart();
  if (!ok) return;

  restarting.value = true;
  restartCountdown.value = 10;
  restartFailed.value = false;

  countdownTimer = setInterval(() => {
    if (restartCountdown.value > 0) restartCountdown.value--;
  }, 1000);

  setTimeout(() => {
    pollTimer = setInterval(async () => {
      try {
        const resp = await fetch('/api/health');
        if (resp.ok) {
          cleanupTimers();
          window.location.reload();
        }
      } catch {
        // server still down
      }
    }, 2000);
  }, 3000);

  setTimeout(() => {
    if (restarting.value) {
      cleanupTimers();
      restartFailed.value = true;
    }
  }, 30000);
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('system.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('system.subtitle') }}</p>
      </div>
    </div>

    <div v-if="loading" class="space-y-4">
      <Skeleton v-for="i in 3" :key="i" class="h-40 w-full" />
    </div>

    <div v-else class="space-y-6 pb-8">
      <!-- General Settings -->
      <Card>
        <CardHeader>
          <CardTitle>{{ t('system.general') }}</CardTitle>
        </CardHeader>
        <CardContent class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t('system.bcryptCost')
            }}</label>
            <Input
              v-model.number="form.bcrypt_cost"
              type="number"
              min="4"
              max="31"
              @input="markDirty('bcrypt_cost')"
            />
            <p class="mt-1 text-xs text-muted-foreground">
              {{ t('system.bcryptCostHint') }}
            </p>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t('system.gameTimer')
            }}</label>
            <Input
              v-model.number="form.game_timer_ms"
              type="number"
              min="1000"
              max="120000"
              step="1000"
              @input="markDirty('game_timer_ms')"
            />
            <p class="mt-1 text-xs text-muted-foreground">
              {{ t('system.gameTimerHint') }}
            </p>
          </div>
          <div class="sm:col-span-2 flex items-center justify-between">
            <div>
              <label class="block text-sm font-medium">{{
                t('system.rateLimiting')
              }}</label>
              <p class="text-xs text-muted-foreground">
                {{ t('system.rateLimitingHint') }}
              </p>
            </div>
            <Button
              :variant="form.rate_limit_enabled ? 'default' : 'outline'"
              size="sm"
              @click="
                form.rate_limit_enabled = !form.rate_limit_enabled;
                markDirty('rate_limit_enabled');
              "
            >
              {{
                form.rate_limit_enabled
                  ? t('system.enabled')
                  : t('system.disabled')
              }}
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Photo Storage -->
      <Card>
        <CardHeader>
          <CardTitle>{{ t('system.photoStorage') }}</CardTitle>
        </CardHeader>
        <CardContent class="grid gap-4">
          <div class="flex gap-2">
            <Button
              :variant="useLocalStorage ? 'default' : 'outline'"
              size="sm"
              @click="switchToLocal"
              >{{ t('storage.localFilesystem') }}</Button
            >
            <Button
              :variant="!useLocalStorage ? 'default' : 'outline'"
              size="sm"
              @click="useLocalStorage = false"
              >{{ t('storage.s3Storage') }}</Button
            >
          </div>
          <div
            v-if="useLocalStorage"
            class="flex items-start gap-2 rounded-md border bg-muted px-4 py-3 text-sm text-muted-foreground"
          >
            <p>
              {{ t('system.localStorageInfo') }}
            </p>
          </div>
          <template v-if="!useLocalStorage">
            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('storage.bucket')
                }}</label>
                <Input
                  v-model="form.s3_bucket"
                  placeholder="my-wedding-bucket"
                  @input="markDirty('s3_bucket')"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('storage.region')
                }}</label>
                <Input
                  v-model="form.s3_region"
                  placeholder="us-east-1"
                  @input="markDirty('s3_region')"
                />
              </div>
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium">{{
                t('storage.endpoint')
              }}</label>
              <Input
                v-model="form.s3_endpoint"
                placeholder="https://s3.example.com (leave empty for AWS)"
                @input="markDirty('s3_endpoint')"
              />
            </div>
            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('storage.accessKey')
                }}</label>
                <Input
                  v-model="form.s3_access_key"
                  type="password"
                  :placeholder="
                    sensitiveSet.s3_access_key
                      ? t('system.configured')
                      : t('system.notSet')
                  "
                  @input="markDirty('s3_access_key')"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('storage.secretKey')
                }}</label>
                <Input
                  v-model="form.s3_secret_key"
                  type="password"
                  :placeholder="
                    sensitiveSet.s3_secret_key
                      ? t('system.configured')
                      : t('system.notSet')
                  "
                  @input="markDirty('s3_secret_key')"
                />
              </div>
            </div>
            <div class="flex items-center gap-3">
              <Button
                variant="outline"
                :disabled="testing || !form.s3_bucket || isDemo"
                @click="handleTestS3"
              >
                <Wifi class="mr-2 h-4 w-4" />
                {{
                  testing ? t('storage.testing') : t('storage.testConnection')
                }}
              </Button>
              <span
                v-if="s3TestResult?.success"
                class="flex items-center gap-1 text-sm text-green-600"
              >
                <CheckCircle2 class="h-4 w-4" />
                {{ t('storage.connectionSuccessful') }}
              </span>
              <span
                v-if="s3TestResult && !s3TestResult.success"
                class="flex items-center gap-1 text-sm text-destructive"
              >
                <AlertTriangle class="h-4 w-4" /> {{ s3TestResult.error }}
              </span>
            </div>
          </template>
        </CardContent>
      </Card>

      <!-- Content Moderation -->
      <Card>
        <CardHeader>
          <CardTitle>{{ t('moderation.title') }}</CardTitle>
          <p class="text-sm text-muted-foreground">
            {{ t('system.moderationRestartHint') }}
          </p>
        </CardHeader>
        <CardContent class="grid gap-4">
          <div class="flex items-center gap-3">
            <Checkbox
              :model-value="form.use_moderation"
              @update:model-value="
                (val: boolean | 'indeterminate') => {
                  form.use_moderation = val === true;
                  markDirty('use_moderation');
                }
              "
            />
            <label class="text-sm font-medium">{{
              t('moderation.enable')
            }}</label>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium">{{
                t('moderation.openaiApiKey')
              }}</label>
              <Input
                v-model="form.openai_api_key"
                type="password"
                :placeholder="
                  sensitiveSet.openai_api_key
                    ? t('system.configured')
                    : t('system.notSet')
                "
                @input="markDirty('openai_api_key')"
              />
            </div>
          </div>

          <Separator class="my-2" />

          <div :class="{ 'opacity-50': !form.use_moderation }">
            <div class="mb-3">
              <h4 class="text-sm font-medium">{{ t('system.thresholds') }}</h4>
              <p class="mt-1 text-xs text-muted-foreground">
                {{ t('system.thresholdsHint') }}
              </p>
            </div>
            <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              <div v-for="category in MODERATION_CATEGORIES" :key="category">
                <label class="mb-1 block text-sm font-medium">{{
                  category
                }}</label>
                <Input
                  v-model.number="form.moderation_thresholds[category]"
                  type="number"
                  min="0"
                  max="1"
                  step="0.01"
                  :disabled="!form.use_moderation"
                  @input="thresholdsDirty = true"
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Result Banners -->
      <div
        v-if="resultBanner?.restart_required"
        class="flex items-center justify-between gap-2 rounded-md border border-yellow-300 bg-yellow-50 px-4 py-3 text-sm text-yellow-800 dark:border-yellow-700 dark:bg-yellow-950 dark:text-yellow-200"
      >
        <div class="flex items-center gap-2">
          <AlertTriangle class="h-4 w-4 shrink-0" />
          {{ t('system.restartRequired') }}
        </div>
        <ConfirmDialog
          :title="t('system.restartConfirmTitle')"
          :description="t('system.restartConfirmDesc')"
          @confirm="handleRestart"
        >
          <Button variant="destructive" size="sm" :disabled="isDemo">
            <RotateCcw class="mr-2 h-4 w-4" />
            {{ t('system.restartButton') }}
          </Button>
        </ConfirmDialog>
      </div>
      <div
        v-if="resultBanner?.s3_reinitialized"
        class="flex items-center gap-2 rounded-md border border-green-300 bg-green-50 px-4 py-3 text-sm text-green-800 dark:border-green-700 dark:bg-green-950 dark:text-green-200"
      >
        <CheckCircle2 class="h-4 w-4 shrink-0" />
        {{ t('system.s3Reinitialized') }}
      </div>

      <!-- Save -->
      <div class="flex items-center gap-3">
        <Button :disabled="saving || !hasChanges || isDemo" @click="handleSave">
          <Save class="mr-2 h-4 w-4" />
          {{ saving ? t('common.saving') : t('common.saveSettings') }}
        </Button>
        <span v-if="isDemo" class="text-sm text-muted-foreground">
          {{ t('demo.disabled') }}
        </span>
        <span
          v-if="
            success &&
            !resultBanner?.restart_required &&
            !resultBanner?.s3_reinitialized
          "
          class="text-sm text-green-600"
        >
          {{ t('common.settingsSaved') }}
        </span>
        <span v-if="error" class="text-sm text-destructive">{{ error }}</span>
      </div>
    </div>

    <!-- Restart overlay -->
    <div
      v-if="restarting"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
    >
      <Card class="w-80 text-center">
        <CardContent class="flex flex-col items-center gap-4 pt-6">
          <template v-if="!restartFailed">
            <Loader2 class="h-10 w-10 animate-spin text-primary" />
            <h2 class="text-lg font-semibold">{{ t('restart.restarting') }}</h2>
            <p class="text-sm text-muted-foreground">
              {{ t('restart.reconnectingIn', { seconds: restartCountdown }) }}
            </p>
          </template>
          <template v-else>
            <AlertTriangle class="h-10 w-10 text-destructive" />
            <h2 class="text-lg font-semibold">{{ t('restart.noResponse') }}</h2>
            <p class="text-sm text-muted-foreground">
              {{ t('restart.reloadManually') }}
            </p>
            <Button @click="reloadPage">{{ t('common.reload') }}</Button>
          </template>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
