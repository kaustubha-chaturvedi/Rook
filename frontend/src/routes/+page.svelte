<script lang="ts">
  import { onMount } from 'svelte';
  import * as monaco from 'monaco-editor';
  import AppHeader from '$lib/components/AppHeader.svelte';
  import FileSidebar from '$lib/components/FileSidebar.svelte';
  import SettingsPanel from '$lib/components/SettingsPanel.svelte';
  import EditorToolbar from '$lib/components/EditorToolbar.svelte';
  import AppFooter from '$lib/components/AppFooter.svelte';
  import { sampleFiles } from '$lib/defaults';
  import { importSnbFile, pickAndImportSnb, saveSnb } from '$lib/notebook-api';
  import { loadSettings, saveSettings } from '$lib/settings';
  import { checkSQL } from '$lib/safety-check';
  import { applyTheme, monacoTheme } from '$lib/theme';
  import type { AppSettings, NotebookFile } from '$lib/types';

  type ConnectionCategory = 'prod' | 'testing' | 'dev';
  type SavedConnection = {
    label: string;
    category: ConnectionCategory;
    host: string;
    port: string;
    database: string;
    user: string;
    password: string;
    ssl: boolean;
  };
  type CellScope = { id: string; sql: string; startOffset: number; endOffset: number };
  type CellMarker = { id: string; label: string; top: number; left: number };

  const CELL_DELIMITER = '\n\n\n';
  const CONNECTIONS_KEY = 'rook.connections.v1';
  const EDITOR_TOP_PADDING = 36;
  const CELL_BADGE_HEIGHT = 24;

  let editorEl: HTMLDivElement;
  let editor: monaco.editor.IStandaloneCodeEditor;
  let cellScopes: CellScope[] = [];
  let cellMarkers: CellMarker[] = [];
  let activeCellId = '';

  let settings: AppSettings = loadSettings();
  let files: NotebookFile[] = [...sampleFiles];
  let selectedFile = 0;
  let safety = 'Safe';
  let txState = 'pending_commit';
  let durationMs = 0;

  let showConnectionPopup = false;
  let showConnectionSecret = false;
  let isConnected = false;
  let connectionStatus = 'Disconnected';
  let connectionError = '';
  let selectedConnectionLabel = '';
  let savedConnections: SavedConnection[] = [];
  let connection: SavedConnection = {
    label: 'Primary',
    category: 'testing',
    host: 'localhost',
    port: '5432',
    database: 'postgres',
    user: 'readonly_user',
    password: '',
    ssl: false
  };

  $: applyTheme(settings.mode);
  $: saveSettings(settings);

  if (typeof localStorage !== 'undefined' && savedConnections.length === 0) {
    const raw = localStorage.getItem(CONNECTIONS_KEY);
    if (raw) {
      try {
        const parsed = JSON.parse(raw) as Array<Record<string, unknown>>;
        savedConnections = parsed.map((c) => ({
          label: String(c.label ?? 'Connection'),
          category: (c.category === 'prod' || c.category === 'dev' || c.category === 'testing'
            ? c.category
            : 'testing') as ConnectionCategory,
          host: String(c.host ?? ''),
          port: String(c.port ?? ''),
          database: String(c.database ?? ''),
          user: String(c.user ?? ''),
          password: String(c.password ?? ''),
          ssl: Boolean(c.ssl)
        }));
      } catch {
        savedConnections = [];
      }
    }
  }

  function loadFile(index: number) {
    selectedFile = index;
    editor?.setValue(files[index].sql);
  }

  function syncEditorToFile() {
    if (!editor) return;
    files[selectedFile] = { ...files[selectedFile], sql: editor.getValue() };
  }

  async function handleImport(file: File) {
    const imported = await importSnbFile(file);
    files = [...files, imported];
    selectedFile = files.length - 1;
    editor?.setValue(imported.sql);
  }

  async function handlePickImport() {
    const imported = await pickAndImportSnb();
    if (!imported?.sql) return;
    files = [...files, imported];
    selectedFile = files.length - 1;
    editor?.setValue(imported.sql);
  }

  async function handleSaveNotebook() {
    syncEditorToFile();
    const current = files[selectedFile];
    const updated = await saveSnb(current, current.sql);
    if (updated) files[selectedFile] = updated;
  }

  function splitCells(sql: string): string[] {
    return sql.split(CELL_DELIMITER).map((part) => part.trim()).filter(Boolean);
  }

  function buildCellScopes(sql: string): CellScope[] {
    const scopes: CellScope[] = [];
    let cursor = 0;
    let idx = 0;
    while (cursor <= sql.length) {
      const next = sql.indexOf(CELL_DELIMITER, cursor);
      const end = next < 0 ? sql.length : next;
      const trimmed = sql.slice(cursor, end).trim();
      if (trimmed) scopes.push({ id: `cell-${idx++}`, sql: trimmed, startOffset: cursor, endOffset: end });
      if (next < 0) break;
      cursor = next + CELL_DELIMITER.length;
    }
    return scopes;
  }

  function rebuildCellMarkers() {
    if (!editor) return;
    const model = editor.getModel();
    if (!model) return;
    const layout = editor.getLayoutInfo();
    const left = layout.contentLeft + Math.round(layout.contentWidth / 2);
    cellMarkers = cellScopes
      .map((cell, i) => {
        const pos = model.getPositionAt(cell.startOffset);
        const visible = editor.getScrolledVisiblePosition(pos);
        if (!visible) return null;
        const top = Math.max(8, visible.top - CELL_BADGE_HEIGHT + 3);
        return { id: cell.id, label: `Cell ${i + 1}`, top, left };
      })
      .filter((marker): marker is CellMarker => marker !== null);
  }

  function currentCellSQL(): string {
    if (!editor) return '';
    const model = editor.getModel();
    const pos = editor.getPosition();
    if (!model || !pos) return '';
    const sql = model.getValue();
    const cursorOffset = model.getOffsetAt(pos);
    const left = sql.lastIndexOf(CELL_DELIMITER, Math.max(0, cursorOffset - 1));
    const right = sql.indexOf(CELL_DELIMITER, cursorOffset);
    const start = left < 0 ? 0 : left + CELL_DELIMITER.length;
    const end = right < 0 ? sql.length : right;
    return sql.slice(start, end).trim();
  }

  function currentCellId(): string {
    if (!editor) return '';
    const model = editor.getModel();
    const pos = editor.getPosition();
    if (!model || !pos) return '';
    const offset = model.getOffsetAt(pos);
    return cellScopes.find((cell) => offset >= cell.startOffset && offset <= cell.endOffset)?.id ?? '';
  }

  async function executeSQL(sql: string) {
    if (!sql.trim()) {
      safety = 'Blocked (Empty cell)';
      txState = 'blocked';
      return;
    }
    const start = performance.now();
    const outcome = checkSQL(sql, settings.rules, connection.category === 'prod');
    safety = outcome.safety;
    txState = outcome.txState;
    if (outcome.blocked) return;
    await new Promise((resolve) => setTimeout(resolve, Math.max(1, settings.autoTimeoutSeconds) * 1000));
    txState = `auto_${settings.autoAction}`;
    durationMs = Math.round(performance.now() - start);
  }

  async function runCurrentCell() {
    await executeSQL(currentCellSQL());
  }

  async function runAllCells() {
    await executeSQL(splitCells(editor?.getValue() ?? '').join('\n'));
  }

  async function runCellById(cellId: string) {
    const found = cellScopes.find((cell) => cell.id === cellId);
    if (found) await executeSQL(found.sql);
  }

  function focusCellById(cellId: string) {
    if (!editor) return;
    const model = editor.getModel();
    const found = cellScopes.find((cell) => cell.id === cellId);
    if (!model || !found) return;
    const pos = model.getPositionAt(found.startOffset);
    editor.setPosition(pos);
    editor.revealPositionInCenter(pos);
    editor.focus();
    activeCellId = cellId;
    rebuildCellMarkers();
  }

  function cancelCell() {
    txState = 'cancelled_rollback';
  }

  function connectToDatabase() {
    connectionError = '';
    if (!connection.host || !connection.database || !connection.user) {
      connectionError = 'Host, database, and user are required.';
      return;
    }
    isConnected = true;
    connectionStatus = connection.label || `${connection.database}@${connection.host}`;
    showConnectionPopup = false;
  }

  function disconnectDatabase() {
    isConnected = false;
    connectionStatus = 'Disconnected';
  }

  function saveConnectionPreset() {
    connectionError = '';
    if (!connection.label?.trim()) {
      connectionError = 'Label is required to save a connection.';
      return;
    }
    savedConnections = savedConnections.filter((item) => item.label !== connection.label).concat({ ...connection });
    localStorage.setItem(CONNECTIONS_KEY, JSON.stringify(savedConnections));
    selectedConnectionLabel = connection.label;
  }

  function selectConnection(label: string) {
    selectedConnectionLabel = label;
    const item = savedConnections.find((saved) => saved.label === label);
    if (!item) return;
    connection = { ...item };
    connectionError = '';
  }

  onMount(() => {
    editor = monaco.editor.create(editorEl, {
      language: 'sql',
      value: files[selectedFile].sql,
      minimap: { enabled: false },
      automaticLayout: true,
      fontSize: 13,
      padding: { top: EDITOR_TOP_PADDING, bottom: 24 },
      scrollBeyondLastLine: false
    });
    cellScopes = buildCellScopes(editor.getValue());
    monaco.editor.setTheme(monacoTheme(settings.mode));
    rebuildCellMarkers();
    activeCellId = currentCellId();

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, runCurrentCell);
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyO, () => void handlePickImport());
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyI, () => {
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = '.snb,application/json';
      input.onchange = async () => {
        const file = input.files?.[0];
        if (file) await handleImport(file);
      };
      input.click();
    });
    editor.addCommand(monaco.KeyCode.Escape, cancelCell);
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
      saveSettings(settings);
      void handleSaveNotebook();
    });
    editor.onDidChangeModelContent(() => {
      cellScopes = buildCellScopes(editor.getValue());
      activeCellId = currentCellId();
      rebuildCellMarkers();
    });
    editor.onDidChangeCursorPosition(() => {
      activeCellId = currentCellId();
    });
    editor.onDidScrollChange(rebuildCellMarkers);
    editor.onDidLayoutChange(rebuildCellMarkers);
    return () => editor?.dispose();
  });

  $: if (editor) monaco.editor.setTheme(monacoTheme(settings.mode));
</script>

<main class="h-screen bg-background text-foreground grid grid-rows-[56px_1fr_44px] text-sm">
  <AppHeader
    bind:mode={settings.mode}
    {connectionStatus}
    connectionEnv={connection.category}
    {isConnected}
    onOpenConnections={() => (showConnectionPopup = true)}
    onToggleSettings={() => (settings.showSettings = !settings.showSettings)}
  />

  <section class="grid grid-cols-[280px_1fr] min-h-0">
    <FileSidebar {files} selectedIndex={selectedFile} onSelect={loadFile} onImport={handleImport} onPickImport={handlePickImport} />
    <div class="p-3 flex flex-col gap-3 min-h-0 bg-background">
      <EditorToolbar {safety} {txState} {durationMs} onRun={runCurrentCell} onRunAll={runAllCells} onCancel={cancelCell} />
      <div class="relative h-full rounded-lg border border-border bg-card overflow-hidden">
        {#each cellMarkers as marker}
          <div
            class="absolute z-10 inline-flex items-center gap-1 rounded-md border border-border bg-black/25 px-2 py-1 text-xs text-zinc-100 backdrop-blur-sm transition-colors hover:bg-black/45"
            style={`top:${marker.top}px; left:${marker.left}px; transform: translateX(-50%);`}
          >
            <button class="font-medium {activeCellId === marker.id ? 'text-primary' : ''}" onclick={() => focusCellById(marker.id)}>{marker.label}</button>
            <button class="inline-flex h-5 w-5 items-center justify-center rounded hover:bg-primary/30" title="Run cell" aria-label="Run cell" onclick={() => runCellById(marker.id)}>
              <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><path d="M8 5v14l11-7z"></path></svg>
            </button>
            <button class="inline-flex h-5 w-5 items-center justify-center rounded hover:bg-destructive/30" title="Stop cell" aria-label="Stop cell" onclick={cancelCell}>
              <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><rect x="7" y="7" width="10" height="10" rx="1"></rect></svg>
            </button>
          </div>
        {/each}
        <div bind:this={editorEl} class="h-full"></div>
      </div>
    </div>
  </section>

  <AppFooter />
</main>

{#if settings.showSettings}
  <div class="fixed inset-0 z-40 bg-black/55 backdrop-blur-[2px] flex items-center justify-center p-6 animate-fade-in" role="button" tabindex="0" onclick={() => (settings.showSettings = false)} onkeydown={(e) => e.key === 'Escape' && (settings.showSettings = false)}>
    <div class="w-full max-w-2xl rounded-xl border border-border bg-card shadow-2xl animate-pop-in" role="dialog" aria-modal="true" tabindex="-1" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()}>
      <div class="flex items-center justify-between border-b border-border px-4 py-3">
        <div class="text-sm font-semibold">Execution and Safety Settings</div>
        <button class="h-8 rounded-md border border-input bg-background px-3 text-xs hover:bg-accent hover:text-accent-foreground" onclick={() => (settings.showSettings = false)}>Close</button>
      </div>
      <div class="p-4 space-y-3">
        <div class="rounded-lg border border-border bg-card p-3">
          <div class="mb-3 text-sm font-semibold">Execution Defaults</div>
          <div class="space-y-2">
            <div class="flex items-center justify-between rounded-md border border-border bg-background/40 px-3 py-2">
              <div class="text-sm">Auto Action</div>
              <div class="inline-flex rounded-md border border-input bg-background p-0.5">
                <button class="h-8 px-3 rounded text-xs font-medium {settings.autoAction === 'rollback' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}" onclick={() => (settings.autoAction = 'rollback')}>Rollback</button>
                <button class="h-8 px-3 rounded text-xs font-medium {settings.autoAction === 'commit' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}" onclick={() => (settings.autoAction = 'commit')}>Commit</button>
              </div>
            </div>
            <label class="block text-sm">Auto Action Timeout (seconds)<input type="number" min="1" bind:value={settings.autoTimeoutSeconds} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
            <p class="text-xs text-muted-foreground">Cell delimiter: press Enter 3 times.</p>
          </div>
        </div>
        <SettingsPanel bind:rules={settings.rules} />
      </div>
    </div>
  </div>
{/if}

{#if showConnectionPopup}
  <div class="fixed inset-0 z-50 bg-black/55 backdrop-blur-[2px] flex items-center justify-center p-6 animate-fade-in" role="button" tabindex="0" onclick={() => (showConnectionPopup = false)} onkeydown={(e) => e.key === 'Escape' && (showConnectionPopup = false)}>
    <div class="w-full max-w-lg rounded-xl border border-border bg-card shadow-2xl animate-pop-in" role="dialog" aria-modal="true" tabindex="-1" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()}>
      <div class="flex items-center justify-between border-b border-border px-4 py-3">
        <div class="text-sm font-semibold">Connection Manager</div>
        <button class="h-8 rounded-md border border-input bg-background px-3 text-xs hover:bg-accent hover:text-accent-foreground" onclick={() => (showConnectionPopup = false)}>Close</button>
      </div>
      <div class="space-y-3 p-4">
        <label class="block text-xs text-muted-foreground">Saved Connections<select class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" bind:value={selectedConnectionLabel} onchange={(e) => selectConnection((e.currentTarget as HTMLSelectElement).value)}><option value="">Select saved connection</option>{#each savedConnections as conn}<option value={conn.label}>{conn.label}</option>{/each}</select></label>
        <label class="block text-xs text-muted-foreground">Label<input bind:value={connection.label} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
        <label class="block text-xs text-muted-foreground">Category<select bind:value={connection.category} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm"><option value="prod">prod</option><option value="testing">testing</option><option value="dev">dev</option></select></label>
        <label class="block text-xs text-muted-foreground">Host<input bind:value={connection.host} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
        <div class="grid grid-cols-2 gap-2">
          <label class="block text-xs text-muted-foreground">Port<input bind:value={connection.port} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
          <label class="block text-xs text-muted-foreground">Database<input bind:value={connection.database} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
        </div>
        <label class="block text-xs text-muted-foreground">User<input bind:value={connection.user} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
        <label class="block text-xs text-muted-foreground">Password<input bind:value={connection.password} type={showConnectionSecret ? 'text' : 'password'} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2 text-sm" /></label>
        <div class="flex items-center justify-between">
          <div class="inline-flex rounded-md border border-input bg-background p-0.5">
            <button class="h-8 px-3 rounded text-xs font-medium {connection.ssl ? 'bg-primary text-primary-foreground' : 'text-muted-foreground'}" onclick={() => (connection.ssl = true)}>SSL On</button>
            <button class="h-8 px-3 rounded text-xs font-medium {!connection.ssl ? 'bg-primary text-primary-foreground' : 'text-muted-foreground'}" onclick={() => (connection.ssl = false)}>SSL Off</button>
          </div>
          <button class="h-8 rounded-md border border-input bg-background px-3 text-xs hover:bg-accent hover:text-accent-foreground" onclick={() => (showConnectionSecret = !showConnectionSecret)}>{showConnectionSecret ? 'Hide Secret' : 'Show Secret'}</button>
        </div>
        {#if connectionError}<p class="text-xs text-destructive">{connectionError}</p>{/if}
        <div class="grid grid-cols-3 gap-2">
          <button class="h-9 rounded-md border border-input bg-background px-3 text-sm hover:bg-accent hover:text-accent-foreground" onclick={saveConnectionPreset}>Save Label</button>
          {#if isConnected}<button class="h-9 rounded-md bg-destructive px-3 text-sm font-medium text-destructive-foreground" onclick={disconnectDatabase}>Disconnect</button>{:else}<button class="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground" onclick={connectToDatabase}>Connect</button>{/if}
          <button class="h-9 rounded-md border border-input bg-background px-3 text-sm hover:bg-accent hover:text-accent-foreground" onclick={() => (showConnectionPopup = false)}>Done</button>
        </div>
      </div>
    </div>
  </div>
{/if}
