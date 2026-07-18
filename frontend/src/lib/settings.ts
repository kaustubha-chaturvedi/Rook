import { SETTINGS_KEY, defaultSettings } from './defaults';
import type { AppSettings } from './types';

export function loadSettings(): AppSettings {
  const raw = localStorage.getItem(SETTINGS_KEY);
  if (!raw) return { ...defaultSettings, rules: { ...defaultSettings.rules } };
  const saved = JSON.parse(raw) as Partial<AppSettings>;
  return {
    ...defaultSettings,
    ...saved,
    rules: { ...defaultSettings.rules, ...(saved.rules ?? {}) }
  };
}

export function saveSettings(settings: AppSettings) {
  localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings));
}
