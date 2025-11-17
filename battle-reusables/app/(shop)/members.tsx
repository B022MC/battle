import { MembersView } from '@/components/(tabs)/members/members-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function MembersScreen() {
  return (
    <RouteGuard anyOf={['shop:member:view', 'shop:member:kick']}>
      <MembersView />
    </RouteGuard>
  );
}
