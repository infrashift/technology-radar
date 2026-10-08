<script setup lang="ts">
import { computed } from 'vue';
import type { PositionedBlip, Quadrant } from '../types/radar';
import { quadrantBounds, toViewBox } from '../config/radar';
import { blipHref } from '../utils/url';
import RadarBackdrop from './radar/RadarBackdrop.vue';
import RadarBlip from './radar/RadarBlip.vue';
import RadarLabel from './radar/RadarLabel.vue';

// One quadrant of the radar with the current blip emphasized and labelled.
// It has no interactivity, so pages render it without a client directive.
const props = defineProps<{
  blips: PositionedBlip[];
  quadrant: Quadrant;
  highlightedBlipId?: string;
}>();

const bounds = computed(() => quadrantBounds(props.quadrant));
const highlighted = computed(() => props.blips.find((b) => b.id === props.highlightedBlipId));
</script>

<template>
  <svg :viewBox="toViewBox(bounds)" class="quadrant-radar-svg" xmlns="http://www.w3.org/2000/svg" role="group" aria-label="Quadrant radar">
    <RadarBackdrop :quadrants="[quadrant]" />
    <RadarBlip
      v-for="blip in blips"
      :key="blip.id"
      :blip="blip"
      :href="blipHref(blip)"
      :emphasized="blip.id === highlightedBlipId"
    />
    <RadarLabel
      v-if="highlighted"
      :x="highlighted.x"
      :y="highlighted.y"
      :text="highlighted.title"
      :right-edge="bounds.x + bounds.size"
      :offset="22"
      bold
    />
  </svg>
</template>

<style scoped>
.quadrant-radar-svg {
  display: block;
  width: 100%;
  height: auto;
}
</style>
