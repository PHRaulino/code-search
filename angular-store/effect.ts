// store/effects/item.effects.ts
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of, interval } from 'rxjs';
import { catchError, map, mergeMap, startWith, switchMap } from 'rxjs/operators';
import { ItemService } from '../../services/item.service';
import { updateItems, updateItemsSuccess, updateItemsFailure } from '../actions/item.actions';

@Injectable()
export class ItemEffects {
  updateItems$ = createEffect(() =>
    this.actions$.pipe(
      ofType(updateItems),
      switchMap(() =>
        this.itemService.getItems().pipe(
          map(items => updateItemsSuccess({ items })),
          catchError(error => of(updateItemsFailure({ error })))
        )
      )
    )
  );

  // Efeito para disparar a ação de atualização a cada 2 minutos
  periodicUpdate$ = createEffect(() =>
    interval(120000).pipe( // 120000 ms = 2 minutos
      startWith(0),
      map(() => updateItems())
    )
  );

  constructor(private actions$: Actions, private itemService: ItemService) {}
}
