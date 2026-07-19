<template>
  <img
    v-if="src"
    :src="src"
    :alt="emoji"
    class="twemoji"
    :style="{ width: size, height: size }"
    draggable="false"
  />
  <span v-else>{{ emoji }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const svgs = import.meta.glob('../assets/emoji/*.svg', {
  eager: true,
  query: '?url',
  import: 'default',
}) as Record<string, string>;

function emojiToCodepoints(emoji: string): string {
  const codepoints: string[] = [];
  for (const char of emoji) {
    const cp = char.codePointAt(0);
    if (cp !== undefined && cp !== 0xfe0f) {
      codepoints.push(cp.toString(16));
    }
  }
  return codepoints.join('-');
}

const props = defineProps<{
  emoji: string;
  size?: string;
}>();

const src = computed(() => {
  const code = emojiToCodepoints(props.emoji);
  const key = `../assets/emoji/${code}.svg`;
  return svgs[key] || null;
});
</script>

<style>
.twemoji {
  display: inline-block;
  vertical-align: -0.1em;
}
</style>
