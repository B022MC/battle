import React from 'react';
import { PermissionsView } from '@/components/(shop)/permissions/permissions-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function PermissionsScreen() {
  return (
    <RouteGuard anyOf={['permission:view']}>
      <PermissionsView />
    </RouteGuard>
  );
}

