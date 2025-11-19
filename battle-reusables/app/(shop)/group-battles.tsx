import { GroupBattlesView } from '@/components/(shop)/battles/group-battles-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function GroupBattlesScreen() {
  return (
    <RouteGuard anyOf={['shop:member:view', 'battles:view']}>
      <GroupBattlesView />
    </RouteGuard>
  );
}

