// core.module.ts
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { ItemEffects } from '../store/effects/item.effects';
import { reducers } from '../store/reducers';
import { loadItems } from '../store/actions/item.actions';

@NgModule({
  imports: [
    CommonModule,
    StoreModule.forFeature('items', reducers),
    EffectsModule.forFeature([ItemEffects])
  ]
})
export class CoreModule {
  constructor(store: Store) {
    store.dispatch(loadItems());
  }
}
