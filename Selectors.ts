// store/selectors/group.selectors.ts
import { createSelector, createFeatureSelector } from '@ngrx/store';
import { GroupState } from '../reducers/group.reducer';

export const selectGroupState = createFeatureSelector<GroupState>('group');

export const selectAllGroups = createSelector(
  selectGroupState,
  (state: GroupState) => state.groups
);

export const selectGroupError = createSelector(
  selectGroupState,
  (state: GroupState) => state.error
);
