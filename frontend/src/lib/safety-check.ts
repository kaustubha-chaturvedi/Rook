import type { Rules } from './types';

export type RunOutcome = { safety: string; txState: string; blocked: boolean };

export function checkSQL(sql: string, rules: Rules, isProd: boolean): RunOutcome {
  const upper = sql.toUpperCase();

  if (rules.readOnlyConnection && /(UPDATE|DELETE|ALTER|DROP|TRUNCATE|INSERT|CREATE)/i.test(upper)) {
    return { safety: 'Blocked (Read-only connection)', txState: 'blocked', blocked: true };
  }
  if (/\bSHUTDOWN\b|\bXP_CMDSHELL\b/i.test(upper)) {
    return { safety: 'Blocked (Dangerous keyword)', txState: 'blocked', blocked: true };
  }
  if (/\bUPDATE\b/i.test(upper) && !/\bWHERE\b/i.test(upper)) {
    return { safety: 'Blocked (UPDATE without WHERE)', txState: 'blocked', blocked: true };
  }
  if (/\bDELETE\b/i.test(upper) && !/\bWHERE\b/i.test(upper)) {
    return { safety: 'Blocked (DELETE without WHERE)', txState: 'blocked', blocked: true };
  }
  if (/--\s*FORCE\b/i.test(sql) && isProd && !rules.allowForceBypassInProd) {
    return { safety: 'Blocked (--FORCE disabled in PROD)', txState: 'blocked', blocked: true };
  }

  const warnings: string[] = [];
  if (/\bSELECT\b/i.test(upper) && !/\bLIMIT\b|\bTOP\b/i.test(upper)) warnings.push('No LIMIT/TOP');
  if (/\bSELECT\s+\*/i.test(upper)) warnings.push('SELECT *');
  if (/\bWHERE\s+1\s*=\s*1\b|\bWHERE\s+TRUE\b|\bIS\s+NOT\s+NULL\b|!=\s*0/i.test(upper)) {
    warnings.push('Non-selective WHERE');
  }
  if (/\bCROSS\s+JOIN\b/i.test(upper) || (/\bJOIN\b/i.test(upper) && !/\bON\b/i.test(upper))) {
    warnings.push('Cartesian JOIN risk');
  }
  if (/\bDROP\b|\bTRUNCATE\b|\bKILL\b/i.test(upper)) warnings.push('Destructive operation');

  return {
    safety: warnings.length ? `Warning: ${warnings.join(', ')}` : 'Safe',
    txState: 'pending_commit',
    blocked: false
  };
}
