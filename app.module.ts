// app.module.ts
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { StoreRouterConnectingModule } from '@ngrx/router-store';
import { StoreLocalStorageModule } from '@ngrx/store-localstorage';
import { AppComponent } from './app.component';
import { reducers, metaReducers } from './store/reducers';
import { GroupEffects } from './store/effects/group.effects';
import { UserEffects } from './store/effects/user.effects';

@NgModule({
  declarations: [AppComponent],
  imports: [
    BrowserModule,
    StoreModule.forRoot(reducers, { metaReducers }),
    EffectsModule.forRoot([GroupEffects, UserEffects]),
    StoreDevtoolsModule.instrument({ maxAge: 25 }),
    StoreRouterConnectingModule.forRoot(),
    StoreLocalStorageModule.forRoot({
      keys: ['itemState'], // Estado(s) que deseja persistir
      rehydrate: true // Carregar o estado persistido na inicialização
    })
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
