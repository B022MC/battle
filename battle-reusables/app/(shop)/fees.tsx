import { FeesView } from '@/components/(shop)/fees/fees-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function FeesScreen() {
  return (
    <RouteGuard anyOf={['shop:fees:view']}>
      <FeesView />
    </RouteGuard>
  );
}


