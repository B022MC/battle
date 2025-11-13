import { useMemo } from 'react';
import { useAuthStore } from './use-auth-store';

export function usePermission() {
  const roles = useAuthStore((s) => s.roles) ?? [];
  const perms = useAuthStore((s) => s.perms) ?? [];
  const userRole = useAuthStore((s) => s.role);

  const isSuperAdmin = useMemo(() => roles.includes(1), [roles]);
  const isStoreAdmin = useMemo(() => userRole === 'store_admin', [userRole]);
  const set = useMemo(() => {
    const s = new Set<string>();
    for (const p of perms) if (p) s.add(String(p).toLowerCase().trim());
    return s;
  }, [perms]);

  const hasPerm = (p: string) => isSuperAdmin || set.has(p.toLowerCase().trim());
  const hasAny = (codes: string[]) => isSuperAdmin || codes.some((c) => hasPerm(c));
  const hasAll = (codes: string[]) => isSuperAdmin || codes.every((c) => hasPerm(c));

  return { isSuperAdmin, isStoreAdmin, hasPerm, hasAny, hasAll, perms, roles };
}


