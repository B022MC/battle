import { useEffect, useMemo, useState } from 'react';
import { platformsPlazaConsts } from '@/services/platforms';

type LabelMaps = {
  modes: Map<number, string>;
  scenes: Map<number, string>;
  user_status: Map<number, string>;
  member_types: Map<number, string>;
  game_genre: Map<number, string>;
  game_kinds: Map<number, string>;
  table_genre: Map<number, string>;
  system_message_types: Map<number, string>;
  urls: Map<string, string>;
};

type NumericLabelGroup = Exclude<keyof LabelMaps, 'urls'>;

// Module-level cache to ensure the constants are fetched only once per app lifecycle
let cachedData: API.PlazaConsts | null = null;
let inflight: Promise<API.PlazaConsts> | null = null;
const listeners = new Set<(data: API.PlazaConsts) => void>();

async function fetchOnce(): Promise<API.PlazaConsts> {
  if (cachedData) return cachedData;
  if (!inflight) {
    inflight = platformsPlazaConsts().then((res) => {
      const data = res.data as API.PlazaConsts;
      cachedData = data;
      inflight = null;
      listeners.forEach((fn) => fn(data));
      return data;
    }).catch((e) => {
      inflight = null;
      throw e;
    });
  }
  return inflight;
}

export function usePlazaConsts() {
  const [data, setData] = useState<API.PlazaConsts | undefined>(cachedData ?? undefined);
  const [loading, setLoading] = useState<boolean>(!cachedData);
  const [error, setError] = useState<Error | undefined>(undefined);

  useEffect(() => {
    if (cachedData) {
      setData(cachedData);
      setLoading(false);
      return;
    }
    const onData = (d: API.PlazaConsts) => {
      setData(d);
      setLoading(false);
    };
    listeners.add(onData);
    fetchOnce().catch((e) => { setError(e as Error); setLoading(false); });
    return () => { listeners.delete(onData); };
  }, []);

  const maps: LabelMaps = useMemo(() => {
    const m: LabelMaps = {
      modes: new Map<number, string>(),
      scenes: new Map<number, string>(),
      user_status: new Map<number, string>(),
      member_types: new Map<number, string>(),
      game_genre: new Map<number, string>(),
      game_kinds: new Map<number, string>(),
      table_genre: new Map<number, string>(),
      system_message_types: new Map<number, string>(),
      urls: new Map<string, string>(),
    };
    if (data) {
      data.modes?.forEach(i => m.modes.set(i.value, i.label));
      data.scenes?.forEach(i => m.scenes.set(i.value, i.label));
      data.user_status?.forEach(i => m.user_status.set(i.value, i.label));
      data.member_types?.forEach(i => m.member_types.set(i.value, i.label));
      data.game_genre?.forEach(i => m.game_genre.set(i.value, i.label));
      data.game_kinds?.forEach(i => m.game_kinds.set(i.value, i.label));
      data.table_genre?.forEach(i => m.table_genre.set(i.value, i.label));
      data.system_message_types?.forEach(i => m.system_message_types.set(i.value, i.label));
      data.urls?.forEach(u => m.urls.set(u.key, u.value));
    }
    return m;
  }, [data]);

  function getLabel(group: NumericLabelGroup, value: number): string {
    return maps[group].get(value) ?? String(value);
  }

  function getLoginModeLabel(mode?: string | null): string {
    if (!mode) return '-';
    if (mode === 'account') return '游戏账号';
    if (mode === 'mobile') return '手机号';
    return String(mode);
  }

  const refresh = async () => {
    cachedData = null;
    setLoading(true);
    try {
      const d = await fetchOnce();
      setData(d);
      setError(undefined);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return { data, loading, error, refresh, maps, getLabel, getLoginModeLabel };
}


