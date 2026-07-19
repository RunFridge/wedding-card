<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useAuth } from '@/lib/auth';
import { useSystemSettings } from '@/composables/useSystemSettings';
import api from '@/lib/axios';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Loader2 } from 'lucide-vue-next';
import {
  Wifi,
  AlertTriangle,
  CheckCircle2,
  ChevronRight,
  ExternalLink,
  Info,
} from 'lucide-vue-next';

const router = useRouter();
const { t } = useI18n();
const { passwordNeedsChange, clearPasswordNeedsChange, clearSetupRequired } =
  useAuth();
const { save, testS3, testModeration, restart } = useSystemSettings();

const step = ref(1);
const totalSteps = 3;

// Step 1: Password
const passwordForm = reactive({
  newPassword: '',
  confirm: '',
});
const passwordError = ref('');
const passwordLoading = ref(false);
const showPassword = ref(false);

// Step 2: Storage
const useLocalStorage = ref(true);
const s3Form = reactive({
  s3_bucket: '',
  s3_region: '',
  s3_endpoint: '',
  s3_access_key: '',
  s3_secret_key: '',
});
const s3Error = ref('');
const s3Loading = ref(false);
const s3TestResult = ref<{ success: boolean; error?: string } | null>(null);
const s3Testing = ref(false);

// Step 3: Moderation
const moderationForm = reactive({
  use_moderation: false,
  openai_api_key: '',
});
const moderationError = ref('');
const moderationLoading = ref(false);
const moderationTestResult = ref<{ success: boolean; error?: string } | null>(
  null,
);
const moderationTesting = ref(false);

const completing = ref(false);
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

onMounted(() => {
  if (!passwordNeedsChange.value) {
    step.value = 2;
  }
});

async function handlePasswordSubmit() {
  passwordError.value = '';

  if (!passwordForm.newPassword || !passwordForm.confirm) {
    passwordError.value = t('passwordForm.allFieldsRequired');
    return;
  }
  if (passwordForm.newPassword.length < 8) {
    passwordError.value = t('passwordForm.tooShort');
    return;
  }
  if (passwordForm.newPassword !== passwordForm.confirm) {
    passwordError.value = t('passwordForm.mismatch');
    return;
  }

  passwordLoading.value = true;
  try {
    await api.put('/admin/password', {
      new_password: passwordForm.newPassword,
    });
    clearPasswordNeedsChange();
    step.value = 2;
  } catch {
    passwordError.value = t('setup.passwordSaveFailed');
  } finally {
    passwordLoading.value = false;
  }
}

async function handleTestS3() {
  s3TestResult.value = null;
  s3Testing.value = true;
  try {
    const result = await testS3({
      s3_endpoint: s3Form.s3_endpoint,
      s3_region: s3Form.s3_region,
      s3_bucket: s3Form.s3_bucket,
      s3_access_key: s3Form.s3_access_key,
      s3_secret_key: s3Form.s3_secret_key,
    });
    if (result) s3TestResult.value = result;
  } finally {
    s3Testing.value = false;
  }
}

async function handleTestModeration() {
  moderationTestResult.value = null;
  moderationTesting.value = true;
  try {
    const result = await testModeration({
      openai_api_key: moderationForm.openai_api_key,
    });
    if (result) moderationTestResult.value = result;
  } finally {
    moderationTesting.value = false;
  }
}

async function handleS3Save() {
  s3Error.value = '';
  s3Loading.value = true;
  try {
    if (useLocalStorage.value) {
      await save({
        s3_bucket: '',
        s3_region: '',
        s3_endpoint: '',
        s3_access_key: '',
        s3_secret_key: '',
      });
    } else {
      await save({
        s3_bucket: s3Form.s3_bucket,
        s3_region: s3Form.s3_region,
        s3_endpoint: s3Form.s3_endpoint,
        s3_access_key: s3Form.s3_access_key,
        s3_secret_key: s3Form.s3_secret_key,
      });
    }
    step.value = 3;
  } catch {
    s3Error.value = t('setup.storageSaveFailed');
  } finally {
    s3Loading.value = false;
  }
}

async function handleModerationSave() {
  moderationError.value = '';
  moderationLoading.value = true;
  try {
    await save({
      use_moderation: moderationForm.use_moderation,
      openai_api_key: moderationForm.openai_api_key,
    });
    await completeSetup(moderationForm.use_moderation);
  } catch {
    moderationError.value = t('setup.moderationSaveFailed');
  } finally {
    moderationLoading.value = false;
  }
}

async function completeSetup(shouldRestart = false) {
  completing.value = true;
  try {
    await api.post('/admin/setup/complete');
    clearSetupRequired();

    if (shouldRestart) {
      const ok = await restart();
      if (ok) {
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
                window.location.href = '/-/admin/';
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
        return;
      }
    }

    router.push({ name: 'dashboard' });
  } catch {
    clearSetupRequired();
    router.push({ name: 'dashboard' });
  } finally {
    completing.value = false;
  }
}

function stepClass(s: number) {
  if (s < step.value) return 'bg-primary text-primary-foreground';
  if (s === step.value)
    return 'bg-primary text-primary-foreground ring-2 ring-primary ring-offset-2';
  return 'bg-muted text-muted-foreground';
}

function lineClass(s: number) {
  return s < step.value ? 'bg-primary' : 'bg-muted';
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-lg space-y-6">
      <!-- Stepper -->
      <div class="flex items-center justify-center gap-2">
        <template v-for="s in totalSteps" :key="s">
          <div
            :class="[
              'flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold transition-colors',
              stepClass(s),
            ]"
          >
            {{ s }}
          </div>
          <div
            v-if="s < totalSteps"
            :class="['h-0.5 w-12 transition-colors', lineClass(s)]"
          />
        </template>
      </div>

      <!-- Step 1: Password -->
      <Card v-if="step === 1">
        <CardHeader>
          <CardTitle class="text-2xl">{{ t('setup.passwordTitle') }}</CardTitle>
          <CardDescription>
            {{ t('setup.passwordDesc') }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form class="space-y-4" @submit.prevent="handlePasswordSubmit">
            <div class="space-y-2">
              <label class="text-sm font-medium" for="setup-new-password">{{
                t('passwordForm.newPassword')
              }}</label>
              <div class="relative">
                <Input
                  id="setup-new-password"
                  v-model="passwordForm.newPassword"
                  :type="showPassword ? 'text' : 'password'"
                  :placeholder="t('passwordForm.atLeast8')"
                  autofocus
                />
                <button
                  type="button"
                  tabindex="-1"
                  class="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-muted-foreground hover:text-foreground"
                  @click="showPassword = !showPassword"
                >
                  {{ showPassword ? t('setup.hide') : t('setup.show') }}
                </button>
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium" for="setup-confirm-password">{{
                t('setup.confirmPassword')
              }}</label>
              <div class="relative">
                <Input
                  id="setup-confirm-password"
                  v-model="passwordForm.confirm"
                  :type="showPassword ? 'text' : 'password'"
                  :placeholder="t('passwordForm.repeatNew')"
                />
                <button
                  type="button"
                  tabindex="-1"
                  class="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-muted-foreground hover:text-foreground"
                  @click="showPassword = !showPassword"
                >
                  {{ showPassword ? t('setup.hide') : t('setup.show') }}
                </button>
              </div>
            </div>
            <p v-if="passwordError" class="text-sm text-destructive">
              {{ passwordError }}
            </p>
            <Button type="submit" class="w-full" :disabled="passwordLoading">
              <template v-if="passwordLoading">{{
                t('passwordForm.changing')
              }}</template>
              <template v-else>
                {{ t('setup.setPasswordContinue') }}
                <ChevronRight class="ml-2 h-4 w-4" />
              </template>
            </Button>
          </form>
        </CardContent>
      </Card>

      <!-- Step 2: Photo Storage -->
      <Card v-if="step === 2">
        <CardHeader>
          <CardTitle class="text-2xl">{{ t('setup.storageTitle') }}</CardTitle>
          <CardDescription>
            {{ t('setup.storageDesc') }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex gap-2">
              <Button
                :variant="useLocalStorage ? 'default' : 'outline'"
                size="sm"
                @click="useLocalStorage = true"
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
              <Info class="h-4 w-4 mt-0.5 shrink-0" />
              <p>
                {{ t('setup.localStorageInfo') }}
              </p>
            </div>

            <template v-if="!useLocalStorage">
              <div
                class="flex items-start gap-2 rounded-md border bg-muted px-4 py-3 text-sm text-muted-foreground"
              >
                <Info class="h-4 w-4 mt-0.5 shrink-0" />
                <div>
                  <p>
                    {{ t('setup.s3Compatible') }}
                  </p>
                  <i18n-t keypath="setup.s3Requirements" tag="p" class="mt-1">
                    <template #bucket>
                      <strong>{{ t('setup.reqBucket') }}</strong>
                    </template>
                    <template #region>
                      <strong>{{ t('setup.reqRegion') }}</strong>
                    </template>
                    <template #keys>
                      <strong>{{ t('setup.reqKeys') }}</strong>
                    </template>
                  </i18n-t>
                  <p class="mt-1">
                    <a
                      href="https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html"
                      target="_blank"
                      rel="noopener"
                      class="inline-flex items-center gap-1 underline hover:text-foreground"
                    >
                      {{ t('setup.s3Guide') }} <ExternalLink class="h-3 w-3" />
                    </a>
                  </p>
                </div>
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{
                    t('storage.bucket')
                  }}</label>
                  <Input
                    v-model="s3Form.s3_bucket"
                    placeholder="my-wedding-bucket"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{
                    t('storage.region')
                  }}</label>
                  <Input v-model="s3Form.s3_region" placeholder="us-east-1" />
                </div>
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium">{{
                  t('storage.endpoint')
                }}</label>
                <Input
                  v-model="s3Form.s3_endpoint"
                  placeholder="https://s3.example.com (leave empty for AWS)"
                />
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{
                    t('storage.accessKey')
                  }}</label>
                  <Input
                    v-model="s3Form.s3_access_key"
                    type="password"
                    :placeholder="t('storage.accessKey')"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{
                    t('storage.secretKey')
                  }}</label>
                  <Input
                    v-model="s3Form.s3_secret_key"
                    type="password"
                    :placeholder="t('storage.secretKey')"
                  />
                </div>
              </div>
              <div class="flex items-center gap-3">
                <Button
                  variant="outline"
                  :disabled="s3Testing || !s3Form.s3_bucket"
                  @click="handleTestS3"
                >
                  <Wifi class="mr-2 h-4 w-4" />
                  {{
                    s3Testing
                      ? t('storage.testing')
                      : t('storage.testConnection')
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
            <p v-if="s3Error" class="text-sm text-destructive">{{ s3Error }}</p>
            <div class="flex gap-3">
              <Button
                :disabled="s3Loading || (!useLocalStorage && !s3Form.s3_bucket)"
                @click="handleS3Save"
              >
                <template v-if="s3Loading">{{ t('common.saving') }}</template>
                <template v-else>
                  {{
                    useLocalStorage ? t('setup.next') : t('setup.saveAndNext')
                  }}
                  <ChevronRight class="ml-2 h-4 w-4" />
                </template>
              </Button>
              <Button
                v-if="!useLocalStorage"
                variant="outline"
                @click="step = 3"
              >
                {{ t('setup.skip') }}
                <ChevronRight class="ml-2 h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Step 3: Moderation -->
      <Card v-if="step === 3">
        <CardHeader>
          <CardTitle class="text-2xl">{{ t('moderation.title') }}</CardTitle>
          <CardDescription>
            {{ t('setup.moderationDesc') }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <Checkbox
                :model-value="moderationForm.use_moderation"
                @update:model-value="
                  (val: boolean | 'indeterminate') => {
                    moderationForm.use_moderation = val === true;
                  }
                "
              />
              <label class="text-sm font-medium">{{
                t('moderation.enable')
              }}</label>
            </div>
            <div v-if="moderationForm.use_moderation" class="grid gap-4">
              <div class="space-y-2">
                <label class="text-sm font-medium">{{
                  t('moderation.openaiApiKey')
                }}</label>
                <Input
                  v-model="moderationForm.openai_api_key"
                  type="password"
                  placeholder="sk-..."
                />
                <p class="text-xs text-muted-foreground">
                  {{ t('setup.getApiKey') }}
                  <a
                    href="https://platform.openai.com/api-keys"
                    target="_blank"
                    rel="noopener"
                    class="inline-flex items-center gap-0.5 underline hover:text-foreground"
                  >
                    platform.openai.com/api-keys
                    <ExternalLink class="h-3 w-3" />
                  </a>
                </p>
                <div class="flex items-center gap-3 mt-2">
                  <Button
                    variant="outline"
                    size="sm"
                    :disabled="
                      moderationTesting || !moderationForm.openai_api_key
                    "
                    @click="handleTestModeration"
                  >
                    <Wifi class="mr-2 h-4 w-4" />
                    {{
                      moderationTesting
                        ? t('storage.testing')
                        : t('setup.testApiKey')
                    }}
                  </Button>
                  <span
                    v-if="moderationTestResult?.success"
                    class="flex items-center gap-1 text-sm text-green-600"
                  >
                    <CheckCircle2 class="h-4 w-4" />
                    {{ t('setup.apiKeyValid') }}
                  </span>
                  <span
                    v-if="moderationTestResult && !moderationTestResult.success"
                    class="flex items-center gap-1 text-sm text-destructive"
                  >
                    <AlertTriangle class="h-4 w-4" />
                    {{ moderationTestResult.error }}
                  </span>
                </div>
              </div>
            </div>
            <div
              v-if="moderationForm.use_moderation"
              class="flex items-center gap-2 rounded-md border border-yellow-300 bg-yellow-50 px-4 py-3 text-sm text-yellow-800 dark:border-yellow-700 dark:bg-yellow-950 dark:text-yellow-200"
            >
              <AlertTriangle class="h-4 w-4 shrink-0" />
              {{ t('setup.restartNotice') }}
            </div>
            <p v-if="moderationError" class="text-sm text-destructive">
              {{ moderationError }}
            </p>
            <div class="flex gap-3">
              <Button
                :disabled="moderationLoading || completing"
                @click="
                  moderationForm.use_moderation
                    ? handleModerationSave()
                    : completeSetup()
                "
              >
                <template v-if="moderationLoading || completing">{{
                  t('setup.finishing')
                }}</template>
                <template v-else>
                  {{
                    moderationForm.use_moderation
                      ? t('setup.saveAndComplete')
                      : t('setup.completeSetup')
                  }}
                </template>
              </Button>
              <Button
                v-if="moderationForm.use_moderation"
                variant="outline"
                :disabled="completing"
                @click="completeSetup()"
              >
                {{
                  completing ? t('setup.finishing') : t('setup.skipAndComplete')
                }}
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
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
