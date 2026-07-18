<script lang="ts">
  import type { ThemeMode } from '../types';

  let {
    mode = $bindable(),
    connectionStatus,
    connectionEnv,
    isConnected,
    onOpenConnections,
    onToggleSettings
  }: {
    mode: ThemeMode;
    connectionStatus: string;
    connectionEnv: 'prod' | 'testing' | 'dev';
    isConnected: boolean;
    onOpenConnections: () => void;
    onToggleSettings: () => void;
  } = $props();
</script>

<header class="h-14 border-b border-border bg-card px-4 flex items-center justify-between">
  <div class="flex items-center gap-4">
    <span class="text-sm font-semibold tracking-tight">Rook</span>
    <button
      class="inline-flex items-center gap-2 rounded-md border px-2.5 py-1 text-xs transition-colors {isConnected
        ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300 hover:bg-emerald-500/20'
        : 'border-border bg-secondary text-muted-foreground hover:bg-accent hover:text-accent-foreground'}"
      onclick={onOpenConnections}
      title="Open connection manager"
    >
      <span class="h-1.5 w-1.5 rounded-full {isConnected ? 'bg-emerald-400' : 'bg-zinc-500'}"></span>
      {connectionStatus}
    </button>
    <span
      class="inline-flex items-center rounded-full border px-2 py-0.5 text-[11px] font-medium {connectionEnv ===
      'prod'
        ? 'border-destructive/30 bg-destructive/10 text-destructive'
        : connectionEnv === 'testing'
          ? 'border-amber-500/30 bg-amber-500/10 text-amber-300'
          : 'border-sky-500/30 bg-sky-500/10 text-sky-300'}"
    >
      {connectionEnv.toUpperCase()}
    </span>
  </div>

  <div class="flex items-center gap-2">
    <button
      class="inline-flex h-9 w-9 items-center justify-center rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground"
      onclick={onToggleSettings}
      title="Open settings"
      aria-label="Open settings"
    >
      <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current" stroke-width="2">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M19.4 15a1.7 1.7 0 0 0 .34 1.87l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.7 1.7 0 0 0-1.87-.34 1.7 1.7 0 0 0-1 1.55V21a2 2 0 1 1-4 0v-.09a1.7 1.7 0 0 0-1-1.55 1.7 1.7 0 0 0-1.87.34l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.7 1.7 0 0 0 4.6 15a1.7 1.7 0 0 0-1.55-1H3a2 2 0 1 1 0-4h.09a1.7 1.7 0 0 0 1.55-1A1.7 1.7 0 0 0 4.3 7.13l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.7 1.7 0 0 0 9 4.6a1.7 1.7 0 0 0 1-1.55V3a2 2 0 1 1 4 0v.09a1.7 1.7 0 0 0 1 1.55 1.7 1.7 0 0 0 1.87-.34l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.7 1.7 0 0 0 19.4 9c.24.62.85 1.03 1.51 1H21a2 2 0 1 1 0 4h-.09c-.66 0-1.27.41-1.51 1z"></path>
      </svg>
    </button>
    <button
      class="inline-flex h-9 w-9 items-center justify-center rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground"
      onclick={() => (mode = mode === 'dark' ? 'light' : 'dark')}
      title="Toggle theme"
      aria-label="Toggle theme"
    >
      {#if mode === 'dark'}
        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current" stroke-width="2">
          <circle cx="12" cy="12" r="4"></circle>
          <path d="M12 2v2"></path><path d="M12 20v2"></path><path d="M2 12h2"></path><path d="M20 12h2"></path>
          <path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path>
          <path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path>
        </svg>
      {:else}
        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current">
          <path d="M21 12.79A9 9 0 1 1 11.21 3c0 .26-.01.52-.01.79A7 7 0 0 0 21 12.79z"></path>
        </svg>
      {/if}
    </button>
  </div>
</header>
