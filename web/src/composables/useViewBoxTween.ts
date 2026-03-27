import { ref } from 'vue';

function parseViewBox(vb: string): [number, number, number, number] {
  return vb.split(/\s+/).map(Number) as [number, number, number, number];
}

function formatViewBox(vb: [number, number, number, number]): string {
  return vb.map((v) => Math.round(v * 100) / 100).join(' ');
}

function easeOutCubic(t: number): number {
  return 1 - Math.pow(1 - t, 3);
}

export function useViewBoxTween(defaultViewBox: string) {
  const currentViewBox = ref(defaultViewBox);
  let animationId: number | null = null;

  function tweenTo(target: string, duration: number = 400) {
    if (animationId !== null) {
      cancelAnimationFrame(animationId);
    }

    const from = parseViewBox(currentViewBox.value);
    const to = parseViewBox(target);
    const startTime = performance.now();

    function animate(now: number) {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const eased = easeOutCubic(progress);

      const current: [number, number, number, number] = [
        from[0] + (to[0] - from[0]) * eased,
        from[1] + (to[1] - from[1]) * eased,
        from[2] + (to[2] - from[2]) * eased,
        from[3] + (to[3] - from[3]) * eased,
      ];

      currentViewBox.value = formatViewBox(current);

      if (progress < 1) {
        animationId = requestAnimationFrame(animate);
      } else {
        animationId = null;
      }
    }

    animationId = requestAnimationFrame(animate);
  }

  return { currentViewBox, tweenTo };
}
