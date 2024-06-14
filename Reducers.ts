// store/reducers/group.reducer.ts
import { createReducer, on } from '@ngrx/store';
import { loadGroups, loadGroupsSuccess, loadGroupsFailure } from '../actions/group.actions';
import { Group } from '../../shared/interfaces/group.model';

export interface GroupState {
  groups: Group[];
  error: any;
}

const initialState: GroupState = {
  groups: [],
  error: null
};

export const groupReducer = createReducer(
  initialState,
  on(loadGroupsSuccess, (state, { groups }) => ({ ...state, groups, error: null })),
  on(loadGroupsFailure, (state, { error }) => ({ ...state, error }))
);
