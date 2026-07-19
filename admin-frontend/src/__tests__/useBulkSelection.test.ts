import { describe, it, expect } from 'vitest';
import { ref } from 'vue';
import { useBulkSelection } from '../composables/useBulkSelection';

function makeItems(n: number) {
  return ref(Array.from({ length: n }, (_, i) => ({ id: i + 1 })));
}

describe('useBulkSelection', () => {
  it('starts with no selection', () => {
    const items = makeItems(5);
    const { selectedCount, headerChecked } = useBulkSelection(items);
    expect(selectedCount.value).toBe(0);
    expect(headerChecked.value).toBe(false);
  });

  it('toggleOne selects and deselects', () => {
    const items = makeItems(5);
    const { toggleOne, isSelected, selectedCount } = useBulkSelection(items);
    toggleOne(1, true);
    expect(isSelected(1)).toBe(true);
    expect(selectedCount.value).toBe(1);
    toggleOne(1, false);
    expect(isSelected(1)).toBe(false);
    expect(selectedCount.value).toBe(0);
  });

  it('toggleAll selects all items', () => {
    const items = makeItems(5);
    const { toggleAll, selectedCount, headerChecked } = useBulkSelection(items);
    toggleAll(true);
    expect(selectedCount.value).toBe(5);
    expect(headerChecked.value).toBe(true);
  });

  it('toggleAll(false) clears selection', () => {
    const items = makeItems(5);
    const { toggleAll, selectedCount } = useBulkSelection(items);
    toggleAll(true);
    toggleAll(false);
    expect(selectedCount.value).toBe(0);
  });

  it('headerChecked returns indeterminate for partial selection', () => {
    const items = makeItems(5);
    const { toggleOne, headerChecked } = useBulkSelection(items);
    toggleOne(1, true);
    toggleOne(3, true);
    expect(headerChecked.value).toBe('indeterminate');
  });

  it('selectedItems returns matching items', () => {
    const items = makeItems(5);
    const { toggleOne, selectedItems } = useBulkSelection(items);
    toggleOne(2, true);
    toggleOne(4, true);
    expect(selectedItems.value).toEqual([{ id: 2 }, { id: 4 }]);
  });

  it('clearSelection empties everything', () => {
    const items = makeItems(5);
    const { toggleAll, clearSelection, selectedCount } = useBulkSelection(items);
    toggleAll(true);
    clearSelection();
    expect(selectedCount.value).toBe(0);
  });
});
