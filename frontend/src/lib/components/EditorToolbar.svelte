<script lang="ts">
  export let safety: string;
  export let txState: string;
  export let durationMs: number;
  export let onRun: () => void;
  export let onRunAll: () => void;
  export let onCancel: () => void;
</script>

<div class="flex items-center gap-2">
  <button class="inline-flex h-9 w-9 items-center justify-center rounded-md bg-primary text-primary-foreground hover:opacity-90" onclick={onRun} title="Run current cell" aria-label="Run current cell">
    <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M8 5v14l11-7z"></path></svg>
  </button>
  <button class="inline-flex h-9 w-9 items-center justify-center rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground" onclick={onRunAll} title="Run whole file" aria-label="Run whole file">
    <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current" stroke-width="2"><path d="M4 6h16"></path><path d="M4 12h16"></path><path d="M4 18h10"></path><path d="m16 16 4 2-4 2z"></path></svg>
  </button>
  <button class="inline-flex h-9 w-9 items-center justify-center rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground" onclick={onCancel} title="Stop execution" aria-label="Stop execution">
    <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><rect x="7" y="7" width="10" height="10" rx="1"></rect></svg>
  </button>
  <span class="inline-flex h-9 items-center rounded-md bg-secondary px-3 text-xs text-secondary-foreground">Safety: {safety}</span>
  <span class="inline-flex h-9 items-center rounded-md bg-secondary px-3 text-xs text-secondary-foreground">Tx: {txState}</span>
  <span class="inline-flex h-9 items-center rounded-md bg-secondary px-3 text-xs text-secondary-foreground">{durationMs}ms</span>
</div>
