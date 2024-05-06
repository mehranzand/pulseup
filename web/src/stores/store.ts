import { configureStore } from '@reduxjs/toolkit'
import containerReducer from './slices/containerSlice'                
import hostReducer from './slices/hostSlice';
           
export const store = configureStore({
  reducer: {
    containers: containerReducer,
    host: hostReducer
  },
});
export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;