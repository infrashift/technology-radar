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

function prefersReducedMotion(): boolean {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
}

export function useViewBoxTween(defaultViewBox: string) {
  const currentViewBox = ref(defaultViewBox);
  let animationId: number | null = null;

  /** Animate to target; a duration of 0 (or reduced motion) jumps immediately. */
  function tweenTo(target: string, duration: number = 400) {
    if (animationId !== null) {
      cancelAnimationFrame(animationId);
      animationId = null;
    }

    if (duration <= 0 || prefersReducedMotion()) {
      currentViewBox.value = target;
      return;
    }

    const from = parseViewBox(currentViewBox.value);
    const to = parseViewBox(target);
    const startTime = performance.now();

    function animate(now: number) {
      const progress = Math.min((now - startTime) / duration, 1);
      const eased = easeOutCubic(progress);
      currentViewBox.value = formatViewBox(from.map((f, i) => f + (to[i] - f) * eased) as [number, number, number, number]);
      animationId = progress < 1 ? requestAnimationFrame(animate) : null;
    }

    animationId = requestAnimationFrame(animate);
  }

  return { currentViewBox, tweenTo };
}
