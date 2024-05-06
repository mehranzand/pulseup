import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface HostState {
  current:  string | undefined
}

const initialState: HostState = {
  current: undefined
}

const hostSlice = createSlice({
  name: 'host',
  initialState,
  reducers: {
    setCurrent: (state, payloadAction: PayloadAction<any>) => {
      var { name } = payloadAction.payload
      state.current = name
    }
  }
})

export const { setCurrent } = hostSlice.actions

export default hostSlice.reducer