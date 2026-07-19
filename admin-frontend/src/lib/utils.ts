import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const MIN_LOADING_MS = 400;

export function minLoadingDelay() {
  return new Promise((r) => setTimeout(r, MIN_LOADING_MS));
}
