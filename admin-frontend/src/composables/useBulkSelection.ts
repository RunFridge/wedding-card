import { computed, type Ref, ref } from 'vue';

export function useBulkSelection<T extends { id: number }>(items: Ref<T[]>) {
  const selectedIds = ref(new Set<number>()) as Ref<Set<number>>;

  const selectedCount = computed(() => selectedIds.value.size);

  const headerChecked = computed<boolean | 'indeterminate'>(() => {
    if (selectedIds.value.size === 0) return false;
    if (selectedIds.value.size === items.value.length) return true;
    return 'indeterminate';
  });

  const selectedItems = computed(() =>
    items.value.filter((item) => selectedIds.value.has(item.id)),
  );

  function isSelected(id: number) {
    return selectedIds.value.has(id);
  }

  function toggleOne(id: number, checked: boolean) {
    const next = new Set(selectedIds.value);
    if (checked) next.add(id);
    else next.delete(id);
    selectedIds.value = next;
  }

  function toggleAll(checked: boolean) {
    if (checked) {
      selectedIds.value = new Set(items.value.map((item) => item.id));
    } else {
      selectedIds.value = new Set();
    }
  }

  function clearSelection() {
    selectedIds.value = new Set();
  }

  return {
    selectedIds,
    selectedCount,
    headerChecked,
    selectedItems,
    isSelected,
    toggleOne,
    toggleAll,
    clearSelection,
  };
}
