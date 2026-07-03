export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface LoginRequest {
  number: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  number: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface RegisterResponse {
  message: string;
}

export interface VerifyEmailRequest {
  code: string;
}

export interface AccountResponse {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  number: number;
  balance: number;
  createdAt: string;
}

export interface TransferRequest {
  toAccount: number;
  amount: number;
}

export interface TransferResponse {
  status: string;
  transactionId: string;
  transferredAt: string;
}

export interface TransactionRecord {
  id: number;
  transactionId: string;
  fromAccountId: number;
  toAccountId: number;
  amount: number;
  transactionType: string;
  status: string;
  createdAt: string;
}

export interface OfferRequest {
  couponCode: string;
}

export interface OfferResponse {
  status: string;
}

export type ProfileUpdateAction = 'profile' | 'email_request' | 'email_verify' | 'password';

export interface UpdateProfileRequest {
  action: ProfileUpdateAction;
  firstName?: string;
  lastName?: string;
  newEmail?: string;
  password?: string;
  otp?: string;
  currentPassword?: string;
  newPassword?: string;
  confirmPassword?: string;
}

export interface WebAuthnRegistrationBeginRequest {
  email: string;
}

export interface WebAuthnLoginBeginRequest {
  email: string;
}

export interface WebAuthnLoginResponse {
  token: string;
  number: number;
}

export interface PublicKeyCredentialCreationOptionsJSON {
  publicKey: {
    challenge: string;
    rp?: {
      name: string;
      id?: string;
    };
    user?: {
      id: string;
      name: string;
      displayName: string;
    };
    pubKeyCredParams?: Array<{
      type: string;
      alg: number;
    }>;
    timeout?: number;
    excludeCredentials?: Array<{
      id: string;
      type: string;
      transports?: string[];
    }>;
    authenticatorSelection?: Record<string, unknown>;
    attestation?: string;
    extensions?: Record<string, unknown>;
  };
}

export interface PublicKeyCredentialRequestOptionsJSON {
  publicKey: {
    challenge: string;
    timeout?: number;
    rpId?: string;
    allowCredentials?: Array<{
      id: string;
      type: string;
      transports?: string[];
    }>;
    userVerification?: string;
    extensions?: Record<string, unknown>;
  };
}

