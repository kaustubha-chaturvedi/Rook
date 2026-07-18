/** Matches Go notebook.File / storage.ParseJSON. */
type SnbDoc = {
  version?: number;
  cells?: { id?: string; type?: string; connection?: string; query?: string }[];
};

export function sqlFromSnbDoc(doc: SnbDoc): string {
  const cells = doc.cells ?? [];
  return cells
    .map((c) => (c.query ?? '').trim())
    .filter(Boolean)
    .join('\n\n');
}

export function snbDocFromSQL(sql: string): SnbDoc {
  const parts = sql
    .split(/\n\n+/)
    .map((s) => s.trim())
    .filter(Boolean);
  return {
    version: 1,
    cells: parts.map((query, i) => ({
      id: `cell-${i + 1}`,
      type: 'sql',
      connection: '',
      query
    }))
  };
}

export function parseSnbFile(name: string, content: string): { name: string; sql: string } {
  const doc = JSON.parse(content) as SnbDoc;
  return { name, sql: sqlFromSnbDoc(doc) };
}

export function downloadSnb(name: string, sql: string) {
  const blob = new Blob([JSON.stringify(snbDocFromSQL(sql), null, 2)], {
    type: 'application/json'
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = name.endsWith('.snb') ? name : `${name}.snb`;
  a.click();
  URL.revokeObjectURL(url);
}
