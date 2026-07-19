import { describe, it, expect } from 'vitest';
import { runWithConcurrency } from '../upload';

describe('runWithConcurrency', () => {
  it('processes every item', async () => {
    const done: number[] = [];
    await runWithConcurrency([1, 2, 3, 4, 5], 2, async (n) => {
      done.push(n);
    });
    expect(done.sort()).toEqual([1, 2, 3, 4, 5]);
  });

  it('never exceeds the concurrency limit', async () => {
    let active = 0;
    let peak = 0;
    await runWithConcurrency(Array.from({ length: 10 }, (_, i) => i), 3, () => {
      active++;
      peak = Math.max(peak, active);
      return new Promise((resolve) =>
        setTimeout(() => {
          active--;
          resolve();
        }, 5),
      );
    });
    expect(peak).toBeLessThanOrEqual(3);
    expect(peak).toBeGreaterThan(1);
  });
});
