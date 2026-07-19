import { describe, it, expect } from 'vitest';
import { ref } from 'vue';
import { useSort } from '../composables/useSort';

interface TestItem {
  name: string;
  age: number;
  active: boolean;
  score: number | null;
}

function makeItems(): TestItem[] {
  return [
    { name: 'Charlie', age: 30, active: true, score: 80 },
    { name: 'Alice', age: 25, active: false, score: null },
    { name: 'Bob', age: 35, active: true, score: 95 },
  ];
}

describe('useSort', () => {
  it('sorts by string field ascending', () => {
    const items = ref(makeItems());
    const { sortedItems, toggleSort } = useSort(items);
    toggleSort('name');
    expect(sortedItems.value.map((i) => i.name)).toEqual([
      'Alice',
      'Bob',
      'Charlie',
    ]);
  });

  it('sorts by string field descending', () => {
    const items = ref(makeItems());
    const { sortedItems, toggleSort } = useSort(items);
    toggleSort('name');
    toggleSort('name');
    expect(sortedItems.value.map((i) => i.name)).toEqual([
      'Charlie',
      'Bob',
      'Alice',
    ]);
  });

  it('resets to default on third toggle', () => {
    const items = ref(makeItems());
    const { sortedItems, sortKey, toggleSort } = useSort(items);
    toggleSort('name');
    toggleSort('name');
    toggleSort('name');
    expect(sortKey.value).toBeNull();
    expect(sortedItems.value.map((i) => i.name)).toEqual([
      'Charlie',
      'Alice',
      'Bob',
    ]);
  });

  it('switches field resets to ascending', () => {
    const items = ref(makeItems());
    const { sortedItems, sortKey, sortDir, toggleSort } = useSort(items);
    toggleSort('name');
    toggleSort('name');
    expect(sortDir.value).toBe('desc');
    toggleSort('age');
    expect(sortKey.value).toBe('age');
    expect(sortDir.value).toBe('asc');
    expect(sortedItems.value.map((i) => i.age)).toEqual([25, 30, 35]);
  });

  it('sorts by number field', () => {
    const items = ref(makeItems());
    const { sortedItems, toggleSort } = useSort(items);
    toggleSort('age');
    expect(sortedItems.value.map((i) => i.age)).toEqual([25, 30, 35]);
    toggleSort('age');
    expect(sortedItems.value.map((i) => i.age)).toEqual([35, 30, 25]);
  });

  it('sorts by boolean field', () => {
    const items = ref(makeItems());
    const { sortedItems, toggleSort } = useSort(items);
    toggleSort('active');
    expect(sortedItems.value.map((i) => i.active)).toEqual([
      false,
      true,
      true,
    ]);
    toggleSort('active');
    expect(sortedItems.value.map((i) => i.active)).toEqual([
      true,
      true,
      false,
    ]);
  });

  it('handles null values', () => {
    const items = ref(makeItems());
    const { sortedItems, toggleSort } = useSort(items);
    toggleSort('score');
    const scores = sortedItems.value.map((i) => i.score);
    expect(scores[0]).toBe(80);
    expect(scores[1]).toBe(95);
    expect(scores[2]).toBeNull();
  });
});
