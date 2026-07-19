<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-4">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">{{ t('game.title') }}</h1>
        <p class="text-wood-dark text-sm break-keep">
          {{ t('game.subtitle') }}
        </p>
      </div>
    </div>
    <div class="flex gap-2 mb-4">
      <router-link
        to="/hall-of-fame"
        class="pixel-btn text-sm py-2 no-underline flex-1 text-center"
      >
        <TwEmoji emoji="🏅" size="1rem" /> {{ t('hallOfFame.title') }}
      </router-link>
      <router-link
        to="/achievements"
        class="pixel-btn text-sm py-2 no-underline flex-1 text-center"
      >
        <TwEmoji emoji="🎯" size="1rem" /> {{ t('achievementsView.title') }}
      </router-link>
    </div>

    <!-- Top 3 rank banner -->
    <div v-if="topRank !== null" class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner px-6 py-4 text-center">
        <TwEmoji
          :emoji="topRank === 1 ? '🥇' : topRank === 2 ? '🥈' : '🥉'"
          size="2rem"
          class="mb-1"
        />
        <p class="text-primary text-lg font-bold">
          {{ t('game.rankAchieved', { rank: topRank }) }}
        </p>
        <p class="text-wood-dark text-sm mt-1">{{ t('game.congrats') }}</p>
      </div>
    </div>

    <!-- Idle State -->
    <div v-if="gameState === 'idle'" class="text-center">
      <!-- NPC hint for first-time players -->
      <div
        v-if="!showGallery && !isEmojiMode && npcText"
        class="wooden-panel p-2 mb-6"
      >
        <div class="wooden-panel-inner px-5 py-4">
          <div class="flex items-start gap-3">
            <TwEmoji emoji="🧑‍🌾" size="1.75rem" class="shrink-0 mt-0.5" />
            <p
              class="text-wood-dark text-sm text-left leading-relaxed break-keep"
            >
              {{ npcText }}
            </p>
          </div>
        </div>
      </div>

      <div class="parchment-bg p-6 mb-6">
        <p class="text-wood-dark text-base mb-4">
          {{ t('game.instruction', { seconds: CARD_GAME_TIMER / 1000 }) }}
        </p>
        <button @click="startGame" class="pixel-btn text-lg px-10 py-3">
          {{ t('game.start') }}
        </button>
      </div>

      <div v-if="rankings.length > 0" class="mb-6">
        <div class="wooden-panel p-2">
          <div class="wooden-panel-inner p-4">
            <h2 class="text-secondary text-center mb-3">
              {{ t('game.leaderboard') }}
            </h2>
            <div class="stardew-divider mb-3"></div>
            <div
              v-for="(score, index) in rankings"
              :key="score.id"
              class="flex justify-between items-center px-3 py-2 mb-1"
              :class="index % 2 === 0 ? 'bg-parchment-dark/50' : ''"
            >
              <span class="text-secondary text-sm flex items-center gap-2">
                <span class="text-primary">{{ index + 1 }}.</span>
                <BoringAvatar :name="score.nickname" :size="20" square />
                {{ score.nickname }}
              </span>
              <span class="text-primary text-sm"
                >{{ (score.time_ms / 1000).toFixed(2) }}s</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Photo Gallery (shown after score submission, only in photo mode) -->
      <div v-if="showGallery && !isEmojiMode" class="wooden-panel p-2">
        <div class="wooden-panel-inner p-3">
          <h2 class="text-secondary text-center text-sm mb-3">
            {{ t('game.viewAllPhotos') }}
          </h2>
          <div class="grid grid-cols-4 gap-1.5">
            <div
              v-for="(photo, idx) in gameDetailPhotos"
              :key="photo.id"
              @click="lightboxIndex = idx"
              class="aspect-square cursor-pointer overflow-hidden parchment-bg hover:scale-105 transition-transform"
            >
              <img
                :src="photo.thumbnail_url"
                :alt="String(photo.id)"
                class="w-full h-full object-cover"
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Playing State -->
    <div
      ref="playingArea"
      v-else-if="gameState === 'playing'"
      class="flex flex-col gap-2"
    >
      <div class="wooden-panel p-2 shrink-0">
        <div
          class="wooden-panel-inner px-4 py-2 flex justify-between items-center"
        >
          <span
            class="text-lg font-bold transition-colors duration-300"
            :class="
              timeLeft <= 5000
                ? 'timer-critical'
                : timeLeft <= 10000
                  ? 'timer-warning'
                  : 'text-primary'
            "
            >{{ (timeLeft / 1000).toFixed(2) }}s</span
          >
          <span class="text-secondary text-sm">{{
            t('game.pairCounter', { matched: matchedPairs, total: 6 })
          }}</span>
        </div>
      </div>

      <div class="grid grid-cols-3 gap-2 flex-1 min-h-0">
        <div
          v-for="(card, index) in cards"
          :key="index"
          @click="flipCard(index)"
          class="aspect-square cursor-pointer transition-transform hover:scale-105 overflow-hidden"
          :class="{ 'pointer-events-none': card.flipped || card.matched }"
        >
          <!-- Card Front (revealed) -->
          <div
            v-if="card.flipped || card.matched"
            class="w-full h-full parchment-bg overflow-hidden flex items-center justify-center"
            :class="card.matched ? 'border-accent border-3' : ''"
          >
            <img
              v-if="!card.emoji"
              :src="card.photo"
              alt="card"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-4xl">{{ card.emoji }}</span>
          </div>
          <!-- Card Back -->
          <div
            v-else
            class="w-full h-full card-back flex items-center justify-center"
          >
            <span class="text-parchment text-2xl">?</span>
          </div>
        </div>
      </div>

      <!-- Match popup -->
      <div
        v-if="showMatchPopup"
        class="fixed inset-0 flex items-center justify-center z-40 pointer-events-none"
      >
        <div class="match-popup pointer-events-none">{{ t('game.match') }}</div>
      </div>
    </div>

    <!-- Won State -->
    <div v-else-if="gameState === 'won'" class="text-center">
      <div class="wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-6 py-6">
          <h2 class="text-2xl text-primary mb-3">{{ t('game.won') }}</h2>
          <p class="text-wood-dark mb-4">
            {{
              t('game.elapsed', {
                seconds: ((CARD_GAME_TIMER - finalTime) / 1000).toFixed(2),
              })
            }}
          </p>

          <div class="mb-4">
            <label class="block text-secondary text-sm mb-2">{{
              t('game.initialsLabel')
            }}</label>
            <input
              v-model="playerName"
              maxlength="3"
              class="pixel-input w-24 text-center text-xl uppercase"
              placeholder="AAA"
              @input="
                playerName = playerName.toUpperCase().replace(/[^A-Z]/g, '')
              "
            />
          </div>

          <div class="flex gap-3 justify-center">
            <button
              @click="submitScore"
              :disabled="playerName.length !== 3 || submittingScore"
              class="pixel-btn"
            >
              {{ submittingScore ? t('game.saving') : t('game.saveScore') }}
            </button>
            <button @click="resetGame" class="pixel-btn-outline bg-parchment">
              {{ t('game.playAgain') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Photo Gallery (only in photo mode) -->
      <div v-if="!isEmojiMode" class="wooden-panel p-2">
        <div class="wooden-panel-inner p-3">
          <h2 class="text-secondary text-center text-sm mb-3">
            {{ t('game.viewAllPhotos') }}
          </h2>
          <div class="grid grid-cols-4 gap-1.5">
            <div
              v-for="(photo, idx) in gameDetailPhotos"
              :key="photo.id"
              @click="lightboxIndex = idx"
              class="aspect-square cursor-pointer overflow-hidden parchment-bg hover:scale-105 transition-transform"
            >
              <img
                :src="photo.thumbnail_url"
                :alt="String(photo.id)"
                class="w-full h-full object-cover"
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Lost State -->
    <div v-else-if="gameState === 'lost'" class="text-center">
      <div class="wooden-panel p-2 mb-6">
        <div class="wooden-panel-inner px-6 py-6">
          <h2 class="text-2xl text-secondary mb-3">{{ t('game.timeUp') }}</h2>
          <p class="text-wood-dark mb-4">
            {{ t('game.pairsMatched', { matched: matchedPairs, total: 6 }) }}
          </p>
          <button @click="resetGame" class="pixel-btn text-lg px-10 py-3">
            {{ t('game.retry') }}
          </button>
        </div>
      </div>

      <!-- Photo Gallery (only in photo mode) -->
      <div v-if="!isEmojiMode" class="wooden-panel p-2">
        <div class="wooden-panel-inner p-3">
          <h2 class="text-secondary text-center text-sm mb-3">
            {{ t('game.viewAllPhotos') }}
          </h2>
          <div class="grid grid-cols-4 gap-1.5">
            <div
              v-for="(photo, idx) in gameDetailPhotos"
              :key="photo.id"
              @click="lightboxIndex = idx"
              class="aspect-square cursor-pointer overflow-hidden parchment-bg hover:scale-105 transition-transform"
            >
              <img
                :src="photo.thumbnail_url"
                :alt="String(photo.id)"
                class="w-full h-full object-cover"
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <LightboxViewer
      v-if="!isEmojiMode"
      v-model="lightboxIndex"
      :items="lightboxItems"
    />

    <StardewDialog
      :visible="dialogVisible"
      :message="dialogMessage"
      @close="dialogVisible = false"
    />

    <!-- Confetti -->
    <Teleport to="body">
      <div
        v-if="showConfetti"
        class="fixed inset-0 pointer-events-none z-[200] overflow-hidden"
      >
        <span
          v-for="i in 40"
          :key="i"
          class="confetti-piece"
          :style="{
            left: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 2}s`,
            animationDuration: `${2.5 + Math.random() * 2}s`,
            fontSize: `${14 + Math.random() * 14}px`,
          }"
          >{{ ['🎊', '🎉', '✨', '⭐', '💫', '🌟'][i % 6] }}</span
        >
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  getGamePhotos,
  getGameRankings,
  recordGameBeat,
  submitGameScore,
} from '../services/api';
import StardewDialog from '../components/StardewDialog.vue';
import LightboxViewer from '../components/LightboxViewer.vue';
import BoringAvatar from '../components/BoringAvatar.vue';
import { useAchievements } from '../composables/useAchievements';
import {
  getGameCompletion,
  setGameCompletion,
  clearGameCompletion,
  clearAll as clearAllStore,
} from '../lib/store';
import TwEmoji from '../components/TwEmoji.vue';
import { CARD_GAME_TIMER, GAME_NPC_MESSAGE } from '../config/wedding';
import type { GameState, Card, GameScore, GamePhotoEntry } from '../types';

const { t } = useI18n();
const { award: awardAchievement, clearAll: clearAchievements } =
  useAchievements();

const playingArea = ref<HTMLElement | null>(null);

const lightboxIndex = ref<number | null>(null);
const viewedPhotos = new Set<number>();

watch(lightboxIndex, (idx) => {
  if (idx !== null && gameDetailPhotos.value.length > 0) {
    viewedPhotos.add(idx);
    if (viewedPhotos.size >= gameDetailPhotos.value.length) {
      awardAchievement('gallery_all');
    }
  }
});

const isEmojiMode = ref(false);
const gameDetailPhotos = ref<GamePhotoEntry[]>([]);

const lightboxItems = computed(() =>
  gameDetailPhotos.value.map((p) => ({
    thumbnail: p.thumbnail_url,
    src: p.detail_url,
  })),
);

const dialogVisible = ref(false);
const dialogMessage = ref('');

function showDialog(msg: string) {
  dialogMessage.value = msg;
  dialogVisible.value = true;
}

const gameState = ref<GameState>('idle');
const cards = ref<Card[]>([]);
const timeLeft = ref(CARD_GAME_TIMER);
const matchedPairs = ref(0);
const finalTime = ref(0);
const playerName = ref('');
const rankings = ref<GameScore[]>([]);
const submittingScore = ref(false);
const showMatchPopup = ref(false);
const showGallery = ref(false);
const topRank = ref<number | null>(null);
const showConfetti = ref(false);
const gameToken = ref('');

let timer: ReturnType<typeof setInterval> | null = null;
let flippedCards: number[] = [];
let pendingMismatch: number[] = [];

const NPC_MESSAGE = GAME_NPC_MESSAGE;
const npcText = ref('');
let npcTimer: ReturnType<typeof setInterval> | null = null;

function startNpcTyping() {
  if (showGallery.value || isEmojiMode.value) return;
  let i = 0;
  npcTimer = setInterval(() => {
    if (i < NPC_MESSAGE.length) {
      npcText.value = NPC_MESSAGE.slice(0, ++i);
    } else {
      if (npcTimer) clearInterval(npcTimer);
      npcTimer = null;
    }
  }, 50);
}

async function loadRankings() {
  try {
    const data = await getGameRankings();
    rankings.value = data.rankings;
    const completedAt = getGameCompletion();
    const resetAt = Number(data.game_reset_at || '0');
    let gameCompleted = completedAt !== null && completedAt > resetAt;
    if (completedAt !== null && !gameCompleted) {
      clearAllStore();
      clearAchievements();
    }
    if (gameCompleted) {
      showGallery.value = true;
      loadGamePhotos();
    } else {
      startNpcTyping();
    }
  } catch (e) {
    console.error('Failed to load rankings:', e);
  }
}

async function loadGamePhotos() {
  if (gameDetailPhotos.value.length > 0) return;
  try {
    const data = await getGamePhotos();
    if (data.type === 'photo') {
      isEmojiMode.value = false;
      gameDetailPhotos.value = data.photos;
    }
  } catch {
    // ignore — gallery just won't show
  }
}

function preloadImages(urls: string[]): Promise<Event[]> {
  return Promise.all(
    urls.map(
      (url) =>
        new Promise<Event>((resolve, reject) => {
          const img = new Image();
          img.onload = resolve;
          img.onerror = reject;
          img.src = url;
        }),
    ),
  );
}

async function startGame() {
  try {
    const data = await getGamePhotos();
    gameToken.value = data.game_token || '';

    let shuffled: Card[];

    if (data.type === 'photo') {
      isEmojiMode.value = false;
      gameDetailPhotos.value = data.photos;

      const rng = () => Math.random() - 0.5;
      const selected = [...data.photos].sort(rng).slice(0, 6);

      await preloadImages(selected.map((p) => p.thumbnail_url));

      shuffled = selected
        .flatMap((photo, i) => [
          {
            photo: photo.thumbnail_url,
            id: i * 2,
            pairId: String(photo.id),
            flipped: false,
            matched: false,
          },
          {
            photo: photo.thumbnail_url,
            id: i * 2 + 1,
            pairId: String(photo.id),
            flipped: false,
            matched: false,
          },
        ])
        .sort(() => Math.random() - 0.5);
    } else {
      isEmojiMode.value = true;
      gameDetailPhotos.value = [];

      const rng = () => Math.random() - 0.5;
      const selected = [...data.photos].sort(rng).slice(0, 6);

      shuffled = selected
        .flatMap((entry, i) => [
          {
            photo: '',
            emoji: entry.emoji,
            id: i * 2,
            pairId: entry.id,
            flipped: false,
            matched: false,
          },
          {
            photo: '',
            emoji: entry.emoji,
            id: i * 2 + 1,
            pairId: entry.id,
            flipped: false,
            matched: false,
          },
        ])
        .sort(rng);
    }

    cards.value = shuffled;
    timeLeft.value = CARD_GAME_TIMER;
    matchedPairs.value = 0;
    flippedCards = [];
    gameState.value = 'playing';
    await nextTick();
    playingArea.value?.scrollIntoView({ behavior: 'smooth', block: 'start' });

    timer = setInterval(() => {
      timeLeft.value -= 10;
      if (timeLeft.value <= 0) {
        endGame(false);
      }
    }, 10);
  } catch {
    showDialog(t('game.photoLoadFailed'));
  }
}

function flipCard(index: number) {
  const card = cards.value[index];
  if (card.flipped || card.matched) return;

  if (pendingMismatch.length === 2) {
    cards.value[pendingMismatch[0]].flipped = false;
    cards.value[pendingMismatch[1]].flipped = false;
    pendingMismatch = [];
  }

  card.flipped = true;
  flippedCards.push(index);
  navigator.vibrate?.(10);

  if (flippedCards.length === 2) {
    const [first, second] = flippedCards;
    const firstCard = cards.value[first];
    const secondCard = cards.value[second];
    flippedCards = [];

    if (firstCard.pairId === secondCard.pairId) {
      firstCard.matched = true;
      secondCard.matched = true;
      matchedPairs.value++;
      navigator.vibrate?.([20, 40, 20]);

      showMatchPopup.value = true;
      setTimeout(() => {
        showMatchPopup.value = false;
      }, 600);

      if (matchedPairs.value === 6) {
        endGame(true);
      }
    } else {
      pendingMismatch = [first, second];
      setTimeout(() => {
        if (pendingMismatch[0] === first && pendingMismatch[1] === second) {
          firstCard.flipped = false;
          secondCard.flipped = false;
          pendingMismatch = [];
        }
      }, 500);
    }
  }
}

function endGame(won: boolean) {
  if (timer) clearInterval(timer);
  timer = null;

  if (won) {
    finalTime.value = timeLeft.value;
    gameState.value = 'won';
    setGameCompletion();
    recordGameBeat();

    awardAchievement('gallery');

    const timeMs = CARD_GAME_TIMER - timeLeft.value;
    const rank = rankings.value.filter((s) => s.time_ms < timeMs).length;
    if (rank < 3) {
      topRank.value = rank + 1;
      showConfetti.value = true;
      setTimeout(() => {
        showConfetti.value = false;
      }, 4000);
      const medals = ['gold', 'silver', 'bronze'];
      for (let i = rank; i < 3; i++) {
        awardAchievement(medals[i]);
      }
    }

    if (timeMs <= 10000) {
      awardAchievement('speed');
    }
  } else {
    gameState.value = 'lost';
  }
}

async function submitScore() {
  if (playerName.value.length !== 3) return;

  submittingScore.value = true;
  try {
    const timeMs = CARD_GAME_TIMER - finalTime.value;
    await submitGameScore(playerName.value, timeMs, gameToken.value);
    await loadRankings();
    topRank.value = null;
    gameState.value = 'idle';
    playerName.value = '';
    showGallery.value = true;
  } catch {
    showDialog(t('game.scoreSaveFailed'));
  } finally {
    submittingScore.value = false;
  }
}

function resetGame() {
  gameState.value = 'idle';
  playerName.value = '';
  showGallery.value = false;
  topRank.value = null;
  showConfetti.value = false;
}

onMounted(loadRankings);

onUnmounted(() => {
  if (timer) clearInterval(timer);
  if (npcTimer) clearInterval(npcTimer);
});
</script>

<style scoped>
.card-back {
  background-color: var(--btn-bg);
  border: 3px solid var(--btn-border);
  box-shadow:
    inset 2px 2px 0 var(--btn-highlight),
    inset -2px -2px 0 var(--btn-shadow);
  background-image: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 4px,
    rgba(92, 58, 14, 0.15) 4px,
    rgba(92, 58, 14, 0.15) 8px
  );
}

.border-3 {
  border-width: 3px;
}

.match-popup {
  font-size: 3rem;
  font-weight: bold;
  color: #ffeb3b;
  text-shadow:
    0 0 10px rgba(255, 235, 59, 0.8),
    0 0 20px rgba(255, 165, 0, 0.6),
    0 0 40px rgba(255, 100, 0, 0.4),
    3px 3px 0 #5c3a0e,
    -1px -1px 0 #5c3a0e;
  background: rgba(92, 58, 14, 0.85);
  border: 4px solid #c4943a;
  box-shadow:
    0 0 20px rgba(255, 235, 59, 0.5),
    inset 0 0 10px rgba(255, 235, 59, 0.2);
  padding: 0.5rem 2rem;
  animation: match-pop 0.6s ease-out;
}

@keyframes match-pop {
  0% {
    transform: scale(0.3);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.timer-warning {
  color: #e65100;
  animation: timer-pulse 1s ease-in-out infinite;
}

.timer-critical {
  color: #d32f2f;
  animation: timer-shake 0.4s ease-in-out infinite;
  text-shadow: 0 0 8px rgba(211, 47, 47, 0.5);
}

@keyframes timer-pulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.15);
    opacity: 0.8;
  }
}

@keyframes timer-shake {
  0%,
  100% {
    transform: translateX(0) scale(1.1);
  }
  25% {
    transform: translateX(-2px) scale(1.15);
  }
  75% {
    transform: translateX(2px) scale(1.15);
  }
}

.confetti-piece {
  position: absolute;
  top: -30px;
  animation: confetti-fall linear forwards;
}

@keyframes confetti-fall {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    transform: translateY(100vh) rotate(720deg);
    opacity: 0;
  }
}
</style>
