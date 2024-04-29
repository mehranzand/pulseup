import { PayloadAction, createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { Container } from '../../types/Container'
import { API_URLS } from "../../configs/api"
import axios from "axios"

interface ContainerState {
  loading: boolean
  data: Container[],
  error: string | null
}

export const fetchContainers = createAsyncThunk("container/fetchData", async (host: string, thunkApi) => {
  try {
    var url = API_URLS.containers_url.replace(':host', host)
    const response = await axios.get<Container[]>(url)
    return response.data
  } catch (err: any) {
    return thunkApi.rejectWithValue(err.message)
  }
})

const initialState: ContainerState = {
  loading: false,
  data: [],
  error: null
}

const containerSlice = createSlice({
  name: 'containers',
  initialState,
  reducers: {
    handleSourceEvent: (state, payloadAction: PayloadAction<any>) => {
      var { action, container } = payloadAction.payload
      if (action == 'create') {
        state.data?.unshift(container)
        return
      }

      var index = state.data?.findIndex(a => a.id == container.id)
      if (index != -1 && index !== undefined) {
        if (action == 'destroy') {
          state.data?.splice(index, 1)
        }
        else { //die, start
          state.data[index] = container
        }
      }
    }
  },
  extraReducers: (builder) => {
    builder.addCase(fetchContainers.pending, (state) => {
      state.loading = true
    })

    builder.addCase(fetchContainers.fulfilled, (state, action: PayloadAction<Container[]>) => {
      state.loading = false
      state.data = action.payload
    })

    builder.addCase(fetchContainers.rejected, (state, action: PayloadAction<any>) => {
      state.loading = false
      state.error = action.payload
    })
  }
})

export const { handleSourceEvent } = containerSlice.actions

export default containerSlice.reducer