// store/actions/item.actions.ts
import { createAction, props } from '@ngrx/store';
import { Item } from '../../models/item.model';

export const updateItems = createAction('[Item] Update Items');
export const updateItemsSuccess = createAction('[Item] Update Items Success', props<{ items: Item[] }>());
export const updateItemsFailure = createAction('[Item] Update Items Failure', props<{ error: any }>());
