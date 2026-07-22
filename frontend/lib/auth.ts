const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

export interface StoredUser {
  number: number;
}

export const auth = {
  setToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(TOKEN_KEY, token);
    }
  },

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(TOKEN_KEY);
    }
    return null;
  },

  clearToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
    }
  },

  setUser(user: StoredUser): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(USER_KEY, JSON.stringify(user));
    }
  },

  getUser(): StoredUser | null {
    if (typeof window !== 'undefined') {
      const user = localStorage.getItem(USER_KEY);
      if (!user) {
        return null;
      }

      try {
        return JSON.parse(user) as StoredUser;
      } catch {
        localStorage.removeItem(USER_KEY);
        return null;
      }
    }
    return null;
  },

  isAuthenticated(): boolean {
    return Boolean(this.getToken());
  },

  logout(): void {
    this.clearToken();
  },
};
