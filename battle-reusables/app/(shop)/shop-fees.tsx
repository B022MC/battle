import { ShopFeesView } from '@/components/(shop)/shop-fees/shop-fees-view';
import { RouteGuard } from '@/components/auth/RouteGuard';

export default function ShopFeesScreen() {
  return (
    <RouteGuard anyOf={['shop:fees:view']}>
      <ShopFeesView />
    </RouteGuard>
  );
}
