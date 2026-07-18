<script lang="ts">
  import type { Rules } from '../types';

  let { rules = $bindable() }: { rules: Rules } = $props();
</script>

<div class="rounded-lg border border-border bg-card p-3">
  <div class="mb-2 text-sm font-semibold">Safety Rules</div>
  <div class="space-y-2.5">
    <div class="flex items-center justify-between rounded-md border border-border bg-background/40 px-2.5 py-1.5">
      <div class="text-sm">Strict PROD</div>
      <button type="button" role="switch" aria-checked={rules.strictProd} aria-label="Toggle strict PROD" class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {rules.strictProd ? 'bg-primary' : 'bg-secondary'}" onclick={() => (rules.strictProd = !rules.strictProd)}>
        <span class="inline-block h-5 w-5 transform rounded-full bg-background shadow transition-transform {rules.strictProd ? 'translate-x-5' : 'translate-x-1'}"></span>
      </button>
    </div>
    <div class="flex items-center justify-between rounded-md border border-border bg-background/40 px-2.5 py-1.5">
      <div class="text-sm">Allow FORCE in PROD</div>
      <button type="button" role="switch" aria-checked={rules.allowForceBypassInProd} aria-label="Toggle FORCE in PROD" class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {rules.allowForceBypassInProd ? 'bg-primary' : 'bg-secondary'}" onclick={() => (rules.allowForceBypassInProd = !rules.allowForceBypassInProd)}>
        <span class="inline-block h-5 w-5 transform rounded-full bg-background shadow transition-transform {rules.allowForceBypassInProd ? 'translate-x-5' : 'translate-x-1'}"></span>
      </button>
    </div>
    <div class="flex items-center justify-between rounded-md border border-border bg-background/40 px-2.5 py-1.5">
      <div class="text-sm">Read-only connection</div>
      <button type="button" role="switch" aria-checked={rules.readOnlyConnection} aria-label="Toggle read-only connection" class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {rules.readOnlyConnection ? 'bg-primary' : 'bg-secondary'}" onclick={() => (rules.readOnlyConnection = !rules.readOnlyConnection)}>
        <span class="inline-block h-5 w-5 transform rounded-full bg-background shadow transition-transform {rules.readOnlyConnection ? 'translate-x-5' : 'translate-x-1'}"></span>
      </button>
    </div>
    <div class="grid grid-cols-1 gap-2 md:grid-cols-2">
      <label class="text-sm">Max Query Sec<input type="number" bind:value={rules.maxQueryTimeSeconds} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
      <label class="text-sm">Warn Rows<input type="number" bind:value={rules.largeRowsWarnThreshold} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
      <label class="text-sm">Danger Rows<input type="number" bind:value={rules.largeRowsDangerThreshold} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
      <label class="text-sm">Block Rows<input type="number" bind:value={rules.largeRowsBlockThreshold} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
      <label class="text-sm">Max Concurrency<input type="number" bind:value={rules.maxConcurrentQueries} class="mt-1 h-9 w-full rounded-md border border-input bg-background px-2" /></label>
    </div>
  </div>
</div>
