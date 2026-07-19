<script setup lang="ts">
import { computed, ref } from 'vue';

interface DailyCount {
  date: string;
  count: number;
}

const props = defineProps<{ data: DailyCount[] }>();

const W = 600;
const H = 200;
const PAD = { top: 20, right: 20, bottom: 30, left: 45 };

const chartW = W - PAD.left - PAD.right;
const chartH = H - PAD.top - PAD.bottom;

const maxCount = computed(() => Math.max(...props.data.map((d) => d.count), 1));

const points = computed(() =>
  props.data.map((d, i) => ({
    x: PAD.left + (i / Math.max(props.data.length - 1, 1)) * chartW,
    y: PAD.top + chartH - (d.count / maxCount.value) * chartH,
    date: d.date,
    count: d.count,
  })),
);

const linePath = computed(() =>
  points.value.map((p, i) => `${i === 0 ? 'M' : 'L'}${p.x},${p.y}`).join(' '),
);

const areaPath = computed(() => {
  if (!points.value.length) return '';
  const first = points.value[0];
  const last = points.value[points.value.length - 1];
  return `${linePath.value} L${last.x},${PAD.top + chartH} L${first.x},${PAD.top + chartH} Z`;
});

const yTicks = computed(() => {
  const max = maxCount.value;
  const step = Math.max(1, Math.ceil(max / 4));
  const ticks = [];
  for (let v = 0; v <= max; v += step) {
    ticks.push({
      value: v,
      y: PAD.top + chartH - (v / max) * chartH,
    });
  }
  return ticks;
});

const xLabels = computed(() => {
  const pts = points.value;
  if (pts.length <= 7) return pts;
  const step = Math.max(1, Math.floor(pts.length / 6));
  return pts.filter((_, i) => i % step === 0 || i === pts.length - 1);
});

const hovered = ref<number | null>(null);
const svgRef = ref<SVGSVGElement | null>(null);

function updateHovered(clientX: number) {
  if (!svgRef.value || !points.value.length) return;
  const rect = svgRef.value.getBoundingClientRect();
  const svgX = ((clientX - rect.left) / rect.width) * W;
  if (svgX < PAD.left || svgX > W - PAD.right) {
    hovered.value = null;
    return;
  }
  let closest = 0;
  let minDist = Infinity;
  for (let i = 0; i < points.value.length; i++) {
    const dist = Math.abs(points.value[i].x - svgX);
    if (dist < minDist) {
      minDist = dist;
      closest = i;
    }
  }
  if (hovered.value !== closest) {
    navigator.vibrate?.(5);
  }
  hovered.value = closest;
}

function onMouseMove(e: MouseEvent) {
  updateHovered(e.clientX);
}

function onMouseLeave() {
  hovered.value = null;
}

function onTouchStart(e: TouchEvent) {
  e.preventDefault();
  updateHovered(e.touches[0].clientX);
}

function onTouchMove(e: TouchEvent) {
  e.preventDefault();
  updateHovered(e.touches[0].clientX);
}

function onTouchEnd() {
  hovered.value = null;
}

function formatDate(dateStr: string) {
  const [, m, d] = dateStr.split('-');
  return `${parseInt(m)}/${parseInt(d)}`;
}

const tooltipWidth = computed(() => {
  if (hovered.value === null) return 0;
  return String(points.value[hovered.value].count).length * 6 + 16;
});

const flipBelow = computed(() => {
  if (hovered.value === null) return false;
  return points.value[hovered.value].y < PAD.top + 30;
});
</script>

<template>
  <svg
    ref="svgRef"
    :viewBox="`0 0 ${W} ${H}`"
    class="w-full h-auto"
    preserveAspectRatio="xMidYMid meet"
    style="user-select: none; cursor: crosshair"
    @mousemove="onMouseMove"
    @mouseleave="onMouseLeave"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
  >
    <!-- Grid lines -->
    <line
      v-for="tick in yTicks"
      :key="tick.value"
      :x1="PAD.left"
      :y1="tick.y"
      :x2="PAD.left + chartW"
      :y2="tick.y"
      stroke="currentColor"
      stroke-opacity="0.1"
      stroke-dasharray="4 4"
    />

    <!-- Y-axis labels -->
    <text
      v-for="tick in yTicks"
      :key="'t' + tick.value"
      :x="PAD.left - 8"
      :y="tick.y + 4"
      text-anchor="end"
      fill="currentColor"
      fill-opacity="0.5"
      font-size="10"
    >
      {{ tick.value }}
    </text>

    <!-- Area fill -->
    <path
      v-if="areaPath"
      :d="areaPath"
      fill="hsl(var(--primary))"
      fill-opacity="0.1"
    />

    <!-- Line -->
    <path
      v-if="linePath"
      :d="linePath"
      fill="none"
      stroke="hsl(var(--primary))"
      stroke-width="2"
      stroke-linejoin="round"
    />

    <!-- Hover crosshair + dot + tooltip -->
    <template v-if="hovered !== null">
      <!-- Vertical line -->
      <line
        :x1="points[hovered].x"
        :y1="PAD.top"
        :x2="points[hovered].x"
        :y2="PAD.top + chartH"
        stroke="hsl(var(--primary))"
        stroke-opacity="0.3"
        stroke-width="1"
        stroke-dasharray="3 3"
      />

      <!-- Dot on the line -->
      <circle
        :cx="points[hovered].x"
        :cy="points[hovered].y"
        r="4"
        fill="hsl(var(--primary))"
      />

      <!-- Tooltip (flips below point when near top edge) -->
      <g :transform="`translate(${
        points[hovered].x + tooltipWidth / 2 + 4 > W - PAD.right
          ? points[hovered].x - tooltipWidth / 2 - 4
          : points[hovered].x
      }, ${flipBelow
        ? points[hovered].y + 14
        : points[hovered].y - 14
      })`">
        <rect
          :x="-tooltipWidth / 2"
          :y="flipBelow ? 0 : -20"
          :width="tooltipWidth"
          height="20"
          rx="4"
          fill="hsl(var(--popover))"
          stroke="hsl(var(--border))"
          stroke-width="0.5"
        />
        <text
          text-anchor="middle"
          :y="flipBelow ? 13 : -7"
          font-size="9"
          font-weight="600"
          fill="hsl(var(--popover-foreground))"
        >
          {{ points[hovered].count }}
        </text>
      </g>
    </template>

    <!-- X-axis labels -->
    <text
      v-for="label in xLabels"
      :key="'x' + label.date"
      :x="label.x"
      :y="PAD.top + chartH + 18"
      text-anchor="middle"
      fill="currentColor"
      :fill-opacity="hovered !== null ? 0.2 : 0.5"
      font-size="10"
    >
      {{ formatDate(label.date) }}
    </text>

    <!-- Highlighted date on hover -->
    <text
      v-if="hovered !== null"
      :x="points[hovered].x"
      :y="PAD.top + chartH + 18"
      text-anchor="middle"
      fill="hsl(var(--primary))"
      font-size="10"
      font-weight="600"
    >
      {{ formatDate(points[hovered].date) }}
    </text>
  </svg>
</template>
