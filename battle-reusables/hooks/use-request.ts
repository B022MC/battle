import { ResponseStructure } from '@/utils/request';
import { useCallback, useEffect, useMemo, useReducer } from 'react';

type State<T> = {
  data?: T;
  loading: boolean;
  error?: Error;
};

type Action<T> =
  | { type: 'REQUEST' }
  | { type: 'SUCCESS'; payload?: T }
  | { type: 'FAILURE'; payload?: Error };

// Reducer 纯函数，管理状态的变更逻辑
const reducer = <T>(state: State<T>, action: Action<T>): State<T> => {
  switch (action.type) {
    case 'REQUEST':
      return {
        ...state,
        loading: true,
        error: undefined,
      };
    case 'SUCCESS':
      return {
        ...state,
        loading: false,
        data: action.payload,
      };
    case 'FAILURE':
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    default:
      return state;
  }
};

type Service<T, P extends any[]> = (...args: P) => Promise<ResponseStructure<T>>;

type Options<T, P extends any[]> = {
  params?: P;
  manual?: boolean;
  onSuccess?: (data?: T) => void;
  onError?: (error: Error) => void;
};

export const useRequest = <T, P extends any[]>(
  service: Service<T, P>,
  options: Options<T, P> = {}
) => {
  const { params, manual, onSuccess, onError } = options;

  const [state, dispatch] = useReducer(reducer, { loading: false } as State<T>);

  const run = useCallback(
    async (...params: P) => {
      dispatch({ type: 'REQUEST' });
      try {
        const { data } = await service(...params);
        dispatch({ type: 'SUCCESS', payload: data });
        onSuccess?.(data);
        return data;
      } catch (err: any) {
        dispatch({ type: 'FAILURE', payload: err });
        onError?.(err);
        throw err;
      }
    },
    [service, onSuccess, onError]
  );

  const paramsStr = useMemo(() => JSON.stringify(params), [params]);

  useEffect(() => {
    if (!manual) run(...((params ?? []) as P));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [manual, run, paramsStr]);

  return { ...state, run };
};
