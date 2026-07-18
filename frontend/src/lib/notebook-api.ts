import { downloadSnb, parseSnbFile } from './snb';
import type { NotebookFile } from './types';

export type NotebookView = NotebookFile;

/** Wails-bound notebook service (present when running inside the desktop app). */
type WailsNotebook = {
  ImportContent: (name: string, content: string) => Promise<NotebookView>;
  ImportPath: (path: string) => Promise<NotebookView>;
  PickAndImport: () => Promise<NotebookView>;
  Save: (path: string, sql: string) => Promise<void>;
  PickAndSave: (sql: string) => Promise<string>;
};

function wailsNotebook(): WailsNotebook | undefined {
  const g = globalThis as {
    go?: { notebook?: { Service?: WailsNotebook } };
  };
  return g.go?.notebook?.Service;
}

export async function importSnbFile(file: File): Promise<NotebookView> {
  const content = await file.text();
  const name = file.name;
  const api = wailsNotebook();
  if (api?.ImportContent) {
    return api.ImportContent(name, content);
  }
  return { ...parseSnbFile(name, content) };
}

export async function pickAndImportSnb(): Promise<NotebookView | null> {
  const api = wailsNotebook();
  if (!api?.PickAndImport) return null;
  const view = await api.PickAndImport();
  return view?.sql ? view : null;
}

export async function saveSnb(file: NotebookFile, sql: string): Promise<NotebookView | null> {
  const api = wailsNotebook();
  if (file.path && api?.Save) {
    await api.Save(file.path, sql);
    return { ...file, sql };
  }
  if (api?.PickAndSave) {
    const path = await api.PickAndSave(sql);
    if (path) return { name: file.name, sql, path };
  }
  downloadSnb(file.name, sql);
  return { ...file, sql };
}
