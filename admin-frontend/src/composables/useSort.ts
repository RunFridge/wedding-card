import { ref, computed, type Ref } from 'vue';

export type SortDir = 'asc' | 'desc';

export function useSort<T>(items: Ref<T[]>) {
  const sortKey = ref<string | null>(null);
  const sortDir = ref<SortDir>('asc');

  function toggleSort(key: string) {
    if (sortKey.value === key) {
      if (sortDir.value === 'asc') {
        sortDir.value = 'desc';
      } else {
        sortKey.value = null;
        sortDir.value = 'asc';
      }
    } else {
      sortKey.value = key;
      sortDir.value = 'asc';
    }
  }

  const sortedItems = computed(() => {
    if (!sortKey.value) return items.value;
    const key = sortKey.value;
    const dir = sortDir.value === 'asc' ? 1 : -1;
    return [...items.value].sort((a, b) => {
      const av = (a as Record<string, unknown>)[key];
      const bv = (b as Record<string, unknown>)[key];
      if (av == null && bv == null) return 0;
      if (av == null) return dir;
      if (bv == null) return -dir;
      if (typeof av === 'number' && typeof bv === 'number') return (av - bv) * dir;
      if (typeof av === 'boolean' && typeof bv === 'boolean') return (Number(av) - Number(bv)) * dir;
      return String(av).localeCompare(String(bv)) * dir;
    });
  });

  return { sortKey, sortDir, toggleSort, sortedItems };
}
