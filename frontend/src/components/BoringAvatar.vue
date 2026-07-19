<template>
  <svg
    :viewBox="`0 0 ${SIZE} ${SIZE}`"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    :width="size"
    :height="size"
  >
    <mask
      :id="maskId"
      maskUnits="userSpaceOnUse"
      x="0"
      y="0"
      :width="SIZE"
      :height="SIZE"
    >
      <rect
        :width="SIZE"
        :height="SIZE"
        :rx="square ? undefined : SIZE * 2"
        fill="#FFFFFF"
      />
    </mask>
    <g :mask="`url(#${maskId})`">
      <rect :width="SIZE" :height="SIZE" :fill="data.backgroundColor" />
      <rect
        x="0"
        y="0"
        :width="SIZE"
        :height="SIZE"
        :transform="`translate(${data.wrapperTranslateX} ${data.wrapperTranslateY}) rotate(${data.wrapperRotate} ${SIZE / 2} ${SIZE / 2}) scale(${data.wrapperScale})`"
        :fill="data.wrapperColor"
        :rx="data.isCircle ? SIZE : SIZE / 6"
      />
      <g
        :transform="`translate(${data.faceTranslateX} ${data.faceTranslateY}) rotate(${data.faceRotate} ${SIZE / 2} ${SIZE / 2})`"
      >
        <path
          v-if="data.isMouthOpen"
          :d="`M15 ${19 + data.mouthSpread}c2 1 4 1 6 0`"
          :stroke="data.faceColor"
          fill="none"
          stroke-linecap="round"
        />
        <path
          v-else
          :d="`M13,${19 + data.mouthSpread} a1,0.75 0 0,0 10,0`"
          :fill="data.faceColor"
        />
        <rect
          :x="14 - data.eyeSpread"
          y="14"
          width="1.5"
          height="2"
          rx="1"
          stroke="none"
          :fill="data.faceColor"
        />
        <rect
          :x="20 + data.eyeSpread"
          y="14"
          width="1.5"
          height="2"
          rx="1"
          stroke="none"
          :fill="data.faceColor"
        />
      </g>
    </g>
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { AVATAR_COLORS } from '../config/wedding';

const SIZE = 36;

interface AvatarData {
  wrapperColor: string;
  faceColor: string;
  backgroundColor: string;
  wrapperTranslateX: number;
  wrapperTranslateY: number;
  wrapperRotate: number;
  wrapperScale: number;
  isMouthOpen: boolean;
  isCircle: boolean;
  eyeSpread: number;
  mouthSpread: number;
  faceRotate: number;
  faceTranslateX: number;
  faceTranslateY: number;
}

const props = withDefaults(
  defineProps<{
    name: string;
    size?: number;
    square?: boolean;
    colors?: string[];
  }>(),
  {
    size: 28,
    square: false,
    colors: () => AVATAR_COLORS.split(',').map((c) => `#${c}`),
  },
);

const maskId = computed(() => `beam-${hashCode(props.name)}`);

const data = computed<AvatarData>(() => generateData(props.name, props.colors));

function hashCode(name: string): number {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = (hash << 5) - hash + name.charCodeAt(i);
    hash = hash & hash;
  }
  return Math.abs(hash);
}

function getDigit(number: number, ntn: number): number {
  return Math.floor((number / Math.pow(10, ntn)) % 10);
}

function getBoolean(number: number, ntn: number): boolean {
  return !(getDigit(number, ntn) % 2);
}

function getUnit(number: number, range: number, index?: number): number {
  const value = number % range;
  if (index && getDigit(number, index) % 2 === 0) return -value;
  return value;
}

function getRandomColor(
  number: number,
  colors: string[],
  range: number,
): string {
  return colors[number % range];
}

function getContrast(hexcolor: string): string {
  const hex = hexcolor.startsWith('#') ? hexcolor.slice(1) : hexcolor;
  const r = parseInt(hex.substr(0, 2), 16);
  const g = parseInt(hex.substr(2, 2), 16);
  const b = parseInt(hex.substr(4, 2), 16);
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  return yiq >= 128 ? '#000000' : '#FFFFFF';
}

function generateData(name: string, colors: string[]): AvatarData {
  const numFromName = hashCode(name);
  const range = colors.length;
  const wrapperColor = getRandomColor(numFromName, colors, range);
  const preTranslateX = getUnit(numFromName, 10, 1);
  const wrapperTranslateX =
    preTranslateX < 5 ? preTranslateX + SIZE / 9 : preTranslateX;
  const preTranslateY = getUnit(numFromName, 10, 2);
  const wrapperTranslateY =
    preTranslateY < 5 ? preTranslateY + SIZE / 9 : preTranslateY;

  return {
    wrapperColor,
    faceColor: getContrast(wrapperColor),
    backgroundColor: getRandomColor(numFromName + 13, colors, range),
    wrapperTranslateX,
    wrapperTranslateY,
    wrapperRotate: getUnit(numFromName, 360),
    wrapperScale: 1 + getUnit(numFromName, SIZE / 12) / 10,
    isMouthOpen: getBoolean(numFromName, 2),
    isCircle: getBoolean(numFromName, 1),
    eyeSpread: getUnit(numFromName, 5),
    mouthSpread: getUnit(numFromName, 3),
    faceRotate: getUnit(numFromName, 10, 3),
    faceTranslateX:
      wrapperTranslateX > SIZE / 6
        ? wrapperTranslateX / 2
        : getUnit(numFromName, 8, 1),
    faceTranslateY:
      wrapperTranslateY > SIZE / 6
        ? wrapperTranslateY / 2
        : getUnit(numFromName, 7, 2),
  };
}
</script>
