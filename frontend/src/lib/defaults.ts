import type { AppSettings, NotebookFile, Rules } from './types';

export const SETTINGS_KEY = 'rook.settings.v2';
export const IS_PROD = true;

export const defaultRules: Rules = {
  allowForceBypassInProd: false,
  strictProd: true,
  readOnlyConnection: false,
  largeRowsWarnThreshold: 100,
  largeRowsDangerThreshold: 1000,
  largeRowsBlockThreshold: 10000,
  maxQueryTimeSeconds: 60,
  maxConcurrentQueries: 5
};

export const sampleFiles: NotebookFile[] = [
  {
    name: 'prod-checks.snb',
    sql: 'SELECT * FROM users LIMIT 50;\n\n\n--FORCE\nUPDATE users SET active = 0 WHERE id=1;'
  },
  {
    name: 'rollback-plan.snb',
    sql: 'DELETE FROM logs WHERE created_at < "2024-01-01";'
  }
];

export const defaultSettings: AppSettings = {
  mode: 'dark',
  autoAction: 'rollback',
  autoTimeoutSeconds: 5,
  rules: defaultRules,
  showSettings: true
};
