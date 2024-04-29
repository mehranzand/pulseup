import { configureStore } from '@reduxjs/toolkit'
import containerReducer from './slices/containerSlice'             
           

export const store = configureStore({
  reducer: {
    containers: containerReducer
  },
});
export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;