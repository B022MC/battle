import React from 'react';
import { RolesView } from '@/components/(shop)/roles/roles-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function RolesScreen() {
  return (
    <RouteGuard anyOf={['role:view']}>
      <RolesView />
    </RouteGuard>
  );
}

