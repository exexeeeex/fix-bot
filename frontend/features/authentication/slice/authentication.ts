import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { AuthenticationResponse, TokensResponse } from '../types';
import { authenticationApi } from '../api/authentication';

export interface AuthenticationSliceState {
  user: AuthenticationResponse['user'] | null;
  tokens: TokensResponse | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

const initialState: AuthenticationSliceState = {
  isAuthenticated: false,
  user: null,
  tokens: null,
  isLoading: false,
  error: null,
};

const authenticationSlice = createSlice({
  name: 'authentication',
  initialState,
  reducers: {
    setTokens: (state, action: PayloadAction<TokensResponse>) => {
      state.tokens = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addMatcher(
        authenticationApi.endpoints.login.matchFulfilled,
        (state, action: PayloadAction<AuthenticationResponse>) => {
          const response = action.payload;
          state.isLoading = false;

          if (response.success) {
            state.tokens = response.tokens;
            state.user = response.user;
            state.isAuthenticated = true;
            state.error = null;
          } else {
            state.error = 'Login failed';
          }
        }
      )
      .addMatcher(authenticationApi.endpoints.login.matchRejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Login failed';
      })
      .addMatcher(authenticationApi.endpoints.login.matchPending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addMatcher(
        authenticationApi.endpoints.demoLogin.matchFulfilled,
        (state, action: PayloadAction<AuthenticationResponse>) => {
          const response = action.payload;
          state.isLoading = false;

          if (response.success) {
            state.tokens = response.tokens;
            state.user = response.user;
            state.isAuthenticated = true;
            state.error = null;
          } else {
            state.error = 'Login failed';
          }
        }
      )
      .addMatcher(authenticationApi.endpoints.demoLogin.matchRejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Login failed';
      })
      .addMatcher(authenticationApi.endpoints.demoLogin.matchPending, (state) => {
        state.isLoading = true;
        state.error = null;
      });
  },
});

export default authenticationSlice.reducer;
