<script setup lang="ts">
import { computed } from 'vue';

// A title label beside a blip. It flips to the blip's left when it would run
// past the right edge of the visible area.
const props = withDefaults(
  defineProps<{
    x: number;
    y: number;
    text: string;
    rightEdge: number;
    offset?: number;
    bold?: boolean;
  }>(),
  { offset: 16, bold: false },
);

const CHAR_WIDTH = 6.6; // approximate width of an 11px Inter glyph
const PADDING = 8;

const width = computed(() => props.text.length * CHAR_WIDTH + PADDING * 2);
const flip = computed(() => props.x + props.offset + width.value > props.rightEdge);
const rectX = computed(() => (flip.value ? props.x - props.offset - width.value : props.x + props.offset));
</script>

<template>
  <g class="radar-label" pointer-events="none">
    <rect :x="rectX" :y="y - 12" :width="width" height="24" rx="4" class="label-bg" fill-opacity="0.92" />
    <text
      :x="rectX + PADDING"
      :y="y + 1"
      class="label-text"
      font-size="11"
      :font-weight="bold ? 600 : 400"
      dominant-baseline="middle"
    >
      {{ text }}
    </text>
  </g>
</template>

<style scoped>
.label-bg {
  fill: var(--color-radar-tooltip-bg);
}

.label-text {
  fill: var(--color-radar-tooltip-text);
}
</style>
