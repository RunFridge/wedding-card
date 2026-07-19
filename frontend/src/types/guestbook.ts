import type { GuestbookEntry } from './api';

export interface PendingAction {
  type: 'edit' | 'delete';
  entry: GuestbookEntry;
  password?: string;
}

export interface DialogState {
  visible: boolean;
  message: string;
}
