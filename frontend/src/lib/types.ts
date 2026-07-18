export type ThemeMode = 'dark' | 'light';
export type AutoAction = 'commit' | 'rollback';

export type Rules = {
  allowForceBypassInProd: boolean;
  strictProd: boolean;
  readOnlyConnection: boolean;
  largeRowsWarnThreshold: number;
  largeRowsDangerThreshold: number;
  largeRowsBlockThreshold: number;
  maxQueryTimeSeconds: number;
  maxConcurrentQueries: number;
};

export type NotebookFile = { name: string; sql: string; path?: string };

export type AppSettings = {
  mode: ThemeMode;
  autoAction: AutoAction;
  autoTimeoutSeconds: number;
  rules: Rules;
  showSettings: boolean;
};
