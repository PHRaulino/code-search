// store/effects/group.effects.ts
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs/operators';
import { GroupService } from '../../shared/services/group.service';
import { loadGroups, loadGroupsSuccess, loadGroupsFailure } from '../actions/group.actions';

@Injectable()
export class GroupEffects {
  loadGroups$ = createEffect(() =>
    this.actions$.pipe(
      ofType(loadGroups),
      mergeMap(({ userId }) =>
        this.groupService.getGroups(userId).pipe(
          map((groups) => loadGroupsSuccess({ groups })),
          catchError((error) => of(loadGroupsFailure({ error })))
        )
      )
    )
  );

  constructor(private actions$: Actions, private groupService: GroupService) {}
}
