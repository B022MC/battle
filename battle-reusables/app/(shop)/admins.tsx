import { AdminsView } from '@/components/(shop)/admins/admins-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function AdminsScreen() {
  return (
    <RouteGuard anyOf={['shop:admin:view', 'shop:admin:assign', 'shop:admin:revoke']}>
      <AdminsView />
    </RouteGuard>
  );
}

