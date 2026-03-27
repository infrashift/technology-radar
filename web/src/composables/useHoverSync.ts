import { ref, computed } from 'vue';

export function useHoverSync() {
  const hoveredBlipId = ref<string | null>(null);

  function setHover(id: string | null) {
    hoveredBlipId.value = id;
  }

  function isHighlighted(id: string) {
    return computed(() => hoveredBlipId.value === id);
  }

  return { hoveredBlipId, setHover, isHighlighted };
}
