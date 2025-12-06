import { RoomCreditView } from '@/components/(shop)/room-credit/room-credit-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function RoomCreditScreen() {
  return (
    <RouteGuard anyOf={['room:credit:view']}>
      <RoomCreditView />
    </RouteGuard>
  );
}

