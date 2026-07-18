import type { ThemeMode } from './types';

export function applyTheme(mode: ThemeMode) {
  document.documentElement.classList.toggle('dark', mode === 'dark');
}

export function monacoTheme(mode: ThemeMode) {
  return mode === 'dark' ? 'vs-dark' : 'vs';
}
