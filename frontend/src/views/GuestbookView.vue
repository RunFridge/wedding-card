<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">{{ t('guestbook.title') }}</h1>
        <p class="text-wood-dark text-sm break-keep">
          {{ t('guestbook.subtitle') }}
        </p>
      </div>
    </div>

    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-4">
        <form @submit.prevent="submitEntry">
          <div class="mb-3">
            <label class="block text-secondary text-sm mb-1">{{
              t('guestbook.nameLabel')
            }}</label>
            <input
              v-model="nickname"
              type="text"
              maxlength="50"
              required
              class="pixel-input w-full"
              :placeholder="t('common.namePlaceholder')"
            />
          </div>
          <div class="mb-3">
            <label class="block text-secondary text-sm mb-1">{{
              t('guestbook.messageLabel')
            }}</label>
            <textarea
              v-model="message"
              maxlength="500"
              required
              rows="4"
              class="pixel-input w-full resize-none"
              :placeholder="t('guestbook.messagePlaceholder')"
            ></textarea>
            <p class="text-xs text-wood-dark/60 text-right mt-1">
              {{ message.length }}/500
            </p>
          </div>
          <div class="mb-3">
            <label class="block text-secondary text-sm mb-1">{{
              t('guestbook.passwordLabel')
            }}</label>
            <input
              v-model="password"
              type="password"
              minlength="5"
              maxlength="30"
              required
              class="pixel-input w-full"
              :placeholder="t('guestbook.passwordPlaceholder')"
            />
          </div>
          <label
            class="flex items-center gap-2 mb-1 cursor-pointer text-sm text-wood-dark/80 select-none"
            @click.prevent="secret = !secret"
          >
            <span class="pixel-checkbox" :class="{ 'is-checked': secret }">
              <span v-if="secret" class="pixel-checkmark">&#10003;</span>
            </span>
            <TwEmoji emoji="🔒" size="0.875rem" /> {{ t('guestbook.secret') }}
          </label>
          <p class="text-xs text-wood-dark/50 mb-3 ml-8">
            {{ t('guestbook.secretHint') }}
          </p>
          <button type="submit" :disabled="submitting" class="pixel-btn w-full">
            {{ submitting ? t('guestbook.submitting') : t('guestbook.submit') }}
          </button>
        </form>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-6">
      <div class="pixel-spinner"></div>
    </div>

    <div
      v-else-if="entries.length === 0"
      class="parchment-bg p-6 text-center text-wood-dark/80"
    >
      {{ t('guestbook.emptyLine1') }}<br />{{ t('guestbook.emptyLine2') }}
    </div>

    <div v-else class="space-y-3">
      <div v-for="entry in entries" :key="entry.id" class="parchment-bg p-4">
        <div class="flex justify-between items-start mb-2">
          <div class="flex items-center gap-2">
            <TwEmoji
              v-if="entry.secret"
              emoji="🔒"
              size="1.25rem"
              class="flex items-center justify-center w-7 h-7"
            />
            <BoringAvatar
              v-else
              :name="entry.nickname"
              :size="28"
              class="rounded-full border-2 border-wood-dark/20"
            />
            <span class="text-secondary text-sm font-semibold">{{
              entry.nickname
            }}</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="text-xs text-wood-dark/50 mr-1">{{
              formatDate(entry.created_at)
            }}</span>
            <template v-if="!entry.secret">
              <button
                class="text-wood-dark/40 hover:text-secondary text-xs px-1"
                :aria-label="t('guestbook.editAria')"
                @click="startEdit(entry)"
              >
                ✏️
              </button>
              <button
                class="text-wood-dark/40 hover:text-red-600 text-xs px-1"
                :aria-label="t('guestbook.deleteAria')"
                @click="startDelete(entry)"
              >
                🗑️
              </button>
            </template>
          </div>
        </div>
        <p
          class="text-wood-dark text-base whitespace-pre-wrap"
          :class="{ 'italic text-wood-dark/60': entry.secret }"
        >
          {{ entry.message }}
        </p>
      </div>
    </div>

    <!-- Infinite scroll sentinel -->
    <div
      v-if="nextCursor !== null"
      ref="sentinelRef"
      class="flex justify-center py-4"
    >
      <div v-if="loadingMore" class="pixel-spinner"></div>
    </div>

    <!-- Alert dialog (no input) -->
    <StardewDialog
      :visible="alertDialog.visible"
      :message="alertDialog.message"
      @close="alertDialog.visible = false"
    />

    <!-- Password prompt dialog -->
    <StardewDialog
      :visible="passwordDialog.visible"
      :message="passwordDialog.message"
      input-mode="password"
      @close="passwordDialog.visible = false"
      @cancel="passwordDialog.visible = false"
      @confirm="onPasswordConfirm"
    />

    <!-- Edit message dialog -->
    <StardewDialog
      :visible="editDialog.visible"
      :message="t('guestbook.editTitle')"
      input-mode="edit"
      :initial-value="editDialog.message"
      @close="editDialog.visible = false"
      @cancel="editDialog.visible = false"
      @confirm="onEditConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';
import { onBeforeRouteLeave } from 'vue-router';
import { useI18n } from 'vue-i18n';
import {
  getGuestbook,
  createGuestbookEntry,
  verifyGuestbookPassword,
  updateGuestbookEntry,
  deleteGuestbookEntry,
  ApiError,
} from '../services/api';
import { useHearts } from '../composables/useHearts';
import { useAchievements } from '../composables/useAchievements';
import StardewDialog from '../components/StardewDialog.vue';
import BoringAvatar from '../components/BoringAvatar.vue';
import TwEmoji from '../components/TwEmoji.vue';
import type { GuestbookEntry, DialogState, PendingAction } from '../types';

const { t, locale } = useI18n();
const { onContentUpdate, offContentUpdate } = useHearts();
const { award } = useAchievements();

const entries = ref<GuestbookEntry[]>([]);
const nextCursor = ref<number | null>(null);
const loadingMore = ref(false);
const sentinelRef = ref<HTMLElement | null>(null);
let scrollObserver: IntersectionObserver | null = null;
const nickname = ref('');
const message = ref('');
const password = ref('');
const secret = ref(false);
const loading = ref(true);
const submitting = ref(false);

const alertDialog = ref<DialogState>({ visible: false, message: '' });
const passwordDialog = ref<DialogState>({ visible: false, message: '' });
const editDialog = ref<DialogState>({ visible: false, message: '' });

const pendingAction = ref<PendingAction | null>(null);

function showAlert(msg: string) {
  alertDialog.value = { visible: true, message: msg };
}

async function loadEntries() {
  try {
    const data = await getGuestbook();
    entries.value = data.items;
    nextCursor.value = data.next_cursor;
  } catch (e) {
    console.error('Failed to load guestbook:', e);
  } finally {
    loading.value = false;
  }
}

async function loadMore() {
  if (loadingMore.value || nextCursor.value === null) return;
  loadingMore.value = true;
  try {
    const data = await getGuestbook(nextCursor.value);
    entries.value.push(...data.items);
    nextCursor.value = data.next_cursor;
  } catch (e) {
    console.error('Failed to load more:', e);
  } finally {
    loadingMore.value = false;
  }
}

async function submitEntry() {
  if (!nickname.value.trim() || !message.value.trim() || !password.value.trim())
    return;

  if (password.value.length < 5 || password.value.length > 30) {
    showAlert(t('common.passwordLength'));
    return;
  }

  submitting.value = true;
  try {
    const newEntry = await createGuestbookEntry(
      nickname.value,
      message.value,
      password.value,
      secret.value,
    );
    entries.value.unshift(newEntry);
    award('guestbook');
    message.value = '';
    password.value = '';
    secret.value = false;
  } catch {
    showAlert(t('guestbook.submitFailed'));
  } finally {
    submitting.value = false;
  }
}

function startEdit(entry: GuestbookEntry) {
  pendingAction.value = { type: 'edit', entry };
  passwordDialog.value = {
    visible: true,
    message: t('guestbook.enterPassword'),
  };
}

function startDelete(entry: GuestbookEntry) {
  pendingAction.value = { type: 'delete', entry };
  passwordDialog.value = {
    visible: true,
    message: t('common.deletePrompt'),
  };
}

async function onPasswordConfirm(pw: string) {
  passwordDialog.value.visible = false;

  if (!pw) return;

  const action = pendingAction.value;
  if (!action) return;

  if (action.type === 'delete') {
    try {
      await deleteGuestbookEntry(action.entry.id, pw);
      entries.value = entries.value.filter((e) => e.id !== action.entry.id);
    } catch (e) {
      showAlert(
        e instanceof ApiError && e.code === 'wrong_password'
          ? t('common.wrongPassword')
          : t('common.deleteFailed'),
      );
    }
    pendingAction.value = null;
  } else if (action.type === 'edit') {
    try {
      await verifyGuestbookPassword(action.entry.id, pw);
      pendingAction.value = { ...action, password: pw };
      editDialog.value = { visible: true, message: action.entry.message };
    } catch (e) {
      showAlert(
        e instanceof ApiError && e.code === 'wrong_password'
          ? t('common.wrongPassword')
          : t('guestbook.verifyFailed'),
      );
      pendingAction.value = null;
    }
  }
}

async function onEditConfirm(newMessage: string) {
  editDialog.value.visible = false;

  const action = pendingAction.value;
  if (!action || !newMessage?.trim()) {
    pendingAction.value = null;
    return;
  }

  try {
    const updated = await updateGuestbookEntry(
      action.entry.id,
      newMessage,
      action.password!,
    );
    const idx = entries.value.findIndex((e) => e.id === action.entry.id);
    if (idx !== -1) entries.value[idx] = updated;
  } catch (e) {
    showAlert(
      e instanceof ApiError && e.code === 'wrong_password'
        ? t('common.wrongPassword')
        : t('guestbook.updateFailed'),
    );
  }

  pendingAction.value = null;
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

const isDirty = computed(
  () => nickname.value.trim().length > 0 || message.value.trim().length > 0,
);

function onBeforeUnload(e: BeforeUnloadEvent) {
  if (isDirty.value) e.preventDefault();
}

onBeforeRouteLeave(() => {
  if (isDirty.value) {
    return window.confirm(t('guestbook.unsavedConfirm'));
  }
});

function setupScrollObserver() {
  scrollObserver?.disconnect();
  if (!sentinelRef.value) return;
  scrollObserver = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting) loadMore();
    },
    { rootMargin: '200px' },
  );
  scrollObserver.observe(sentinelRef.value);
}

watch(sentinelRef, () => setupScrollObserver());

onMounted(() => {
  loadEntries();
  onContentUpdate('guestbook', loadEntries);
  window.addEventListener('beforeunload', onBeforeUnload);
});

onUnmounted(() => {
  offContentUpdate('guestbook', loadEntries);
  window.removeEventListener('beforeunload', onBeforeUnload);
  scrollObserver?.disconnect();
});
</script>
