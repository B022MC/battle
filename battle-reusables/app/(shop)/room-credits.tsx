import { RoomCreditsView } from '@/components/(shop)/room-credits/room-credits-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function RoomCreditsScreen() {
  return (
    <RouteGuard anyOf={['game:room_credit:view']}>
      <RoomCreditsView />
    </RouteGuard>
  );
}
