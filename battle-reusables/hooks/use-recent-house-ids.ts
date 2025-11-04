import { useCallback, useEffect, useState } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';

const STORAGE_KEY = 'recent:house_ids';

export function useRecentHouseIds(limit: number = 8) {
  const [list, setList] = useState<string[]>([]);

  useEffect(() => {
    (async () => {
      try {
        const v = await AsyncStorage.getItem(STORAGE_KEY);
        if (v) setList(JSON.parse(v));
      } catch {}
    })();
  }, []);

  const persist = useCallback(async (arr: string[]) => {
    setList(arr);
    try { await AsyncStorage.setItem(STORAGE_KEY, JSON.stringify(arr)); } catch {}
  }, []);

  const add = useCallback(async (id: number | string) => {
    const s = String(id).trim();
    if (!s) return;
    const next = [s, ...list.filter((x) => x !== s)].slice(0, Math.max(1, limit));
    await persist(next);
  }, [limit, list, persist]);

  const clear = useCallback(async () => {
    await persist([]);
  }, [persist]);

  const suggestions = useCallback((q: string) => {
    const s = String(q ?? '').trim();
    if (!s) return list.slice(0, limit);
    return list.filter((x) => x.startsWith(s)).slice(0, limit);
  }, [list, limit]);

  return { list, add, suggestions, clear };
}


