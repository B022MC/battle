import { GroupBalancesView } from '@/components/(shop)/battles/group-balances-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function GroupBalancesScreen() {
  return (
    <RouteGuard anyOf={['shop:member:view', 'fund:wallet:view']}>
      <GroupBalancesView />
    </RouteGuard>
  );
}

