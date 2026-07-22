import type { PublicKeyCredentialCreationOptionsJSON, PublicKeyCredentialRequestOptionsJSON } from './types';

const textEncoder = new TextEncoder();
const textDecoder = new TextDecoder();

function base64UrlToBuffer(value: string): ArrayBuffer {
  const base64 = value.replace(/-/g, '+').replace(/_/g, '/');
  const padding = '='.repeat((4 - (base64.length % 4)) % 4);
  const binary = atob(base64 + padding);
  const bytes = new Uint8Array(binary.length);

  for (let index = 0; index < binary.length; index += 1) {
    bytes[index] = binary.charCodeAt(index);
  }

  return bytes.buffer;
}

function bufferToBase64Url(value: ArrayBuffer | ArrayBufferView): string {
  const bytes = value instanceof ArrayBuffer ? new Uint8Array(value) : new Uint8Array(value.buffer, value.byteOffset, value.byteLength);
  let binary = '';

  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }

  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
}

export function webauthnSupported(): boolean {
  return typeof window !== 'undefined' && Boolean(window.PublicKeyCredential);
}

export function unsupportedWebauthnMessage() {
  return 'This browser does not support passkeys or WebAuthn. Use your password to continue.';
}

export function normalizeRegistrationOptions(options: PublicKeyCredentialCreationOptionsJSON): PublicKeyCredentialCreationOptions {
  const publicKey = options.publicKey;

  return {
    ...publicKey,
    challenge: base64UrlToBuffer(publicKey.challenge),
    user: publicKey.user
      ? {
          ...publicKey.user,
          id: base64UrlToBuffer(publicKey.user.id),
        }
      : undefined,
    excludeCredentials: Array.isArray(publicKey.excludeCredentials)
      ? publicKey.excludeCredentials.map((credential) => ({
          ...credential,
          id: base64UrlToBuffer(credential.id),
        }))
      : undefined,
  } as PublicKeyCredentialCreationOptions;
}

export function normalizeLoginOptions(options: PublicKeyCredentialRequestOptionsJSON): PublicKeyCredentialRequestOptions {
  const publicKey = options.publicKey;

  return {
    ...publicKey,
    challenge: base64UrlToBuffer(publicKey.challenge),
    allowCredentials: Array.isArray(publicKey.allowCredentials)
      ? publicKey.allowCredentials.map((credential) => ({
          ...credential,
          id: base64UrlToBuffer(credential.id),
        }))
      : undefined,
  } as PublicKeyCredentialRequestOptions;
}

export function normalizeLoginOptionsWithoutAllowList(options: PublicKeyCredentialRequestOptionsJSON): PublicKeyCredentialRequestOptions {
  const normalized = normalizeLoginOptions(options);

  return {
    ...normalized,
    allowCredentials: undefined,
  };
}

function serializeClientExtensionResults(results: AuthenticationExtensionsClientOutputs) {
  return Object.keys(results).length > 0 ? results : undefined;
}

export function serializeCreationCredential(credential: PublicKeyCredential) {
  const response = credential.response as AuthenticatorAttestationResponse;

  return {
    id: credential.id,
    rawId: bufferToBase64Url(credential.rawId),
    type: credential.type,
    response: {
      attestationObject: bufferToBase64Url(response.attestationObject),
      clientDataJSON: bufferToBase64Url(response.clientDataJSON),
    },
    clientExtensionResults: serializeClientExtensionResults(credential.getClientExtensionResults()),
  };
}

export function serializeRequestCredential(credential: PublicKeyCredential) {
  const response = credential.response as AuthenticatorAssertionResponse;

  return {
    id: credential.id,
    rawId: bufferToBase64Url(credential.rawId),
    type: credential.type,
    response: {
      authenticatorData: bufferToBase64Url(response.authenticatorData),
      clientDataJSON: bufferToBase64Url(response.clientDataJSON),
      signature: bufferToBase64Url(response.signature),
      userHandle: response.userHandle ? bufferToBase64Url(response.userHandle) : null,
    },
    clientExtensionResults: serializeClientExtensionResults(credential.getClientExtensionResults()),
  };
}

export function toUsername(value: string) {
  return textDecoder.decode(textEncoder.encode(value.trim()));
}
