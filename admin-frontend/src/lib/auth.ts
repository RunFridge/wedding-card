import { ref, readonly } from 'vue';
import api from './axios';

const authenticated = ref(false);
const passwordNeedsChange = ref(false);
const setupRequired = ref(false);

export function useAuth() {
  async function login(password: string): Promise<boolean> {
    const res = await api.post('/admin/verify', { password });
    if (res.data.status === 'ok') {
      authenticated.value = true;
      passwordNeedsChange.value = !!res.data.password_needs_change;
      setupRequired.value = !!res.data.setup_required;
      return true;
    }
    return false;
  }

  async function checkSession(): Promise<boolean> {
    try {
      const res = await api.get('/admin/session');
      authenticated.value = true;
      passwordNeedsChange.value = !!res.data.password_needs_change;
      setupRequired.value = !!res.data.setup_required;
      return true;
    } catch {
      authenticated.value = false;
      passwordNeedsChange.value = false;
      setupRequired.value = false;
      return false;
    }
  }

  async function logout() {
    try {
      await api.post('/admin/logout');
    } catch {
      // ignore — cookie may already be expired
    }
    authenticated.value = false;
    passwordNeedsChange.value = false;
    setupRequired.value = false;
  }

  function clearPasswordNeedsChange() {
    passwordNeedsChange.value = false;
  }

  function clearSetupRequired() {
    setupRequired.value = false;
  }

  return {
    authenticated: readonly(authenticated),
    passwordNeedsChange: readonly(passwordNeedsChange),
    setupRequired: readonly(setupRequired),
    login,
    logout,
    checkSession,
    clearPasswordNeedsChange,
    clearSetupRequired,
  };
}
