<script lang="ts">
  import type { NotebookFile } from '../types';

  let {
    files,
    selectedIndex,
    onSelect,
    onImport,
    onPickImport
  }: {
    files: NotebookFile[];
    selectedIndex: number;
    onSelect: (index: number) => void;
    onImport: (file: File) => void;
    onPickImport: () => void;
  } = $props();

  let fileInput: HTMLInputElement;
  let menuOpen = $state(false);
  let menuX = $state(0);
  let menuY = $state(0);

  function onFileChosen(e: Event) {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (file) onImport(file);
    input.value = '';
    menuOpen = false;
  }

  function onSidebarContextMenu(e: MouseEvent) {
    e.preventDefault();
    menuOpen = true;
    menuX = e.clientX;
    menuY = e.clientY;
  }
</script>

<aside
  class="relative border-r border-border bg-card p-3 overflow-auto flex flex-col gap-3"
  oncontextmenu={onSidebarContextMenu}
  onclick={() => (menuOpen = false)}
>
  <div class="text-sm font-semibold tracking-tight">Files</div>
  <input bind:this={fileInput} type="file" accept=".snb,application/json" class="hidden" onchange={onFileChosen} />

  {#if menuOpen}
    <div
      class="fixed z-50 min-w-40 rounded-md border border-border bg-card p-1 shadow-lg"
      style={`left:${menuX}px; top:${menuY}px`}
      role="menu"
      tabindex="-1"
      onclick={(e) => e.stopPropagation()}
    >
      <button class="block w-full rounded px-2 py-1.5 text-left text-xs hover:bg-accent hover:text-accent-foreground" onclick={async () => { await onPickImport(); menuOpen = false; }}>
        Open...
      </button>
      <button class="block w-full rounded px-2 py-1.5 text-left text-xs hover:bg-accent hover:text-accent-foreground" onclick={() => fileInput?.click()}>
        Import .snb
      </button>
    </div>
  {/if}

  <div class="space-y-1">
    {#each files as file, idx}
      <button
        class="w-full rounded-md border px-2 py-2 text-left text-sm transition-colors {idx === selectedIndex
          ? 'border-primary/30 bg-primary/10 text-foreground'
          : 'border-input bg-background text-muted-foreground hover:bg-accent hover:text-accent-foreground'}"
        type="button"
        onclick={() => onSelect(idx)}
      >
        {file.name}
      </button>
    {/each}
  </div>
</aside>
