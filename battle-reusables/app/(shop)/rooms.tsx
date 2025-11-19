import { CtrlAccountsView } from '@/components/(shop)/ctrl-accounts/ctrl-accounts-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function RoomsScreen() {
  return (
    <RouteGuard anyOf={['game:ctrl:view', 'game:ctrl:create', 'game:ctrl:update', 'game:ctrl:delete']}>
      <CtrlAccountsView />
    </RouteGuard>
  );
}
