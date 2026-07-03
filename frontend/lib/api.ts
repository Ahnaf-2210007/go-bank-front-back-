const DEFAULT_API_BASE_URL = 'http://localhost:3000';

import type {
  AccountResponse,
  ApiResponse,
  LoginRequest,
  LoginResponse,
  OfferRequest,
  OfferResponse,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegisterRequest,
  RegisterResponse,
  TransactionRecord,
  TransferRequest,
  TransferResponse,
  WebAuthnLoginResponse,
  WebAuthnLoginBeginRequest,
  WebAuthnRegistrationBeginRequest,
  UpdateProfileRequest,
  VerifyEmailRequest,
} from './types';

function getApiBaseUrl() {
  return (process.env.NEXT_PUBLIC_API_URL || DEFAULT_API_BASE_URL)
    .trim()
    .replace(/^['"]|['"]$/g, '')
    .replace(/\/+$/g, '');
}

function apiUrl(path: string) {
  try {
    return new URL(path, `${getApiBaseUrl()}/`).toString();
  } catch {
    throw new Error('Invalid NEXT_PUBLIC_API_URL. Use a full URL like http://localhost:3000');
  }
}

async function readError(response: Response, fallback: string) {
  try {
    const error = await response.json();
    return error.error || error.message || fallback;
  } catch {
    return fallback;
  }
}

export const api = {
  async requestJson<T>(path: string, init?: RequestInit): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(apiUrl(path), init);

      if (!response.ok) {
        throw new Error(await readError(response, 'Request failed'));
      }

      return { data: await response.json() };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'An error occurred',
      };
    }
  },

  async login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    try {
      const response = await fetch(apiUrl('/login'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          number: Number(credentials.number),
          password: credentials.password,
        }),
      });

      if (!response.ok) {
        throw new Error(await readError(response, 'Login failed'));
      }

      return { data: await response.json() };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'An error occurred',
      };
    }
  },

  async register(data: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
    try {
      const response = await fetch(apiUrl('/account'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(await readError(response, 'Registration failed'));
      }

      return { data: await response.json() };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'An error occurred',
      };
    }
  },

  async verifyEmail(data: VerifyEmailRequest): Promise<ApiResponse<{ verified: boolean }>> {
    try {
      const response = await fetch(apiUrl('/account/verification'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(await readError(response, 'Verification failed'));
      }

      return { data: await response.json() };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'An error occurred',
      };
    }
  },

  async getAccount(token: string): Promise<ApiResponse<AccountResponse>> {
    try {
      const response = await fetch(apiUrl('/account'), {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(await readError(response, 'Failed to fetch account'));
      }

      return { data: await response.json() };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'An error occurred',
      };
    }
  },

  async transfer(token: string, data: TransferRequest): Promise<ApiResponse<TransferResponse>> {
    return this.requestJson<TransferResponse>('/transfer', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        toAccount: Number(data.toAccount),
        amount: Number(data.amount),
      }),
    });
  },

  async getTransactions(
    token: string,
    params: { limit?: number; offset?: number; type?: string; month?: string | number } = {},
  ): Promise<ApiResponse<TransactionRecord[]>> {
    const query = new URLSearchParams();

    if (typeof params.limit === 'number') {
      query.set('limit', String(params.limit));
    }
    if (typeof params.offset === 'number') {
      query.set('offset', String(params.offset));
    }
    if (params.type) {
      query.set('type', params.type);
    }
    if (params.month !== undefined && params.month !== null && `${params.month}`.trim() !== '') {
      query.set('month', String(params.month));
    }

    const suffix = query.toString() ? `?${query.toString()}` : '';

    return this.requestJson<TransactionRecord[]>(`/account/transactions${suffix}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });
  },

  async updateAccount(
    token: string,
    data: UpdateProfileRequest,
  ): Promise<ApiResponse<AccountResponse | { message: string }>> {
    return this.requestJson<AccountResponse | { message: string }>('/account/update', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  async redeemCoupon(
    token: string,
    accountId: number,
    couponCode: string,
  ): Promise<ApiResponse<OfferResponse>> {
    return this.requestJson<OfferResponse>(`/account/${accountId}/offer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ couponCode }),
    });
  },

  async webauthnRegisterBegin(
    token: string,
    data: WebAuthnRegistrationBeginRequest,
  ): Promise<ApiResponse<PublicKeyCredentialCreationOptionsJSON>> {
    return this.requestJson<PublicKeyCredentialCreationOptionsJSON>('/webauthn/register/begin', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  async webauthnRegisterFinish(
    token: string,
    credential: unknown,
  ): Promise<ApiResponse<{ status: string }>> {
    return this.requestJson<{ status: string }>('/webauthn/register/finish', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(credential),
    });
  },

  async webauthnLoginBegin(
    data: WebAuthnLoginBeginRequest,
  ): Promise<ApiResponse<PublicKeyCredentialRequestOptionsJSON>> {
    return this.requestJson<PublicKeyCredentialRequestOptionsJSON>('/webauthn/login/begin', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
  },

  async webauthnLoginFinish(
    email: string,
    credential: unknown,
  ): Promise<ApiResponse<WebAuthnLoginResponse>> {
    return this.requestJson<WebAuthnLoginResponse>(`/webauthn/login/finish/${encodeURIComponent(email)}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credential),
    });
  },

  async healthCheck(): Promise<boolean> {
    try {
      const response = await fetch(apiUrl('/health'));
      return response.ok;
    } catch {
      return false;
    }
  },
};
