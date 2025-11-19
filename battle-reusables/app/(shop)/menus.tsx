import React from 'react';
import { MenusView } from '@/components/(shop)/menus/menus-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function MenusScreen() {
  return (
    <RouteGuard anyOf={['menu:view']}>
      <MenusView />
    </RouteGuard>
  );
}


