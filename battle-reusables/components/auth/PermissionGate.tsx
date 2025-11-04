import React from 'react';
import { usePermission } from '@/hooks/use-permission';

type Props = {
  anyOf?: string[];
  allOf?: string[];
  fallback?: React.ReactNode;
  children: React.ReactNode;
};

export function PermissionGate({ anyOf, allOf, fallback = null, children }: Props) {
  const { isSuperAdmin, hasAny, hasAll } = usePermission();
  const ok =
    isSuperAdmin ||
    ((allOf && allOf.length ? hasAll(allOf) : true) &&
      (anyOf && anyOf.length ? hasAny(anyOf) : true));

  return ok ? <>{children}</> : <>{fallback}</>;
}


