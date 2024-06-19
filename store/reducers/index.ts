// store/reducers/index.ts
import {
  ActionReducerMap,
  MetaReducer
} from '@ngrx/store';
import { environment } from '../../environments/environment';
import { localStorageSync } from 'ngrx-store-localstorage';
import { itemReducer, ItemState } from './item.reducer';

export interface State {
  itemState: ItemState;
}

export const reducers: ActionReducerMap<State> = {
  itemState: itemReducer
};

export function localStorageSyncReducer(reducer: ActionReducer<any>): ActionReducer<any> {
  return localStorageSync({ keys: ['itemState'], rehydrate: true })(reducer);
}

export const metaReducers: MetaReducer<State>[] = !environment.production
  ? [localStorageSyncReducer]
  : [localStorageSyncReducer];
